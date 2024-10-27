package memorystorage

import (
	"context"
	"sync"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type Storage struct {
	ctx    context.Context
	events map[string]model.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) FindAll() ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]model.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) FindByID(id string) (*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	if !ok {
		return nil, model.ErrEventNotFound
	}
	return &event, nil
}

func (s *Storage) Create(event model.Event) (string, error) {
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

func (s *Storage) Update(event model.Event) error {
	_, err := s.FindByID(event.ID)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()

	return nil
}

func (s *Storage) Delete(event model.Event) error {
	s.mu.Lock()
	delete(s.events, event.ID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) Connect(ctx context.Context) error {
	s.ctx = ctx
	s.events = make(map[string]model.Event, 50)
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = nil
	s.ctx = nil
	return nil
}
