package main

// TODO implement mobile version

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
	// Получение флагов addr и dsn командной строки
	addr := flag.String("addr", ":8080", "HTTP address")
	dsn := flag.String("dsn", "port=5432 user=postgres password=qwerty dbname=MD sslmode=disable", "PostgreSQL input string")
	flag.Parse()

	// Создание логгеров
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	// Создание подключения к БД
	dBase, err := db.OpenViaDsn(*dsn)
	if err != nil {
		dBase, err = db.OpenViaEnvVars()

		if err != nil {
			errorLog.Fatal(err)
		}
	}
	defer dBase.Close()

	// Создание кэша шаблонов
	templatesCache, err := htmltemplates.CreateNewCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Создание модели базы данных заметок
	notesModel, err := postgres.GetNotesModel(dBase)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Создание модели приложения и его запуск
	appModel := app.CreateModel(*addr, infoLog, errorLog, notesModel, templatesCache)
	appModel.Run()
}
