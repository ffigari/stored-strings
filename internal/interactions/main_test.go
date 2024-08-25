package interactions_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/basesuite"
	"github.com/ffigari/stored-strings/internal/interactions"
)

func (s *InteractionsSuite) TestClicks() {
	ctx := context.Background()

	conn, err := s.dbPool.Acquire(ctx)
	s.Require().NoError(err)
	defer conn.Release()

	fn := func(body string) *http.Response {
		req, err := http.NewRequest(
			"POST",
			s.server.URL+"/interactions/clicks",
			strings.NewReader(body),
		)
		s.Require().NoError(err)

		return s.SendReq(req)
	}

	s.Run("new clicks can be persisted", func() {
		res := fn(`{
			"x": 100,
			"y": 50
		}`)

		s.Require().Equal(http.StatusCreated, res.StatusCode)

		rows, err := conn.Query(ctx, `
			SELECT x, y
			FROM clicks
		;`)
		s.Require().NoError(err)

		count := 0
		for rows.Next() {
			s.Require().Equal(0, count)
			var x, y int

			s.Require().NoError(rows.Scan(&x, &y))

			s.Require().Equal(100, x)
			s.Require().Equal(50, y)

			count++
		}

		s.Require().Equal(1, count)
	})

	s.Run("clicks can be retrieved", func() {
		_, err := conn.Exec(context.Background(), `
			DELETE FROM clicks;
			INSERT INTO clicks (x, y)
			VALUES (100, 200), (-50, 0);
		`)
		s.Require().NoError(err)

		req, err := http.NewRequest(
			"GET",
			s.server.URL+"/interactions/clicks",
			nil,
		)
		s.Require().NoError(err)

		res := s.SendReq(req)
		s.Require().Equal(http.StatusOK, res.StatusCode)
		s.Require().JSONEq(`[{
			"x": 100, "y": 200
		}, {
			"x": -50, "y": 0
		}]`, s.GetBody(res))
	})
}

func (s *InteractionsSuite) SetupSuite() {
	s.dbPool = s.SetupDB(context.Background(), "test_db")

	r := mux.NewRouter()

	interactions.AttachTo(r, s.dbPool)

	s.server = httptest.NewServer(r)
}

func (s *InteractionsSuite) TearDownSuite() {
	s.server.Close()
	s.dbPool.Close()
}

type InteractionsSuite struct {
	basesuite.BaseSuite
	server        *httptest.Server
	dbPool *pgxpool.Pool
}

func TestInteractions(t *testing.T) {
	suite.Run(t, new(InteractionsSuite))
}
