package webapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/parse"
)

type handle struct {
	authentication bool
	connn          bool
	params         []string
}

func NewHandle() *handle {
	return &handle{
		authentication: false,
		connn:          false,
		params:         []string{},
	}
}

func (x *handle) WithStorageConn() *handle {
	x.connn = true
	return x
}

func (x *handle) WithParams(ps []string) *handle {
	x.params = ps
	return x
}

func (x *handle) Authed() *handle {
	x.authentication = true
	return x
}

func (x *handle) Finish(
	authenticator *auth.Authenticator,
	dbPool *pgxpool.Pool,
	cb func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		conn *pgxpool.Conn,
		params map[string]string,
	),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = context.Background()
			params = map[string]string{}
			conn   *pgxpool.Conn
		)

		if x.authentication {
			authorizationTokenCookie, err := r.Cookie("authorization")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf(
					"[webapi] missing authorization cookie at '%s'", r.URL.Path,
				)
				return
			}

			if !authenticator.IsValidToken(authorizationTokenCookie.Value) {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("[webapi] received invalid token at '%s'\n", r.URL.Path)
				return
			}
		}

		if len(x.params) != 0 {
			receivedParams := parse.BodyParams(r)

			missingParams := []string{}
			for _, paramName := range x.params {
				if _, ok := receivedParams[paramName]; !ok {
					missingParams = append(missingParams, paramName)
				}
			}

			if len(missingParams) != 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, fmt.Sprintf(
					"missing required body params [%s]",
					strings.Join(missingParams, ", "),
				))

				return
			}

			emptyParams := []string{}
			incorrectlyEncodedParams := []string{}

			for _, paramName := range x.params {
				trimmedParamValue := strings.TrimSpace(receivedParams[paramName])

				if trimmedParamValue == "" {
					emptyParams = append(emptyParams, paramName)
					continue
				}

				escapedParamValue, err := url.QueryUnescape(trimmedParamValue)
				if err != nil {
					incorrectlyEncodedParams = append(
						incorrectlyEncodedParams,
						fmt.Sprintf("%s %s", paramName, trimmedParamValue),
					)
					continue
				}

				params[paramName] = escapedParamValue
			}

			if len(emptyParams) != 0 || len(incorrectlyEncodedParams) != 0 {
				if len(emptyParams) != 0 {
					log.Printf(
						"[webapi] found empty value for required body parameters [%s]",
						strings.Join(emptyParams, ", "),
					)
				}

				if len(incorrectlyEncodedParams) != 0 {
					log.Printf(
						"[webapi] could not decode required body parameters [%s]",
						strings.Join(incorrectlyEncodedParams, ", "),
					)
				}

				w.WriteHeader(http.StatusBadRequest)

				return
			}
		}

		if x.connn {
			if c, err := dbPool.Acquire(ctx); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf(
					"[webapi] acquiring connection: %s", err.Error(),
				)
			} else {
				conn = c
				defer conn.Release()
			}
		}

		cb(ctx, w, r, conn, params)
	}
}
