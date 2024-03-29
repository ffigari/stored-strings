package auth_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/auth/mocks"
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
