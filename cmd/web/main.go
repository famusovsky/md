package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP address")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srvr := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start server on http://127.0.0.1%s\n", *addr)

	err := srvr.ListenAndServe()

	errorLog.Fatal(err)
}

// TODO handle database
