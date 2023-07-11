package app

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/famusovsky/md/internal/htmltemplates"
	"github.com/famusovsky/md/internal/models"
	"github.com/famusovsky/md/pkg/translator"
	"github.com/go-chi/chi"
	"github.com/shurcooL/github_flavored_markdown"
)

// home - обработчик главной страницы.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.html", &htmltemplates.Data{})
}

// showNote - обработчик страницы заметки.
func (app *application) showNote(w http.ResponseWriter, r *http.Request) {
	encodedId := chi.URLParam(r, "note")

	app.infoLog.Println("Got note request for: " + encodedId)

	id, err := translator.Translate(encodedId)

	if err != nil {
		app.serverError(w, err)
		return
	}

	note, err := app.notesModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.home(w, r)
		} else {
			app.serverError(w, err)
		}
		return
	}

	unsafe := github_flavored_markdown.Markdown([]byte(note.Content))
	rendered := string(unsafe)

	app.render(w, r, "note.page.html", &htmltemplates.Data{Note: note, RenderedNote: rendered})
}

// createNote - обработчик создания заметки.
func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("text")
	d, err := strconv.Atoi(r.FormValue("days"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	id, err := app.notesModel.Add(t, d)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("new note with id = %d successfully created\n", id)

	encryptedId := translator.Encrypt(id)

	http.Redirect(w, &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Path: "/note", RawQuery: "id=" + strconv.Itoa(id)},
	}, "/"+encryptedId, http.StatusSeeOther)
}

// favicon - заглушка обработчика запроса на иконку.
func (app *application) favicon(w http.ResponseWriter, r *http.Request) {
	app.notFound(w)
}
