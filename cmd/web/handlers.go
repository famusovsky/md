package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/famusovsky/md/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.createNote(w, r)
	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "home.page.html", &templateData{})

	app.infoLog.Println("home page successfully rendered")
}

func (app *application) showNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	note, err := app.notes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "note.page.html", &templateData{Note: note})

	app.infoLog.Printf("note page with id = %d successfully rendered\n", id)
}

func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	t := r.FormValue("text")
	d, err := strconv.Atoi(r.FormValue("days"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	id, err := app.notes.Add(t, d)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("new note with id = %d successfully created\n", id)

	http.Redirect(w, &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Path: "/note", RawQuery: "id=" + strconv.Itoa(id)},
	}, "/note?id="+strconv.Itoa(id), http.StatusSeeOther)
}
