package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type Storage struct {
	events map[string]*model.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) GetEventsByDates(eventStart, eventFinish time.Time) ([]*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*model.Event, 0, 50)
	for _, event := range s.events {
		isDateInRange := event.Start.After(eventStart) && event.Start.Before(eventFinish)
		if isDateInRange || event.Start.Equal(eventStart) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) GetRemindEvents(start time.Time) ([]*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start = start.Truncate(24 * time.Hour)
	finish := start.AddDate(0, 0, 1)
	events := make([]*model.Event, 0, 50)
	for _, event := range s.events {
		isDateInRange := event.RemindDate.After(start) && event.RemindDate.Before(finish)
		if isDateInRange || event.RemindDate.Equal(start) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) GetByID(id string) (*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	if !ok {
		return nil, model.ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) Create(event *model.Event) (string, error) {
	eventUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	event.ID = eventUUID.String()

	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()
	return event.ID, nil
}

func (s *Storage) Update(event *model.Event) error {
	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) Delete(event *model.Event) error {
	s.mu.Lock()
	delete(s.events, event.ID)
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteOldEvents() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	start := time.Now().Truncate(24 * time.Hour)
	start = start.AddDate(-1, 0, 0)
	for _, event := range s.events {
		if event.RemindDate.Before(start) || event.RemindDate.Equal(start) {
			delete(s.events, event.ID)
		}
	}
	return nil
}

func (s *Storage) Connect(_ context.Context) error {
	s.events = make(map[string]*model.Event, 50)
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = nil
	return nil
}
