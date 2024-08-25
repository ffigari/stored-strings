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
	"github.com/ffigari/stored-strings/internal/interactions"
	"github.com/ffigari/stored-strings/internal/oos"
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

	r.HandleFunc("/canvas", func(w http.ResponseWriter, r *http.Request) {
		ui.HTMLHeader(w, `
			<style>
			#canvas-container {
				width: 100%;
			}
			#my-canvas {
				border: 2px solid black;
				border-radius: 10px;
			}
			</style>
			<div id="canvas-container">
				<canvas id="my-canvas" ></canvas>
			</div>
			<script>
			const canvas = document.getElementById('my-canvas');
			const ctx = canvas.getContext('2d')

			function resizeCanvas() {
				const canvas = document.getElementById('my-canvas');
				const dpr = window.devicePixelRatio || 1;
				const padding = 10; // Optional: adjust canvas padding
				const rect = canvas.parentElement.getBoundingClientRect();

				h = window.innerHeight - 2 * rect.top

				// Set canvas display size
				canvas.style.width = `+"`"+`${rect.width}px`+"`"+`;
				canvas.style.height = `+"`"+`${h}px`+"`"+`;

				// Set canvas buffer size
				canvas.width = rect.width * dpr;
				canvas.height = h * dpr;
			}
			resizeCanvas();
			window.addEventListener('resize', resizeCanvas);

			let cs = []

			Promise.resolve().then(async () => {
				try {
					const res = await fetch("/interactions/clicks", {
						method: "GET",
					})
					if (res.ok) {
						const l = await res.json()
						cs.push(...l)
					}
				} catch (e) {
					console.log("fail get", e)
				}
			})

			canvas.addEventListener('mousedown', async (e) => {
				const r = canvas.getBoundingClientRect();

				const c = {
					x: e.clientX - r.left,
					y: e.clientY - r.top
				}

				cs.push(c)

				try {
					const res = await fetch("/interactions/clicks", {
						method: "POST",
						body: JSON.stringify(c),
					})
				} catch (e) {
					console.log("error", e)
				}
			})

			animate = (t) => {
				const { width: w, height: h } = canvas

				ctx.clearRect(0, 0, w, h)

				cs.forEach((c) => {
					ctx.beginPath();
					ctx.arc(
						c.x, c.y,
						50 * ((Math.sin(t / 1000) + 1) / 2),
						0, 2 * Math.PI, false,
					);
					ctx.stroke();
				})

				requestAnimationFrame(animate)
			}

			animate()

			</script>
		`)
	})

	if err := calendar.AttachTo(r, "/i", dbPool, authenticator); err != nil {
		return nil, err
	}

	auth.AttachTo(r, password, dbPool, authenticator)

	interactions.AttachTo(rr, dbPool)

	rootPath, err := oos.GetRootPath()
	if err != nil {
		return nil, fmt.Errorf("getting root path: %w", err)
	}

	r.PathPrefix("/").Handler(
		http.StripPrefix("/i", http.FileServer(http.Dir(rootPath+"/i"))),
	)

	rr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/i", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	return rr, nil
}
