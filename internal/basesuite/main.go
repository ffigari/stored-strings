package basesuite

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/dbpool"
	"github.com/ffigari/stored-strings/internal/postgresql"
)

type BaseSuite struct {
	suite.Suite
}

var consecutiveSpacesRegexp = regexp.MustCompile(`\s+`)

func (s *BaseSuite) GetBody(res *http.Response) string {
	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err)

	return consecutiveSpacesRegexp.ReplaceAllString(string(body), " ")
}

func (s *BaseSuite) SendReq(req *http.Request) *http.Response {
	res, err := (&http.Client{}).Do(req)
	s.Require().NoError(err)

	return res
}

func (s *BaseSuite) GetCookieValue(res *http.Response, name string) string {
	for _, cookie := range res.Cookies() {
		if cookie.Name != name {
			continue
		}

		return cookie.Value
	}

	s.Require().Empty(fmt.Sprintf("expected cookie '%s' not found", name))
	return ""
}


func (s *BaseSuite) SetupDB(ctx context.Context, testDBName string) *pgxpool.Pool {
	s.Require().NoError(postgresql.CreateEmptyDB(testDBName))

	dbPool, err := dbpool.NewFromConfig(ctx, testDBName)
	s.Require().NoError(err)

	conn, err := dbPool.Acquire(ctx)
	s.Require().NoError(err)
	defer conn.Release()

	s.Require().NoError(postgresql.RunMigrations(ctx, conn))

	return dbPool
}

func (s *BaseSuite) SetupAuthenticator() *auth.Authenticator {
	authenticator, err := auth.NewFromConfig()
	s.Require().NoError(err)

	return authenticator
}
