package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Get("/", app.home)
	mux.Get("/{note}", app.showNote)
	mux.Get("/favicon.ico", app.favicon)

	mux.Post("/", app.createNote)

	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
