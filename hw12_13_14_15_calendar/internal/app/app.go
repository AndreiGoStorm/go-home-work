package app

import (
	"context"
	"errors"
	"fmt"

	ep "github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/api/pb/event"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model/converter"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logg    *logger.Logger
	storage storage.Storage

	ep.UnimplementedEventServiceServer
}

func New(logg *logger.Logger, storage storage.Storage) *App {
	return &App{logg: logg, storage: storage}
}

func (a *App) GetEventByID(ctx context.Context, req *ep.GetEventByIDRequest) (res *ep.GetEventByIDResponse, err error) {
	defer ctx.Done()

	event, err := a.storage.GetByID(req.GetID())
	if err != nil {
		a.logg.Warn("GetEventByID", err)
		return nil, model.ErrEventNotFound
	}

	return &ep.GetEventByIDResponse{
		Event: converter.EventModelToProto(event),
	}, nil
}

func (a *App) GetEventsByDay(ctx context.Context, req *ep.GetEventsByDayRequest) (
	*ep.GetEventsByDayResponse,
	error,
) {
	defer ctx.Done()

	if req.GetDate() == nil {
		err := errors.New("wrong date")
		a.logg.Warn("GetEventsByDay", err)
		return nil, err
	}
	eventStart, eventFinish := converter.GetDayDatesFromProto(req)
	events, err := a.storage.GetEventsByDates(eventStart, eventFinish)
	if err != nil {
		a.logg.Warn("GetEventsByDay", err)
		return nil, model.ErrEventNotFound
	}

	return &ep.GetEventsByDayResponse{
		Events: converter.EventModelsToProtos(events),
	}, nil
}

func (a *App) GetEventsByWeek(ctx context.Context, req *ep.GetEventsByWeekRequest) (
	*ep.GetEventsByWeekResponse,
	error,
) {
	defer ctx.Done()

	if req.GetDate() == nil {
		err := errors.New("wrong date")
		a.logg.Warn("GetEventsByWeek", err)
		return nil, err
	}
	eventStart, eventFinish := converter.GetWeekDatesFromProto(req)
	events, err := a.storage.GetEventsByDates(eventStart, eventFinish)
	if err != nil {
		a.logg.Warn("GetEventsByWeek", err)
		return nil, model.ErrEventNotFound
	}

	return &ep.GetEventsByWeekResponse{
		Events: converter.EventModelsToProtos(events),
	}, nil
}

func (a *App) GetEventsByMonth(ctx context.Context, req *ep.GetEventsByMonthRequest) (
	*ep.GetEventsByMonthResponse,
	error,
) {
	defer ctx.Done()

	if req.GetDate() == nil {
		err := errors.New("wrong date")
		a.logg.Warn("GetEventsByMonth", err)
		return nil, err
	}
	eventStart, eventFinish := converter.GetMonthDatesFromProto(req)
	events, err := a.storage.GetEventsByDates(eventStart, eventFinish)
	if err != nil {
		a.logg.Warn("GetEventsByMonth", err)
		return nil, model.ErrEventNotFound
	}

	return &ep.GetEventsByMonthResponse{
		Events: converter.EventModelsToProtos(events),
	}, nil
}

func (a *App) CreateEvent(ctx context.Context, req *ep.CreateEventRequest) (*ep.CreateEventResponse, error) {
	defer ctx.Done()

	event := converter.CreateEventRequestToModel(req)
	if event == nil {
		err := fmt.Errorf("can not convert proto created model to event model")
		a.logg.Warn("CreateEvent", err)
		return nil, err
	}

	id, err := a.storage.Create(event)
	if err != nil {
		a.logg.Warn("CreateEvent", err)
		return nil, err
	}

	return &ep.CreateEventResponse{
		ID: id,
	}, nil
}

func (a *App) UpdateEvent(ctx context.Context, req *ep.UpdateEventRequest) (*ep.UpdateEventResponse, error) {
	defer ctx.Done()

	existed, err := a.storage.GetByID(req.GetID())
	if err != nil {
		a.logg.Warn("UpdateEvent", err)
		return nil, model.ErrEventNotFound
	}

	event := converter.UpdateEventRequestToModel(req)
	if event == nil {
		err := fmt.Errorf("can not convert proto updated model to event model")
		a.logg.Warn("UpdateEvent", err)
		return nil, err
	}

	event.Start = existed.Start
	event.UserID = existed.UserID
	err = a.storage.Update(event)
	if err != nil {
		a.logg.Warn("UpdateEvent", err)
		return nil, err
	}
	return &ep.UpdateEventResponse{ID: event.ID}, nil
}

func (a *App) DeleteEvent(ctx context.Context, req *ep.DeleteEventRequest) (*ep.DeleteEventResponse, error) {
	defer ctx.Done()

	event, err := a.storage.GetByID(req.GetID())
	if err != nil {
		a.logg.Warn("DeleteEvent", err)
		return nil, model.ErrEventNotFound
	}

	err = a.storage.Delete(event)
	if err != nil {
		a.logg.Warn("DeleteEvent", err)
		return nil, err
	}

	return &ep.DeleteEventResponse{ID: event.ID}, nil
}
