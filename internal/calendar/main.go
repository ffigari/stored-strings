package calendar

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/ui"
	"github.com/ffigari/stored-strings/internal/webapi"
)

func ForEach_old(
	ctx context.Context, conn *pgxpool.Conn, cb func(date, event string),
) error {
	rows, err := conn.Query(ctx, `SELECT date, event FROM calendar;`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var date, event string

		if err := rows.Scan(&date, &event); err != nil {
			return err
		}

		cb(date, event)
	}

	return nil
}

func ForEach(
	ctx context.Context, conn *pgxpool.Conn, cb func(time.Time, string),
) error {
	rows, err := conn.Query(ctx, `SELECT starts_at, description FROM events;`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var (
			startsAt time.Time
			description string
		)

		if err := rows.Scan(&startsAt, &description); err != nil {
			return err
		}

		cb(startsAt, description)
	}

	return nil
}

func AttachTo(
	r *mux.Router, baseRouterPrefix string, dbPool *pgxpool.Pool,
	authenticator *auth.Authenticator,
) error {
	cr := r.PathPrefix("/calendar").Subrouter()

	buenosAiresLocation, err := time.LoadLocation(
		"America/Argentina/Buenos_Aires",
	)
	if err != nil {
		return fmt.Errorf("loading location: %w", err)
	}

	webapi.At("").Of(cr).Serve(map[string]func(
		http.ResponseWriter, *http.Request,
	){
		"GET": webapi.NewHandle().
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
				if err := ForEach_old(ctx, conn, func(date, event string) {
					baseHTML += fmt.Sprintf(
						`<li class="list-group-item">%s: %s</li>`,
						date,
						event,
					)
				}); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("[webapi] failed iterating old: %s", err.Error())
					return
				}
				if err := ForEach(ctx, conn, func(
					startsAt time.Time, description string,
				) {
					adjustedStartsAt := startsAt.In(buenosAiresLocation)
					y, m, d := adjustedStartsAt.Date()
					baseHTML += fmt.Sprintf(
						`<li class="list-group-item">%s: %s</li>`,
						fmt.Sprintf(
							"%d-%02d-%02d %02d:%02d",
							y, m, d,
							adjustedStartsAt.Hour(), adjustedStartsAt.Minute(),
						),
						description,
					)
				}); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("[webapi] failed iterating: %s", err.Error())
					return
				}
				baseHTML += "</ul>"

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				wrappedHTML := ui.HTMLHeader(baseHTML)
				fmt.Fprint(w, wrappedHTML)
			}),
	})

	webapi.At("/events").Of(cr).Serve(map[string]func(
		http.ResponseWriter, *http.Request,
	) {
		"GET": webapi.NewHandle().
			Authed().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				fmt.Fprint(w, ui.HTMLHeader(ui.Form("Create", []string{
					ui.LabeledInput("Fecha", `
						<input
							id="date-input"
							class="form-control"
							type="datetime-local"
							name="starts_at"
						>
					`),
					ui.LabeledInput("Descripción", `
						<input
							id="description-input"
							class="form-control"
							type="text"
							name="description"
						/>
					`),
				})))
			}),
		"POST": webapi.NewHandle().
			Authed().
			WithParams([]string{"starts_at", "description"}).
			WithStorageConn().
			Finish(authenticator, dbPool, func(
				ctx context.Context,
				w http.ResponseWriter,
				r *http.Request,
				conn *pgxpool.Conn,
				params map[string]string,
			) {
				startsAt, err := webapi.ParseFormDatetime(
					params["starts_at"], buenosAiresLocation,
				)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("[webapi] parsing starts_at: %s", err.Error())
					return
				}

				if _, err := conn.Exec(ctx, `
					INSERT INTO events (starts_at, description)
					VALUES ($1, $2);
				`, startsAt, params["description"]); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("[webapi] failed scanning: %s", err.Error())
					return
				}

				http.Redirect(w, r, baseRouterPrefix+"/calendar", http.StatusSeeOther)
			}),
	})

	return nil
}
