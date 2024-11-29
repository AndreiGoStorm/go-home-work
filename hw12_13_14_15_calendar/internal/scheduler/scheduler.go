package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	storage storage.Storage
	rabbit  *rabbit.Rabbit
	logg    *logger.Logger
}

func New(storage storage.Storage, rabbit *rabbit.Rabbit, logg *logger.Logger) *Scheduler {
	return &Scheduler{
		storage: storage,
		rabbit:  rabbit,
		logg:    logg,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.notify(ctx)
				s.clear()
			case <-ctx.Done():
				s.logg.Error("context done", ctx.Err())
				done <- true
				return
			}
		}
	}()

	<-done
	return nil
}

func (s *Scheduler) notify(ctx context.Context) {
	events, err := s.storage.GetRemindEvents(time.Now())
	if err != nil {
		s.logg.Error("scheduler notify", err)
		return
	}

	if len(events) == 0 {
		s.logg.Info("scheduler no events for notifying")
		return
	}

	if err := s.rabbit.Notify(ctx, events); err != nil {
		s.logg.Error("scheduler push", err)
	}
	s.logg.Info(fmt.Sprintf("scheduler notifying: %d", len(events)))
}

func (s *Scheduler) clear() {
	if err := s.storage.DeleteOldEvents(); err != nil {
		s.logg.Error("scheduler clear", err)
		return
	}
	s.logg.Info("scheduler cleared")
}
