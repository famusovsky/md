package postgres

import (
	"database/sql"
	"errors"

	"github.com/famusovsky/md/pkg/models"
)

type NoteModel struct {
	db *sql.DB
}

func GetNoteModel(db *sql.DB) *NoteModel {
	return &NoteModel{db}
}

func (m *NoteModel) Get(id int) (*models.Note, error) {
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

func (m *NoteModel) Add(text string, daysTillExpire int) (int, error) {
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
