package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	// create router
	mux := chi.NewRouter()

	// setup middlewares
	mux.Use(middleware.Recoverer)
	mux.Get("/", app.HomePage)
	return mux
}
