package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Application) BindRoutes() http.Handler {
	api.Router.Use(
		middleware.RequestID,
		middleware.Recoverer,
		middleware.Logger,
	)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/tasks", api.handleCreateTask)
			// r.Get("/tasks", api.handleListTasks)
			// r.Get("/tasks/{id}", api.handleGetTask)
			// r.Put("/tasks/{id}", api.handleUpdateTask)
		})
	})

	return api.Router
}
