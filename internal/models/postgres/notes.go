// Пакет postgres реализует взаимодействие с базой данных PostgreSQL, хранящей данные из пакета models.
package postgres

import (
	"database/sql"
	"errors"
	"log"

	"github.com/famusovsky/md/internal/models"
)

// NotesModel - модель базы данных заметок.
type NotesModel struct {
	db *sql.DB
}

// GetNotesModel - создание модели базы данных заметок.
// Принимает базу данных.
// Возвращает модель базы данных заметок и ошибку.
func GetNotesModel(db *sql.DB) (*NotesModel, error) {
	err := checkDB(db)
	if err != nil {
		return nil, err
	}

	return &NotesModel{db}, nil
}

// createDB - создание базы данных заметок.
func createDB(db *sql.DB) error {
	q := `CREATE TABLE notes (
        id SERIAL NOT NULL PRIMARY KEY, 
		content TEXT NOT NULL,
		created TIMESTAMP NOT NULL,
		expires TIMESTAMP NOT NULL
    );`

	_, err := db.Exec(q)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// checkDB - проверка базы данных заметок.
func checkDB(db *sql.DB) error {
	q :=
		`SELECT COUNT(*) = 4 AS proper
			FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = 'notes'
			AND (
				(column_name = 'id' AND data_type = 'integer')
				OR (column_name = 'content' AND data_type = 'text')
				OR (column_name = 'created' AND data_type = 'timestamp without time zone')
				OR (column_name = 'expires' AND data_type = 'timestamp without time zone')
			);`
	var proper bool
	db.QueryRow(q).Scan(&proper)

	if !proper {
		// TODO make normal migration
		q = `DROP TABLE IF EXISTS notes`
		_, err := db.Exec(q)
		if err != nil {
			return errors.New("cannot drop incorrect 'notes' table in the database")
		}

		err = createDB(db)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get - получение заметки по id.
// Принимает id заметки.
// Возвращает модель заметки и ошибку.
func (m *NotesModel) Get(id int) (*models.Note, error) {
	q :=
		`SELECT Content, Created, Expires FROM notes 
		WHERE id = $1 AND expires > CURRENT_TIMESTAMP;`

	note := &models.Note{ID: id}

	err := m.db.QueryRow(q, id).Scan(&note.Content, &note.Created, &note.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return note, nil
}

// Add - добавление заметки.
// Принимает текст заметки и количество дней до истечения срока действия заметки.
// Возвращает id заметки и ошибку.
func (m *NotesModel) Add(text string, daysTillExpire int) (int, error) {
	q :=
		`INSERT INTO notes (content, created, expires) VALUES (
		$1,
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP + INTERVAL '1 DAY' * $2
	) RETURNING id`

	var id int
	err := m.db.QueryRow(q, text, daysTillExpire).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// delete - удаление заметки по id.
func (m *NotesModel) delete(id int) error {
	q :=
		`DELETE FROM notes
		WHERE id = $1;`

	_, err := m.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

// Tidy - удаление просроченных заметок.
// Возвращает ошибку.
func (m *NotesModel) Tidy() error {
	q :=
		`SELECT id FROM notes
		WHERE expires < CURRENT_TIMESTAMP;
	`
	rows, err := m.db.Query(q)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		err = m.delete(id)
		if err != nil {
			return err
		}
	}

	return nil
}
