package api

import "github.com/go-chi/chi/v5"

type Application struct {
	Router *chi.Mux
}
