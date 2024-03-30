package personalwebsite

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/calendar"
	"github.com/ffigari/stored-strings/internal/dbpool"
	"github.com/ffigari/stored-strings/internal/webapi"
)

type methods struct {
	Get  func(w http.ResponseWriter, r *http.Request)
	Post func(w http.ResponseWriter, r *http.Request)
}

func at(mux *http.ServeMux, path string, ms methods) {
	mux.HandleFunc(path, func(
		w http.ResponseWriter, r *http.Request,
	) {
		if ms.Get != nil && r.Method == "GET" {
			ms.Get(w, r)
			return
		}

		// TODO: trigger this unsafe pointer usage with test
		if r.Method == "POST" {
			ms.Post(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}

// TODO: the authentication should be kept in one single package. This includes
// being able to validate the token and being "attachable" to this mux. The web
// password should be kept in there
// TODO: Borrar logs. Tratar el log como algo que puede crecer y ser un problema
func NewMux(
	ctx context.Context, dbName string, authenticator *auth.Authenticator,
	password string,
) (*http.ServeMux, error) {
	// TODO: use gorilla mux instead
	mux := http.NewServeMux()

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
			mux.HandleFunc("/i", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, addHTMLHeader(string(homePageBytes)))
			})
		}
	}

	mux.Handle(
		"/i/",
		http.StripPrefix("/i/", http.FileServer(http.Dir("i"))),
	)

	for _, path := range []string{
		"/favicon.ico",
		"/status",
	} {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}

	at(mux, "/i/login", methods{
		Get: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, addHTMLHeader(`
				<form method="POST">
					<div class="mb-3">
						<label for="username-input" class="form-label">
							Username
						</label>
						<input
							id="username-input" class="form-control"
							type="text" name="username" required
						/>
					</div>
					<div class="mb-3">
						<label for="password-input" class="form-label">
							Password
						</label>
						<input
							id="password-input" class="form-control"
							type="password" name="password" required
						/>
					</div>
					<button class="btn btn-primary" type="submit">
						Login
					</button>
				</form>
			`))
			return
		},
		Post: webapi.NewHandle().
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
				fmt.Fprintf(w, addHTMLHeader("logged in"))
				log.Println("[webapi] successful login")
				return
			}),
	})

	at(mux, "/i/calendar", methods{
		Get: webapi.NewHandle().
			Authed().
			WithStorageConn().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {

				baseHTML := `<ul class="list-group">`
				if err := calendar.ForEach(ctx, conn, func(date, event string) {
					baseHTML += fmt.Sprintf(
						`<li class="list-group-item">%s: %s</li>`,
						date,
						event,
					)
				}); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("[webapi] failed iterating: %s", err.Error())
					return
				}
				baseHTML += "</ul>"

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				wrappedHTML := addHTMLHeader(baseHTML)
				fmt.Fprint(w, wrappedHTML)
			}),
	})

	at(mux, "/i/calendar/events", methods{
		Get: webapi.NewHandle().
			Authed().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				fmt.Fprint(w, addHTMLHeader(`
					<form method="POST">
						<div class="row mb-3">
							<label
								for="date-input"
								class="col-sm-2 col-form-label"
							>
								Fecha
							</label>
							<div class="col-sm-10">
								<input
									id="date-input"
									class="form-control"
									type="text"
									name="date"
								/>
							</div>
						</div>
						<div class="row mb-3">
							<label
								for="description-input"
								class="col-sm-2 col-form-label"
							>
								Descripción
							</label>
							<div class="col-sm-10">
								<input
									id="description-input"
									class="form-control"
									type="text"
									name="description"
								/>
							</div>
						</div>
						<button class="btn btn-primary" type='submit'>
							Create
						</button>
					</form>
				`))
			}),
		Post: webapi.NewHandle().
			Authed().
			WithParams([]string{"date", "description"}).
			WithStorageConn().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				if _, err := conn.Exec(ctx, `
					INSERT INTO calendar (date, event)
					VALUES ($1, $2);
				`, params["date"], params["description"]); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("[webapi] failed scanning: %s", err.Error())
					return
				}

				http.Redirect(w, r, "/i/calendar", http.StatusSeeOther)
			}),
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/i", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	return mux, nil
}

// TODO: esta funcion deberia recibir el response y además setearle el header de
// content type
func addHTMLHeader(baseHTML string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
				<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
				<meta name="viewport" content="width=device-width, initial-scale=1">
			</head>
			<body>
				<div class="container">
					<div class="my-3">
						%s
					</div>
				</div>
				<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
			</body>
		</html>
	`, baseHTML)
}
