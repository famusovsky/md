package main

// TODO add tests

// TODO improve link shortaning, maibe make it more random

import (
	"flag"
	"log"
	"os"

	"github.com/famusovsky/md/internal/app"
	"github.com/famusovsky/md/internal/htmltemplates"
	"github.com/famusovsky/md/internal/models/postgres"
	"github.com/famusovsky/md/pkg/db"
	_ "github.com/lib/pq"
)

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

	templatesCache, err := htmltemplates.CreateNewTemplatesCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	notesModel, err := postgres.GetNotesModel(dBase)
	if err != nil {
		errorLog.Fatal(err)
	}

	appModel := app.CreateModel(*addr, infoLog, errorLog, notesModel, templatesCache)

	appModel.Run()
}
