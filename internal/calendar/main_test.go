package calendar_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/basesuite"
	"github.com/ffigari/stored-strings/internal/calendar"
)

func (s *CalendarSuite) TestForEach() {
	ctx := context.Background()

	conn, err := s.dbPool.Acquire(ctx)
	s.Require().NoError(err)
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM events`)
	s.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		INSERT INTO events (starts_at, description)
		VALUES
			('2024-06-21 18:50:00', 'merendar'),
			('2024-06-21 12:30:00', '108 saludos al sol'),
			('2024-08-26 09:30:00', 'pasear en bici'),
			('2024-06-26 22:30:00', 'ver pelis')
		;
	`)
	s.Require().NoError(err)

	s.Run("events are sorted by start time", func () {
		ctx := context.Background()

		conn, err := s.dbPool.Acquire(ctx)
		s.Require().NoError(err)
		defer conn.Release()

		sortedDescriptions := []string{
			"108 saludos al sol", "merendar", "ver pelis", "pasear en bici",
		}
		i := 0
		s.Require().NoError(calendar.ForEach(ctx, conn, func(
			_ time.Time, description string,
		) {
			s.Require().True(i <= 3)
			s.Require().Equal(sortedDescriptions[i], description)
			i++
		}))

		s.Require().Equal(4, i)
	})
}

func (s *CalendarSuite) TestAttachTo() {
	ctx := context.Background()

	conn, err := s.dbPool.Acquire(ctx)
	s.Require().NoError(err)
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM events`)
	s.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		INSERT INTO events (starts_at, description)
		VALUES
			('2024-06-21 18:50:00', 'merendar'),
			('2024-06-21 12:30:00', '108 saludos al sol'),
			('2024-08-26 09:30:00', 'pasear en bici'),
			('2024-06-26 22:30:00', 'ver pelis')
		;
	`)
	s.Require().NoError(err)

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
		res := doLoggedInRequest("GET", "/calendar", nil)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.GetBody(res)
		s.Require().Contains(body, "2024-06-21 09:30: 108 saludos al sol")
		s.Require().Contains(body, "2024-06-26 19:30: ver pelis")
	})

	s.Run("form to add an event is provided", func() {
		res := doLoggedInRequest("GET", "/calendar/events", nil)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.GetBody(res)
		s.Require().Contains(body, `form method="POST"`)
		s.Require().Contains(body, `name="starts_at"`)
		s.Require().Contains(body, `name="description"`)
	})

	s.Run("new events can be created", func() {
		newDescription := "pasarla lindo"

		ctx := context.Background()

		conn, err := s.dbPool.Acquire(ctx)
		s.Require().NoError(err)
		defer conn.Release()

		found := false
		calendar.ForEach(ctx, conn, func(_ time.Time, description string) {
			found = found || description == newDescription
		})
		s.Require().False(found)

		res := doLoggedInRequest(
			"POST",
			"/calendar/events",
			strings.NewReader(fmt.Sprintf(
				"starts_at=%s&description=%s",
				"2024-07-23T19:15",
				newDescription,
			)),
		)

		body := s.GetBody(res)
		s.Require().Contains(body, "2024-06-21 09:30: 108 saludos al sol")
		s.Require().Contains(body, "2024-06-26 19:30: ver pelis")
		s.Require().Contains(body, "2024-07-23 19:15: pasarla lindo")

		s.Require().Equal(http.StatusOK, res.StatusCode)

		found = false
		calendar.ForEach(ctx, conn, func(_ time.Time, description string) {
			found = found || description == newDescription
		})
		s.Require().True(found)
	})
}


type CalendarSuite struct {
	basesuite.BaseSuite
	server        *httptest.Server
	authenticator *auth.Authenticator
	dbPool        *pgxpool.Pool
}

func (s *CalendarSuite) SetupSuite() {
	ctx := context.Background()
	testDBName := "test_db"

	s.authenticator = s.SetupAuthenticator()

	s.dbPool = s.SetupDB(ctx, testDBName)

	r := mux.NewRouter()

	s.Require().NoError(calendar.AttachTo(r, "", s.dbPool, s.authenticator))

	s.server = httptest.NewServer(r)
}

func (s *CalendarSuite) TearDownSuite() {
	s.dbPool.Close()
	s.server.Close()
}

func TestCalendar(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
