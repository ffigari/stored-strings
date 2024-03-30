package personalwebsite_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/calendar"
	"github.com/ffigari/stored-strings/internal/dbpool"
	"github.com/ffigari/stored-strings/internal/personalwebsite"
	"github.com/ffigari/stored-strings/internal/postgresql"
)

var consecutiveSpacesRegexp = regexp.MustCompile(`\s+`)

func (s *Suite) getBody(res *http.Response) string {
	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err)

	return consecutiveSpacesRegexp.ReplaceAllString(string(body), " ")
}

func (s *Suite) getCookieValue(res *http.Response, name string) string {
	for _, cookie := range res.Cookies() {
		if cookie.Name != name {
			continue
		}

		return cookie.Value
	}

	s.Require().Empty(fmt.Sprintf("expected cookie '%s' not found", name))
	return ""
}

func (s *Suite) TestHome() {
	assertOKResponse := func(res *http.Response) {
		s.Require().Equal(http.StatusOK, res.StatusCode)
		body := s.getBody(res)
		s.Require().Contains(body, "desde entonces vivo en Buenos Aires")
	}

	s.Run("is served at '/i'", func() {
		res, err := http.Get(s.server.URL + "/i")
		s.Require().NoError(err)
		assertOKResponse(res)
	})

	s.Run("is redirected to from '/'", func() {
		res, err := http.Get(s.server.URL + "/")
		s.Require().NoError(err)
		assertOKResponse(res)
	})

	s.Run("is not redirected to from '/.+'", func() {
		res, err := http.Get(s.server.URL + "/lalala")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, res.StatusCode)
	})
}

func (s *Suite) TestLogin() {
	loginPath := s.server.URL + "/i/login"

	s.Run("login form is provided", func() {
		res, err := http.Get(loginPath)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.getBody(res)
		s.Assert().Contains(body, `form method="POST"`)
		s.Assert().Contains(body, `name="username"`)
		s.Assert().Contains(body, `name="password"`)
	})

	// TODO: Ver si tiene mas sentido delegar estos chequeos de params al
	// paquete `webapi`. Si no siempre tendria que chequear el msg del body y el
	// status code cuando en verdad estoy usando siempre el mismo codigo
	s.Run("username is requested", func() {
		res, err := http.Post(
			loginPath,
			"application/x-www-form-urlencoded",
			strings.NewReader("username=foo"),
		)
		s.Require().NoError(err)

		body := s.getBody(res)
		s.Require().Contains(body, "missing required body params [password]")
	})

	s.Run("invalid credentials", func() {
		res, err := http.Post(
			loginPath,
			"application/x-www-form-urlencoded",
			strings.NewReader("username=foo&password=bar"),
		)
		s.Require().NoError(err)
		s.Require().NotNil(res)
		s.Require().Equal(http.StatusUnauthorized, res.StatusCode)

		body := s.getBody(res)
		s.Require().Contains(body, "invalid credentials")
	})

	s.Run("valid token is returned on successful login", func() {
		res, err := http.Post(
			loginPath,
			"application/x-www-form-urlencoded",
			strings.NewReader(fmt.Sprintf(
				"username=admin&password=%s", s.password,
			)),
		)
		s.Require().NoError(err)
		s.Require().NotNil(res)
		s.Assert().Equal(http.StatusOK, res.StatusCode)

		body := s.getBody(res)
		s.Require().Contains(body, "logged in")

		authorizationToken := s.getCookieValue(res, "authorization")
		s.Require().NotEmpty(authorizationToken)

		s.Require().True(s.authenticator.IsValidToken(authorizationToken))
	})
}

func (s *Suite) TestCalendar() {
	doLoggedInRequest := func(method, path string, body io.Reader) *http.Response {
		req, err := http.NewRequest(method, s.server.URL+path, body)
		s.Require().NoError(err)

		req.AddCookie(&http.Cookie{
			Name:  "authorization",
			Value: s.authenticator.GenerateToken(),
		})

		res, err := (&http.Client{}).Do(req)
		s.Require().NoError(err)

		return res
	}

	s.Run("calendar can be retrieved", func() {
		res := doLoggedInRequest("GET", "/i/calendar", nil)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.getBody(res)
		s.Require().Contains(body, "8 de marzo: comer rico")
		s.Require().Contains(body, "25 de mayo: tomar mate")
	})

	s.Run("form to add an event is provided", func() {
		res := doLoggedInRequest("GET", "/i/calendar/events", nil)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.getBody(res)
		s.Require().Contains(body, `form method="POST"`)
		s.Require().Contains(body, `name="date"`)
		s.Require().Contains(body, `name="description"`)
	})

	s.Run("new events can be created", func() {
		newDate := "15 de agosto"
		newDescription := "pasarla lindo"

		ctx := context.Background()

		conn, err := s.dbPool.Acquire(ctx)
		s.Require().NoError(err)
		defer conn.Release()

		found := false
		calendar.ForEach(ctx, conn, func(date, event string) {
			found = date == newDate && event == newDescription
		})
		s.Require().False(found)

		rows, err := conn.Query(ctx, `SELECT date, event FROM calendar;`)
		s.Require().NoError(err)
		for rows.Next() {
		}
		res := doLoggedInRequest(
			"POST",
			"/i/calendar/events",
			strings.NewReader(fmt.Sprintf(
				"date=%s&description=%s", newDate, newDescription,
			)),
		)

		body := s.getBody(res)
		s.Require().Contains(body, "8 de marzo: comer rico")
		s.Require().Contains(body, "25 de mayo: tomar mate")
		s.Require().Contains(body, fmt.Sprintf(
			"%s: %s", newDate, newDescription,
		))

		s.Require().Equal(http.StatusOK, res.StatusCode)

		found = false
		calendar.ForEach(ctx, conn, func(date, event string) {
			found = date == newDate && event == newDescription
		})
		s.Require().True(found)
	})
}

func (s *Suite) TestStatusIsOffered() {
	res, err := http.Get(s.server.URL + "/status")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, res.StatusCode)
}

func (s *Suite) TestFaviconIsNotProvided() {
	res, err := http.Get(s.server.URL + "/favicon.ico")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, res.StatusCode)
}

type Suite struct {
	suite.Suite
	server        *httptest.Server
	password      string
	authenticator *auth.Authenticator
	dbPool        *pgxpool.Pool
}

func (s *Suite) SetupSuite() {
	s.password = "password1234"

	authenticator, err := auth.NewFromConfig()
	s.Require().NoError(err)
	s.authenticator = authenticator

	testDBName := "test_db"

	s.Require().NoError(postgresql.CreateEmptyDB(testDBName))

	ctx := context.Background()
	dbPool, err := dbpool.NewFromConfig(ctx, testDBName)
	s.Require().NoError(err)
	s.dbPool = dbPool

	conn, err := dbPool.Acquire(ctx)
	s.Require().NoError(err)
	defer conn.Release()

	s.Require().NoError(postgresql.RunMigrations(ctx, conn))

	_, err = conn.Exec(ctx, `
		INSERT INTO calendar (date, event)
		VALUES
			('8 de marzo', 'comer rico'),
			('25 de mayo', 'tomar mate')
		;
	`)
	s.Require().NoError(err)

	m, err := personalwebsite.NewMux(
		context.Background(), testDBName, authenticator, s.password,
	)
	s.Require().NoError(err)

	s.server = httptest.NewServer(m)
}

func (s *Suite) TearDownSuite() {
	s.dbPool.Close()
	s.server.Close()
}

func TestPersonalWebsite(m *testing.T) {
	suite.Run(m, new(Suite))
}
