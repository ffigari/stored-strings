package personalwebsite_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/basesuite"
	"github.com/ffigari/stored-strings/internal/personalwebsite"
)

func (s *Suite) TestHome() {
	assertOKResponse := func(res *http.Response) {
		s.Require().Equal(http.StatusOK, res.StatusCode)
		body := s.GetBody(res)
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

	s.Run("serves individual media files", func() {
		res, err := http.Get(s.server.URL + "/i/panoramica-de-casa.jpg")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)
	})

	s.Run("lists all media files", func() {
		fmt.Println(s.server.URL)
		res, err := http.Get(s.server.URL + "/i/")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.GetBody(res)
		s.Require().Contains(body, "panoramica-de-casa.jpg")
	})
}

func (s *Suite) TestLogin() {
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
	basesuite.BaseSuite
	server        *httptest.Server
	password      string
	authenticator *auth.Authenticator
	dbPool        *pgxpool.Pool
}

func (s *Suite) SetupSuite() {
	ctx := context.Background()
	testDBName := "test_db"

	s.password = "password1234"
	s.authenticator = s.SetupAuthenticator()
	s.dbPool = s.SetupDB(ctx, testDBName)

	conn, err := s.dbPool.Acquire(ctx)
	s.Require().NoError(err)
	defer conn.Release()

	m, err := personalwebsite.NewMux(
		context.Background(), testDBName, s.authenticator, s.password,
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
