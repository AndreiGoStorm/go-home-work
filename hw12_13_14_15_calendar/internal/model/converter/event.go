package converter

import (
	"database/sql"
	"errors"
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
		RemindDate:  timestamppb.New(event.RemindDate),
	}
}

func EventModelsToProtos(events []*model.Event) (result []*ep.Event) {
	for _, event := range events {
		result = append(result, EventModelToProto(event))
	}
	return result
}

func GetDayDatesFromProto(req *ep.GetEventsByDayRequest) (*time.Time, *time.Time, error) {
	if req.GetDate() == nil {
		return nil, nil, errors.New("wrong date")
	}
	date := req.GetDate().AsTime()
	start := date.Truncate(24 * time.Hour)
	finish := start.Add(24 * time.Hour)
	return &start, &finish, nil
}

func GetWeekDatesFromProto(req *ep.GetEventsByWeekRequest) (*time.Time, *time.Time, error) {
	if req.GetDate() == nil {
		return nil, nil, errors.New("wrong date")
	}
	date := req.GetDate().AsTime()
	start := date.Truncate(7 * 24 * time.Hour)
	finish := start.AddDate(0, 0, 7)
	return &start, &finish, nil
}

func GetMonthDatesFromProto(req *ep.GetEventsByMonthRequest) (*time.Time, *time.Time, error) {
	if req.GetDate() == nil {
		return nil, nil, errors.New("wrong date")
	}
	date := req.GetDate().AsTime()
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	finish := start.AddDate(0, 1, 0)
	return &start, &finish, nil
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

	start := event.GetStart().AsTime()
	remindDate := GetRemindDate(start, int(event.GetRemind()))

	return &model.Event{
		Title:       event.GetTitle(),
		Start:       start,
		Finish:      event.GetFinish().AsTime(),
		Description: description,
		UserID:      event.GetUserID(),
		Remind:      event.GetRemind(),
		RemindDate:  remindDate,
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

func GetRemindDate(date time.Time, remind int) time.Time {
	remindDate := date
	if remind > 0 {
		remindDate = remindDate.AddDate(0, 0, (-1)*remind)
	}
	return remindDate
}
