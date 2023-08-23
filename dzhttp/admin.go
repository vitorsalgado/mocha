package dzhttp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/dzstd"
)

type Admin struct {
	app  *HTTPMockApp
	repo dzstd.MockRepository[*HTTPMock]
}

func (a *Admin) Init() http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/mocks", func(r chi.Router) {
		r.Get("/{id}", a.get)
		r.Post("/", a.add)
	})

	return r
}

func (a *Admin) add(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	mock, err := buildMockFromBytes(a.app, Request(), b, "json", false)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = a.repo.Save(r.Context(), mock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *Admin) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mock, err := a.repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if mock == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Add(httpval.HeaderContentType, httpval.MIMEApplicationCharsetUTF8)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(mock.Describe())
}
