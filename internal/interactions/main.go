package interactions

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gorilla/mux"
)

type click struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func AttachTo(r *mux.Router, dbPool *pgxpool.Pool) {
	persistClick := func(
		w http.ResponseWriter, r *http.Request,
	) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		var c click
		if err := json.Unmarshal(body, &c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

		conn, err := dbPool.Acquire(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer conn.Release()

		if _, err := conn.Exec(ctx, `
			INSERT INTO clicks (x, y)
			VALUES ($1, $2)
		;`, c.X, c.Y); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

	getClicks := func(
		w http.ResponseWriter, r *http.Request,
	) {
		ctx := context.Background()

		conn, err := dbPool.Acquire(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer conn.Release()

		rows, err := conn.Query(ctx, `
			SELECT x, y FROM clicks;
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var cs []click

		for rows.Next() {
			var x, y int

			if err := rows.Scan(&x, &y); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cs = append(cs, click{
				X: x,
				Y: y,
			})
		}

		serializedCs, err := json.Marshal(cs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(serializedCs))
	}

	r.HandleFunc("/interactions/clicks", func(
		w http.ResponseWriter, r *http.Request,
	) {
		switch r.Method {
		case http.MethodGet:
			getClicks(w, r)
		case http.MethodPost:
			persistClick(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
