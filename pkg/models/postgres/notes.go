package postgres

import (
	"database/sql"
	"errors"
	"log"

	"github.com/famusovsky/md/pkg/models"
)

type NotesModel struct {
	db *sql.DB
}

func GetNotesModel(db *sql.DB) (*NotesModel, error) {
	err := checkDB(db)
	if err != nil {
		return nil, err
	}

	return &NotesModel{db}, nil
}

func checkDB(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS notes (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );`

	_, err := db.Exec(q)
	if err != nil {
		log.Fatal(err)
	}

	q =
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
		return errors.New("incorrect 'notes' table in the database")
	}

	return nil
}

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
