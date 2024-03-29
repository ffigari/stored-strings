package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/clock"
	"github.com/ffigari/stored-strings/internal/config"
	"github.com/ffigari/stored-strings/internal/parse"
)

func main() {
	ctx := context.Background()

	config, err := config.Get()
	if err != nil {
		log.Fatal("[webapi] could not read config: ", err)
	}

	conn, err := pgx.Connect(
		ctx,
		config.PostgresServerConnectionString+"/storedstrings",
	)
	if err != nil {
		log.Fatal("[webapi] could not connect to db: ", err)
	}

	var (
		authenticator = auth.New([]byte(config.JWTSecret), clock.New())
	)

	for _, path := range []string{
		"/favicon.ico",
		"/status",
	} {
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}

	wrapHTML := func(baseHTML string) string {
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
			</html>`, baseHTML)
	}

	http.HandleFunc("/i/yo", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/i", http.StatusSeeOther)
	})

	if baseHTMLBytes, err := os.ReadFile("yo.html"); err != nil {
		log.Fatal(err)
	} else {
		wrappedHTML := wrapHTML(string(baseHTMLBytes))
		http.HandleFunc("/i", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-9")
			fmt.Fprint(w, wrappedHTML)
		})
	}

	http.HandleFunc("/i/calendario", func(
		w http.ResponseWriter, r *http.Request,
	) {
		requestCtx := context.Background()

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

		rows, err := conn.Query(requestCtx, `
			SELECT date, event
			FROM calendar;
		`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf(
				"[webapi] failed to get calendar events: %s", err.Error(),
			)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		baseHTML := `<ul class="list-group">`
		for rows.Next() {
			var date, event string

			if err := rows.Scan(&date, &event); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("[webapi] failed scanning: %s", err.Error())
				return
			}

			baseHTML += fmt.Sprintf(
				`<li class="list-group-item">%s: %s</li>`,
				date,
				event,
			)
		}
		baseHTML += "</ul>"
		wrappedHTML := wrapHTML(baseHTML)
		fmt.Fprint(w, wrappedHTML)
	})

	http.HandleFunc("/i/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(w, wrapHTML(`
				<form method="POST">
					<div class="mb-3">
						<label for="username-input" class="form-label">
							Username
						</label>
						<input
							id="username-input" class="form-label"
							type="text" name="username" required
						/>
					</div>
					<div class="mb-3">
						<label for="password-input" class="form-label">
							Password
						</label>
						<input
							id="password-input" class="form-label"
							type='password' name='password' required
						/>
					</div>
					<button class="btn btn-primary" type='submit'>
						Login
					</button>
				</form>
			`))
			return
		}

		if r.Method == "POST" {
			params := parse.BodyParams(r)

			username, okusername := params["username"]
			password, okpassword := params["password"]

			if !okusername || !okpassword {
				missingParams := []string{}

				if !okusername {
					missingParams = append(missingParams, "username")
				}

				if !okpassword {
					missingParams = append(missingParams, "password")
				}

				log.Printf(
					"[webapi] missing params in login request [%s]",
					strings.Join(missingParams, ", "),
				)
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			if username != "admin" || password != config.WebPassword {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf(
					"[webapi] failed login attempt with username '%s'",
					username,
				)
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
			fmt.Fprintf(w, wrapHTML("logged in"))
			log.Println("[webapi] successful login")
			return
		}
	})

	http.Handle(
		"/i/",
		http.StripPrefix("/i/", http.FileServer(http.Dir("i"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/i/yo", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	})

	log.Println("[webapi] about to start http server")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
