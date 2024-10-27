package memorystorage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

var event = model.Event{
	Title:       "title memory event",
	Start:       time.Now(),
	Finish:      time.Now(),
	Description: sql.NullString{String: "description for event", Valid: true},
	UserID:      1,
}

func TestStorage(t *testing.T) {
	t.Run("find event by id", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(event)
		require.NoError(t, err)

		createdEvent, err := s.FindByID(ID)
		require.NoError(t, err)
		require.Equal(t, ID, createdEvent.ID)
		require.Equal(t, event.Title, createdEvent.Title)
		require.Equal(t, event.Start, createdEvent.Start)
		require.Equal(t, event.Finish, createdEvent.Finish)
		require.Equal(t, event.Description, createdEvent.Description)
		require.Equal(t, event.UserID, createdEvent.UserID)
	})

	t.Run("find all events", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		_, err := s.Create(event)
		require.NoError(t, err)

		_, err = s.Create(event)
		require.NoError(t, err)

		events, _ := s.FindAll()
		require.Len(t, events, 2)
	})

	t.Run("create event", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		_, err := s.Create(event)
		require.NoError(t, err)

		events, _ := s.FindAll()
		require.Len(t, events, 1)
	})

	t.Run("update event", func(t *testing.T) {
		newTitle := "new title"
		newUserID := 25
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(event)
		require.NoError(t, err)

		createdEvent, err := s.FindByID(ID)
		require.NoError(t, err)
		createdEvent.Title = newTitle
		createdEvent.UserID = newUserID

		err = s.Update(*createdEvent)
		require.NoError(t, err)

		updatedEvent, err := s.FindByID(ID)
		require.NoError(t, err)
		require.Equal(t, newTitle, updatedEvent.Title)
		require.Equal(t, newUserID, updatedEvent.UserID)
	})

	t.Run("delete event", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(event)
		require.NoError(t, err)

		createdEvent, err := s.FindByID(ID)
		require.NoError(t, err)

		err = s.Delete(*createdEvent)
		require.NoError(t, err)

		events, _ := s.FindAll()
		require.Len(t, events, 0)
	})
}
