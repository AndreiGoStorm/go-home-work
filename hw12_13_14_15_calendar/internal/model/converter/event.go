package converter

import (
	"database/sql"
	"time"

	ep "github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/api/pb/event"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventModelToProto(event *model.Event) *ep.Event {
	return &ep.Event{
		ID:          event.ID,
		Title:       event.Title,
		Start:       timestamppb.New(event.Start),
		Finish:      timestamppb.New(event.Finish),
		Description: event.Description.String,
		UserID:      event.UserID,
		Remind:      event.Remind,
	}
}

func EventModelsToProtos(events []*model.Event) (result []*ep.Event) {
	for _, event := range events {
		result = append(result, EventModelToProto(event))
	}
	return result
}

func GetDayDatesFromProto(req *ep.GetEventsByDayRequest) (eventStart, eventFinish time.Time) {
	date := req.GetDate().AsTime()
	eventStart = date.Truncate(24 * time.Hour)
	eventFinish = eventStart.Add(24 * time.Hour)
	return eventStart, eventFinish
}

func GetWeekDatesFromProto(req *ep.GetEventsByWeekRequest) (eventStart, eventFinish time.Time) {
	date := req.GetDate().AsTime()
	eventStart = date.Truncate(7 * 24 * time.Hour)
	eventFinish = eventStart.AddDate(0, 0, 7)
	return eventStart, eventFinish
}

func GetMonthDatesFromProto(req *ep.GetEventsByMonthRequest) (eventStart, eventFinish time.Time) {
	date := req.GetDate().AsTime()
	eventStart = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	eventFinish = eventStart.AddDate(0, 1, 0)
	return eventStart, eventFinish
}

func CreateEventRequestToModel(req *ep.CreateEventRequest) *model.Event {
	event := req.GetEvent()
	if event == nil {
		return nil
	}

	description := sql.NullString{String: event.GetDescription(), Valid: true}
	if description.String == "" {
		description.Valid = false
	}

	return &model.Event{
		Title:       event.GetTitle(),
		Start:       event.GetStart().AsTime(),
		Finish:      event.GetFinish().AsTime(),
		Description: description,
		UserID:      event.GetUserID(),
		Remind:      event.GetRemind(),
	}
}

func UpdateEventRequestToModel(req *ep.UpdateEventRequest) *model.Event {
	event := req.GetEvent()
	if event == nil {
		return nil
	}

	description := sql.NullString{String: event.GetDescription(), Valid: true}
	if description.String == "" {
		description.Valid = false
	}

	return &model.Event{
		ID:          req.GetID(),
		Title:       event.GetTitle(),
		Finish:      event.GetFinish().AsTime(),
		Description: description,
		Remind:      event.GetRemind(),
	}
}
