package personalwebsite

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gorilla/mux"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/calendar"
	"github.com/ffigari/stored-strings/internal/dbpool"
	"github.com/ffigari/stored-strings/internal/ui"
)

// TODO: Borrar logs. Tratar el log como algo que puede crecer y ser un problema
func NewMux(
	ctx context.Context, dbName string, authenticator *auth.Authenticator,
	password string,
) (*mux.Router, error) {
	rr := mux.NewRouter()
	r := rr.PathPrefix("/i").Subrouter()


	// TODO: This db pool should be closed at graceful shutdown
	dbPool, err := dbpool.NewFromConfig(ctx, dbName)
	if err != nil {
		return nil, err
	}

	if _, filename, _, ok := runtime.Caller(0); !ok {
		return nil, fmt.Errorf("no caller information")
	} else {
		if homePageBytes, err := os.ReadFile(
			path.Dir(filename) + "/home.html",
		); err != nil {
			return nil, err
		} else {
			r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
				ui.HTMLHeader(w, string(homePageBytes))
			})
		}
	}

	r.Handle(
		"/",
		http.StripPrefix("/i/", http.FileServer(http.Dir("i"))),
	)

	for _, path := range []string{
		"/favicon.ico",
		"/status",
	} {
		rr.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}

	for _, p := range []string{
		"/retrocompatible",
		"/retrocompatibilidad",
	} {
		r.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			ui.HTMLHeader(w, ui.Paragraphs([]string{`
				Que lo nuevo siempre pueda existir.
			`, `
				Que lo viejo fluya a la par de lo nuevo.
			`, `
				Que lo eterno se haga presente.
			`, `
				Que el presente se haga eterno.
			`}))
		})
	}

	if err := calendar.AttachTo(r, "/i", dbPool, authenticator); err != nil {
		return nil, err
	}

	auth.AttachTo(r, password, dbPool, authenticator)

	rr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/i", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	return rr, nil
}
