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
	UserID:      "03d70529-132e-4a90-ac1c-a24488fd20c5",
	Remind:      100,
}

func TestStorage(t *testing.T) {
	t.Run("get events on day", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		_, err := s.Create(generateEvent())
		require.NoError(t, err)

		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		event.Start = time.Now().Add(-24 * time.Hour)
		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		startOfDay := time.Now().Truncate(24 * time.Hour)
		eventsOnDay, err := s.GetEventsByDates(startOfDay, startOfDay.Add(24*time.Hour))
		require.NoError(t, err)
		require.Len(t, eventsOnDay, 2)
	})

	t.Run("get events on week", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		_, err := s.Create(generateEvent())
		require.NoError(t, err)

		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		event.Start = time.Now().Add(7 * (-24) * time.Hour)
		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		startOfWeek := time.Now().Truncate(7 * 24 * time.Hour)
		eventsOnWeek, err := s.GetEventsByDates(startOfWeek, startOfWeek.AddDate(0, 0, 7))
		require.NoError(t, err)
		require.Len(t, eventsOnWeek, 3)
	})

	t.Run("get events on month", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		_, err := s.Create(generateEvent())
		require.NoError(t, err)

		now := time.Now()
		event.Start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		event.Start = time.Now().AddDate(0, -1, 1)
		_, err = s.Create(generateEvent())
		require.NoError(t, err)

		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		eventsOnMonth, err := s.GetEventsByDates(startOfMonth, startOfMonth.AddDate(0, 1, 0))
		require.NoError(t, err)
		require.Len(t, eventsOnMonth, 2)
	})

	t.Run("get event by id", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(generateEvent())
		require.NoError(t, err)

		createdEvent, err := s.GetByID(ID)
		require.NoError(t, err)
		require.Equal(t, ID, createdEvent.ID)
		require.Equal(t, event.Title, createdEvent.Title)
		require.Equal(t, event.Start, createdEvent.Start)
		require.Equal(t, event.Finish, createdEvent.Finish)
		require.Equal(t, event.Description, createdEvent.Description)
		require.Equal(t, event.UserID, createdEvent.UserID)
		require.Equal(t, event.Remind, createdEvent.Remind)
	})

	t.Run("create event", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(generateEvent())
		require.NoError(t, err)

		createdEvent, err := s.GetByID(ID)
		require.NoError(t, err)
		require.Equal(t, ID, createdEvent.ID)
	})

	t.Run("update event", func(t *testing.T) {
		newTitle := "new title"
		newUserID := "0c70fc6c-305b-4409-89d5-8bea95d11af7"
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(generateEvent())
		require.NoError(t, err)

		createdEvent, err := s.GetByID(ID)
		require.NoError(t, err)
		createdEvent.Title = newTitle
		createdEvent.UserID = newUserID

		err = s.Update(createdEvent)
		require.NoError(t, err)

		updatedEvent, err := s.GetByID(ID)
		require.NoError(t, err)
		require.Equal(t, newTitle, updatedEvent.Title)
		require.Equal(t, newUserID, updatedEvent.UserID)
	})

	t.Run("delete event", func(t *testing.T) {
		s := New()
		connect := s.Connect(context.Background())
		require.NoError(t, connect)

		ID, err := s.Create(generateEvent())
		require.NoError(t, err)

		createdEvent, err := s.GetByID(ID)
		require.NoError(t, err)

		err = s.Delete(createdEvent)
		require.NoError(t, err)

		deletedEvent, err := s.GetByID(ID)
		require.Nil(t, deletedEvent)
		require.EqualError(t, err, "event not found in storage")
	})
}

func generateEvent() *model.Event {
	return &model.Event{
		Title:       event.Title,
		Start:       event.Start,
		Finish:      event.Finish,
		Description: event.Description,
		UserID:      event.UserID,
		Remind:      event.Remind,
	}
}
