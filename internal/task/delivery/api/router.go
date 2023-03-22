package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Post("/task", h.AddTask)
	r.Get("/task/{id}", h.GetTaskResult)

	r.NotFound(h.NotFoundHandler)
	r.MethodNotAllowed(h.MethodNotAllowedHandler)

	return r
}
