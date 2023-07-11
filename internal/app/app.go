// Пакет app реализует логику приложения.
package app

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/famusovsky/md/internal/models/postgres"
)

// Model - модель приложения.
type Model struct {
	addr          string
	infoLog       *log.Logger
	errorLog      *log.Logger
	notesModel    *postgres.NotesModel
	templateCache map[string]*template.Template
}

// application - структура, которая хранит в себе модель приложения.
type application struct {
	*Model
}

// CreateModel - создание модели приложения.
// Параметры:
// addr - адрес сервера,
// infoLog - логгер информационных сообщений,
// errorLog - логгер ошибок,
// notesModel - модель базы данных заметок,
// templateCache - кэш шаблонов.
// Возвращает модель приложения.
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

// Run - запуск приложения.
// Принимает модель приложения.
func (model *Model) Run() {
	// Создание приложения на основе модели.
	app := &application{
		Model: model,
	}

	// Запуск процесса очистки базы данных от устаревших данных.
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

	// Создание и запуск сервера.
	srvr := &http.Server{
		Addr:     app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	err := srvr.ListenAndServe()

	app.errorLog.Fatal(err)
}
