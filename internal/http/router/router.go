package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(userHandler http.HandlerFunc) http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	if userHandler != nil {
		r.Post("/users", userHandler)
	}

	return r
}
