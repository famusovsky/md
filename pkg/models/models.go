package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: record is not found")

type Note struct {
	ID      int
	Content string
	Created time.Time
	Expires time.Time
}
