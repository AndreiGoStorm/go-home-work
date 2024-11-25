package model

import (
	"time"
)

type Notification struct {
	ID     string
	Title  string
	Start  time.Time
	UserID string
}

func GetNotificationFromEvent(event *Event) *Notification {
	return &Notification{
		ID:     event.ID,
		Title:  event.Title,
		Start:  event.Start,
		UserID: event.UserID,
	}
}
