package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	clockP "github.com/ffigari/stored-strings/internal/clock"
	"github.com/ffigari/stored-strings/internal/config"
	"github.com/ffigari/stored-strings/internal/ui"
	"github.com/ffigari/stored-strings/internal/webapi"
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

func AttachTo(
	r *mux.Router, password string, dbPool *pgxpool.Pool,
	authenticator *Authenticator,
) {
	webapi.At("/login").Of(r).Serve(map[string]func(
		http.ResponseWriter, *http.Request,
	) {
		"GET": func(w http.ResponseWriter, r *http.Request) {
			ui.HTMLHeader(w, ui.Form("Login", []string{
				ui.LabeledInput("Username", `
					<input
						id="username-input"
						class="form-control"
						type="text"
						name="username"
						required
					/>
				`),
				ui.LabeledInput("Password", `
					<input
						id="password-input"
						class="form-control"
						type="password"
						name="password"
						required
					/>
				`),
			}))
		},
		"POST": webapi.NewHandle().
			WithParams([]string{"username", "password"}).
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				incorrectUsername := params["username"] != "admin"
				incorrectPassword := params["password"] != password
				if incorrectUsername || incorrectPassword {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintf(w, "invalid credentials")

					return
				}

				token := authenticator.GenerateToken()
				if token == "" {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("[error] failed to generate auth token")
					return
				}

				w.Header().Set(
					"Set-Cookie",
					fmt.Sprintf("authorization=%s; Secure; HttpOnly", token),
				)

				ui.HTMLHeader(w, "logged in")

				return
			}),
	})
}
