package app

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/famusovsky/md/internal/models/postgres"
)

type Model struct {
	addr          string
	infoLog       *log.Logger
	errorLog      *log.Logger
	notesModel    *postgres.NotesModel
	templateCache map[string]*template.Template
}

type application struct {
	*Model
}

func CreateModel(addr string, infoLog *log.Logger, errorLog *log.Logger,
	notesModel *postgres.NotesModel, templateCache map[string]*template.Template) *Model {
	return &Model{
		addr:          addr,
		infoLog:       infoLog,
		errorLog:      errorLog,
		notesModel:    notesModel,
		templateCache: templateCache,
	}
}

func (model *Model) Run() {
	app := &application{
		Model: model,
	}

	go func(app *application) {
		for {
			err := app.notesModel.Tidy()
			if err != nil {
				app.errorLog.Println(err)
			} else {
				app.infoLog.Println("data get tidied successfully")
			}

			time.Sleep(12 * time.Hour)
		}
	}(app)

	srvr := &http.Server{
		Addr:     app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	// app.infoLog.Printf("Start server on http://127.0.0.1%s\n", app.addr)

	err := srvr.ListenAndServe()

	app.errorLog.Fatal(err)
}
