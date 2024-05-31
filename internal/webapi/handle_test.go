package webapi_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/dbpool"
	"github.com/ffigari/stored-strings/internal/postgresql"
	"github.com/ffigari/stored-strings/internal/webapi"
)

func TestF(t *testing.T) {
	testDBName := "test_db"

	require.NoError(t, postgresql.CreateEmptyDB(testDBName))

	ctx := context.Background()
	dbPool, err := dbpool.NewFromConfig(ctx, testDBName)
	require.NoError(t, err)
	require.NotNil(t, dbPool)

	authenticator, err := auth.NewFromConfig()
	require.NoError(t, err)
	require.NotNil(t, authenticator)

	t.Run("with nothing", func(t *testing.T) {
		handler := webapi.NewHandle().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				require.Nil(t, conn)
				require.Equal(t, 0, len(params))
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "ok")
			})
		require.NotNil(t, handler)

		req, err := http.NewRequest("POST", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "ok", rr.Body.String())
	})

	t.Run("with storage conn", func(t *testing.T) {
		handler := webapi.NewHandle().
			WithStorageConn().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				require.NotNil(t, conn)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "ok")
			})
		require.NotNil(t, handler)

		req, err := http.NewRequest("POST", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "ok", rr.Body.String())
	})

	t.Run("with params sent", func(t *testing.T) {
		handler := webapi.NewHandle().
			WithParams([]string{"p1", "p2"}).
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				require.Nil(t, nil, conn)
				require.Equal(t, "a", params["p1"])
				require.Equal(t, "b", params["p2"])
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "ok")
			})
		require.NotNil(t, handler)

		req, err := http.NewRequest(
			"POST",
			"/",
			bytes.NewBuffer([]byte("p1=a&p2=b")),
		)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "ok", rr.Body.String())
	})

	t.Run("with params not sent", func(t *testing.T) {
		handler := webapi.NewHandle().
			WithParams([]string{"p1", "p2"}).
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "ok")
			})
		require.NotNil(t, handler)

		req, err := http.NewRequest("POST", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.NotEqual(t, "ok", rr.Body.String())
	})

	t.Run("Authed without credentials", func(t *testing.T) {
		handler := webapi.NewHandle().Authed().Finish(authenticator, dbPool, func(
			ctx context.Context,
			w http.ResponseWriter,
			r *http.Request,
			conn *pgxpool.Conn,
			params map[string]string,
		) {
			require.Nil(t, nil, conn)
			require.Equal(t, 0, len(params))

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ok")
		})
		require.NotNil(t, handler)

		req, err := http.NewRequest("POST", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.NotEqual(t, "ok", rr.Body.String())
	})

	t.Run("Authed with credentials", func(t *testing.T) {
		handler := webapi.NewHandle().Authed().Finish(authenticator, dbPool, func(
			ctx context.Context,
			w http.ResponseWriter,
			r *http.Request,
			conn *pgxpool.Conn,
			params map[string]string,
		) {
			require.Nil(t, nil, conn)
			require.Equal(t, 0, len(params))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ok")
		})
		require.NotNil(t, handler)

		req, err := http.NewRequest("POST", "/", nil)
		require.NoError(t, err)
		req.Header.Set(
			"Cookie",
			fmt.Sprintf("authorization=%s", authenticator.GenerateToken()),
		)

		rr := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "ok", rr.Body.String())
	})
}
