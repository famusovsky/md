package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/famusovsky/md/pkg/models/postgres"
	_ "github.com/lib/pq"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	notes         *postgres.NoteModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	dsn := flag.String("dsn", "port=5432 user=postgres password=qwerty dbname=MD sslmode=disable", "PostgreSQL input string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	cache, err := createNewTemplatesCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Println("templates cache is created")

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		notes:         postgres.GetNoteModel(db),
		templateCache: cache,
	}

	srvr := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start server on http://127.0.0.1%s\n", *addr)

	err = srvr.ListenAndServe()

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
