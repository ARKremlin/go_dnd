package rest

import (
	"net/http"

	"github.com/ARKremlin/go_dnd/internal/transport/rest/handler"
	"github.com/ARKremlin/go_dnd/internal/transport/rest/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type RouterDeps struct {
	AuthHandler *handler.AuthHandler
	JWTSecret   []byte
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.AccessLog)
	r.Use(chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.AuthHandler.Register)
			r.Post("/login", deps.AuthHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.JWTSecret))
			r.Get("/me", deps.AuthHandler.Me)
		})
	})
	return r
}
