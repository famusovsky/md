package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/", app.createNote)
	mux.HandleFunc("/note", app.showNote)

	// TODO serve static files

	return mux
}
