package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/auth"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/middleware"
)

func New(
	userHandler http.HandlerFunc,
	loginHandler http.HandlerFunc,
	meHandler http.HandlerFunc,
	conversationHandler http.HandlerFunc,
	messageHandler http.HandlerFunc,
	messageHistoryHandler http.HandlerFunc,
	tokenService *auth.TokenService,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	if userHandler != nil {
		r.Post("/users", userHandler)
	}
	if loginHandler != nil {
		r.Post("/login", loginHandler)
	}
	if meHandler != nil {
		r.With(middleware.RequireAuth(tokenService)).Get("/me", meHandler)
	}
	if conversationHandler != nil {
		r.With(middleware.RequireAuth(tokenService)).
			Post("/conversations", conversationHandler)
	}
	if messageHandler != nil {
		r.With(middleware.RequireAuth(tokenService)).
			Post("/conversations/{conversationID}/messages", messageHandler)
	}
	if messageHistoryHandler != nil {
		r.With(middleware.RequireAuth(tokenService)).
			Get("/conversations/{conversationID}/messages", messageHistoryHandler)
	}

	return r
}
