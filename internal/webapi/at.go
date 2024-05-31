package webapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

type q1 struct {
	path string
}

type q2 struct {
	q1
	r *mux.Router
}

func At(path string) *q1 {
	return &q1{
		path: path,
	}
}

func (x *q1) Of(r *mux.Router) *q2 {
	return &q2{
		q1: q1{x.path},
		r: r,
	}
}

func (x *q2) Serve(
	handlers map[string]func(http.ResponseWriter, *http.Request),
) {
	x.r.HandleFunc(x.path, func(w http.ResponseWriter, r *http.Request) {
		for _, method := range []string{
			"GET", "POST", "PUT", "PATCH",
		} {
			if r.Method != method {
				continue
			}

			if h, ok := handlers[method]; ok {
				h(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
