package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/auth/mocks"
	"github.com/ffigari/stored-strings/internal/basesuite"
)

//go:generate mockgen -package=mocks -source=main.go -destination=mocks/main.go

func TestNewTokensAreValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	clock := mocks.NewMockclock(ctrl)
	authenticator := auth.New([]byte("foo"), clock)

	clock.EXPECT().Now().Return(time.Now())
	token := authenticator.GenerateToken()
	require.NotEqual(t, token, "")

	isValid := authenticator.IsValidToken(token)
	assert.True(t, isValid)
}

func TestTokenOfOneSecretWontBeValidForAnotherSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	oneClock := mocks.NewMockclock(ctrl)
	oneAuthenticator := auth.New([]byte("foo"), oneClock)

	oneClock.EXPECT().Now().Return(time.Now())
	token := oneAuthenticator.GenerateToken()
	require.NotEqual(t, token, "")

	require.True(t, oneAuthenticator.IsValidToken(token))

	anotherAuthenticator := auth.New([]byte("faa"), mocks.NewMockclock(ctrl))
	require.False(t, anotherAuthenticator.IsValidToken(token))
}

func TestTokensCreatedInTheSameMomentAreDifferent(t *testing.T) {
	ctrl := gomock.NewController(t)
	clock := mocks.NewMockclock(ctrl)
	authenticator := auth.New([]byte("foo"), clock)

	now := time.Now()
	clock.EXPECT().Now().Return(now)
	oneToken := authenticator.GenerateToken()

	clock.EXPECT().Now().Return(now)
	anotherToken := authenticator.GenerateToken()

	require.NotEqual(t, oneToken, anotherToken)
}

func TestOldTokensAreMarkedAsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	clock := mocks.NewMockclock(ctrl)
	authenticator := auth.New([]byte("foo"), clock)

	clock.EXPECT().Now().Return(time.Now().AddDate(0, 0, -4))

	token := authenticator.GenerateToken()
	require.NotEqual(t, token, "")

	assert.False(t, authenticator.IsValidToken(token))
}

func (s *Suite) TestLoginService() {
	loginPath := s.server.URL + "/login"

	s.Run("login form is provided", func() {
		res, err := http.Get(loginPath)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)

		body := s.GetBody(res)
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

		body := s.GetBody(res)
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

		body := s.GetBody(res)
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

		body := s.GetBody(res)
		s.Require().Contains(body, "logged in")

		authorizationToken := s.GetCookieValue(res, "authorization")
		s.Require().NotEmpty(authorizationToken)

		s.Require().True(s.authenticator.IsValidToken(authorizationToken))
	})
}

func (s *Suite) SetupSuite() {
	testDBName := "test_db"

	s.authenticator = s.SetupAuthenticator()

	s.dbPool = s.SetupDB(context.Background(), testDBName)

	s.password = "password1234"

	r := mux.NewRouter()

	auth.AttachTo(r, s.password, s.dbPool, s.authenticator)

	s.server = httptest.NewServer(r)
}

type Suite struct {
	basesuite.BaseSuite
	server        *httptest.Server
	password      string
	authenticator *auth.Authenticator
	dbPool        *pgxpool.Pool
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
