package models

import "time"

type note struct {
	ID      int64
	Text    string
	Created time.Time
	Expired time.Time
}
