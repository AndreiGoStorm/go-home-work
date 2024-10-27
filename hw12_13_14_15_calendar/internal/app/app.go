package app

import (
	"context"
	"fmt"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logg    *logger.Logger
	storage storage.Storage
}

func New(logg *logger.Logger, storage storage.Storage) *App {
	return &App{logg, storage}
}

func (a *App) CreateEvent(_ context.Context, title string) error {
	event := model.Event{
		Title: title,
	}
	_, err := a.storage.Create(event)
	if err != nil {
		err := fmt.Errorf("error creating event: %w", err)
		a.logg.Error(err.Error())
		return err
	}

	return nil
}
