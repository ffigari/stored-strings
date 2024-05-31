package webapi_test


import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/webapi"
)

func (s *AtOfServeSuite) TestMethodNotAllowed() {
	r := mux.NewRouter()

	webapi.
		At("/home").
		Of(r).
		Serve(map[string]func(http.ResponseWriter, *http.Request){})

	server := httptest.NewServer(r)

	for _, method := range []string{
		"GET", "POST", "PUT", "PATCH",
	} {
		req, err := http.NewRequest(method, server.URL + "/home", nil)
		s.Require().NoError(err)

		res, err := (&http.Client{}).Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func (s *AtOfServeSuite) TestGET() {
	r := mux.NewRouter()

	webapi.
		At("/home").
		Of(r).
		Serve(map[string]func(http.ResponseWriter, *http.Request){
			"GET": func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		})

	server := httptest.NewServer(r)

	req, err := http.NewRequest("GET", server.URL + "/home", nil)
	s.Require().NoError(err)

	res, err := (&http.Client{}).Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

func (s *AtOfServeSuite) TestGETAndPOST() {
	r := mux.NewRouter()

	webapi.
		At("/home").
		Of(r).
		Serve(map[string]func(http.ResponseWriter, *http.Request){
			"GET": func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			"POST": func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
		})

	server := httptest.NewServer(r)

	req, err := http.NewRequest("GET", server.URL + "/home", nil)
	s.Require().NoError(err)

	res, err := (&http.Client{}).Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)

	req, err = http.NewRequest("POST", server.URL + "/home", nil)
	s.Require().NoError(err)

	res, err = (&http.Client{}).Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, res.StatusCode)
}

func (s *AtOfServeSuite) TestPOST() {
	r := mux.NewRouter()

	webapi.
		At("/home").
		Of(r).
		Serve(map[string]func(http.ResponseWriter, *http.Request){
			"POST": func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		})

	server := httptest.NewServer(r)

	req, err := http.NewRequest("POST", server.URL + "/home", nil)
	s.Require().NoError(err)

	res, err := (&http.Client{}).Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

func (s *AtOfServeSuite) TestPUT() {
	r := mux.NewRouter()

	webapi.
		At("/home").
		Of(r).
		Serve(map[string]func(http.ResponseWriter, *http.Request){
			"PUT": func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		})

	server := httptest.NewServer(r)

	req, err := http.NewRequest("PUT", server.URL + "/home", nil)
	s.Require().NoError(err)

	res, err := (&http.Client{}).Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

func (s *AtOfServeSuite) TestPATCH() {
	r := mux.NewRouter()

	webapi.
		At("/home").
		Of(r).
		Serve(map[string]func(http.ResponseWriter, *http.Request){
			"PATCH": func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		})

	server := httptest.NewServer(r)

	req, err := http.NewRequest("PATCH", server.URL + "/home", nil)
	s.Require().NoError(err)

	res, err := (&http.Client{}).Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

type AtOfServeSuite struct {
	suite.Suite
}

func TestAtOfServe(t *testing.T) {
	suite.Run(t, new(AtOfServeSuite))
}
