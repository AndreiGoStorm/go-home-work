package storage

import (
	"context"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	memorystorage "github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	GetEventsByDates(eventStart, eventFinish time.Time) ([]*model.Event, error)
	GetByID(id string) (*model.Event, error)
	Create(event *model.Event) (string, error)
	Update(event *model.Event) error
	Delete(event *model.Event) error
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(conf *config.Config) Storage {
	switch conf.Storage.Driver {
	case "in-memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New(conf.Database)
	default:
		return nil
	}
}
