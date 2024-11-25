package model

import (
	"database/sql"
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found in storage")

type Event struct {
	ID          string
	Title       string
	Start       time.Time
	Finish      time.Time
	Description sql.NullString
	UserID      string
	Remind      int32
	RemindDate  time.Time
}
