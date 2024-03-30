package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	clockP "github.com/ffigari/stored-strings/internal/clock"
	"github.com/ffigari/stored-strings/internal/config"
)

type clockI interface {
	Now() time.Time
}

type Authenticator struct {
	hmacSampleSecret []byte
	clock            clockI
}

func NewFromConfig() (*Authenticator, error) {
	config, err := config.Get()
	if err != nil {
		return nil, err
	}

	return New([]byte(config.JWTSecret), clockP.New()), nil
}

func New(secret []byte, clock clockI) *Authenticator {
	return &Authenticator{
		hmacSampleSecret: secret,
		clock:            clock,
	}
}

// GenerateToken generates a jwt token meant for users' authentication. An empty
// string will be returned if an internal error occurred.
func (a *Authenticator) GenerateToken() string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": a.clock.Now().AddDate(0, 0, 3).Unix(),
		"jti": uuid.New(),
	})

	tokenString, err := t.SignedString(a.hmacSampleSecret)
	if err != nil {
		log.Println("[auth;error]", err)
		return ""
	}

	return tokenString
}

func (a *Authenticator) IsValidToken(encodedToken string) bool {
	_, err := jwt.Parse(encodedToken, func(
		token *jwt.Token,
	) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(
				"Unexpected signing method: %v",
				token.Header["alg"],
			)
		}

		return a.hmacSampleSecret, nil
	})
	if err != nil {
		log.Println("[auth]", err)
		return false
	}

	return true
}
