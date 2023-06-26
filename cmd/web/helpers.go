package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(msg)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	temp, err := ts.Clone()
	if err != nil {
		app.serverError(w, err)
	}

	tmpl := "{{define \"renderedNote\"}}\n" + td.RenderedNote + "\n{{end}}"

	temp.Parse(tmpl)

	err = temp.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}
