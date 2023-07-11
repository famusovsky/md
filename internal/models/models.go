// Пакет models содержит определения данных для хранения в базе данных.
package models

import (
	"errors"
	"time"
)

// ErrNoRecord - ошибка, которая возникает, когда запись не найдена.
var ErrNoRecord = errors.New("models: record is not found")

// Note - структура заметки.
// ID - идентификатор заметки,
// Content - содержимое заметки,
// Created - дата создания заметки,
// Expires - дата истечения срока действия заметки.
type Note struct {
	ID      int
	Content string
	Created time.Time
	Expires time.Time
}
