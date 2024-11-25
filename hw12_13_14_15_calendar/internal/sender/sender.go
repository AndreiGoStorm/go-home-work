package sender

import (
	"context"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/rabbit"
)

type Sender struct {
	rabbit *rabbit.Rabbit
	logg   *logger.Logger
}

func New(rabbit *rabbit.Rabbit, logg *logger.Logger) *Sender {
	return &Sender{
		rabbit: rabbit,
		logg:   logg,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	go func() {
		if err := s.rabbit.Read(ctx); err != nil {
			s.logg.Error("sender run read", err)
			return
		}
	}()

	<-ctx.Done()
	return nil
}
