package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", app.home)
	mux.Get("/{note}", app.showNote)
	mux.Get("/favicon.ico", app.favicon)

	mux.Post("/", app.createNote)

	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
