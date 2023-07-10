package main

// TODO add tests

// TODO improve link shortaning, maibe make it more random

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/famusovsky/md/pkg/db"
	"github.com/famusovsky/md/pkg/models/postgres"
	_ "github.com/lib/pq"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	notesModel    *postgres.NotesModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	dsn := flag.String("dsn", "port=5432 user=postgres password=qwerty dbname=MD sslmode=disable", "PostgreSQL input string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	dBase, err := db.OpenDB(*dsn)
	if err != nil {
		dBase, err = db.OpenDB("")

		if err != nil {
			errorLog.Fatal(err)
		}
	}
	defer dBase.Close()

	cache, err := createNewTemplatesCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Println("templates cache is created")

	model, err := postgres.GetNotesModel(dBase)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		notesModel:    model,
		templateCache: cache,
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
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start server on http://127.0.0.1%s\n", *addr)

	err = srvr.ListenAndServe()

	errorLog.Fatal(err)
}
