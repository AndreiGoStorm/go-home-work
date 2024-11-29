package app

import (
	"context"
	"path"
	"path/filepath"
	"testing"
	"time"

	ep "github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/api/pb/event"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestApp(t *testing.T) {
	app, ctx := initApp()
	t.Run("get event by id", func(t *testing.T) {
		createdReq := &ep.CreateEventRequest{Event: generateCreateEvent()}
		createdRes, err := app.CreateEvent(ctx, createdReq)
		require.NoError(t, err)

		getRes, err := app.GetEventByID(ctx, &ep.GetEventByIDRequest{ID: createdRes.ID})
		require.NoError(t, err)
		require.Equal(t, createdReq.Event.Title, getRes.Event.Title)
		expectedStart := createdReq.Event.Start.AsTime().Format(time.DateTime)
		actualStart := getRes.Event.Start.AsTime().Format(time.DateTime)
		require.Equal(t, expectedStart, actualStart)
		expectedFinish := createdReq.Event.Finish.AsTime().Format(time.DateTime)
		actualFinish := getRes.Event.Finish.AsTime().Format(time.DateTime)
		require.Equal(t, expectedFinish, actualFinish)
		require.Equal(t, createdReq.Event.Description, getRes.Event.Description)
		require.Equal(t, createdReq.Event.UserID, getRes.Event.UserID)
		require.Equal(t, createdReq.Event.Remind, getRes.Event.Remind)
		require.Equal(t, createdReq.Event.Remind, getRes.Event.Remind)
	})

	t.Run("create event by id", func(t *testing.T) {
		createdReq := &ep.CreateEventRequest{Event: generateCreateEvent()}
		createdRes, err := app.CreateEvent(ctx, createdReq)
		require.NoError(t, err)

		getRes, err := app.GetEventByID(ctx, &ep.GetEventByIDRequest{ID: createdRes.ID})
		require.NoError(t, err)
		require.Equal(t, createdRes.ID, getRes.Event.ID)
	})

	t.Run("update event", func(t *testing.T) {
		createdReq := &ep.CreateEventRequest{Event: generateCreateEvent()}
		createdRes, err := app.CreateEvent(ctx, createdReq)
		require.NoError(t, err)

		updatedReq := &ep.UpdateEventRequest{ID: createdRes.ID, Event: generateUpdateEvent()}
		updatedRes, err := app.UpdateEvent(ctx, updatedReq)
		require.NoError(t, err)
		require.Equal(t, createdRes.ID, updatedRes.ID)

		getRes, err := app.GetEventByID(ctx, &ep.GetEventByIDRequest{ID: updatedRes.ID})
		require.NoError(t, err)
		require.Equal(t, updatedReq.Event.Title, getRes.Event.Title)
		expectedFinish := updatedReq.Event.Finish.AsTime().Format(time.DateTime)
		actualFinish := getRes.Event.Finish.AsTime().Format(time.DateTime)
		require.Equal(t, expectedFinish, actualFinish)
		require.Equal(t, updatedReq.Event.Description, getRes.Event.Description)
		require.Equal(t, updatedReq.Event.Remind, getRes.Event.Remind)
		require.Equal(t, createdReq.Event.UserID, getRes.Event.UserID)
		expectedStart := createdReq.Event.Start.AsTime().Format(time.DateTime)
		actualStart := getRes.Event.Start.AsTime().Format(time.DateTime)
		require.Equal(t, expectedStart, actualStart)
		remind := int(updatedReq.Event.Remind)
		start := getRes.Event.Start.AsTime()
		expectedStart = start.AddDate(0, 0, (-1)*remind).Format(time.DateTime)
		actualFinish = getRes.Event.RemindDate.AsTime().Format(time.DateTime)
		require.Equal(t, expectedStart, actualFinish)
	})

	t.Run("update event with not existed ID", func(t *testing.T) {
		createdReq := &ep.CreateEventRequest{Event: generateCreateEvent()}
		createdRes, err := app.CreateEvent(ctx, createdReq)
		require.NotNil(t, createdRes)
		require.NoError(t, err)

		updatedReq := &ep.UpdateEventRequest{ID: "not existed ID", Event: generateUpdateEvent()}
		updatedRes, err := app.UpdateEvent(ctx, updatedReq)
		require.Nil(t, updatedRes)
		require.EqualError(t, err, "event not found in storage")
	})

	t.Run("delete event", func(t *testing.T) {
		createdReq := &ep.CreateEventRequest{Event: generateCreateEvent()}
		createdRes, err := app.CreateEvent(ctx, createdReq)
		require.NoError(t, err)

		delRes, err := app.DeleteEvent(ctx, &ep.DeleteEventRequest{ID: createdRes.ID})
		require.NoError(t, err)
		require.Equal(t, createdRes.ID, delRes.ID)

		getRes, err := app.GetEventByID(ctx, &ep.GetEventByIDRequest{ID: delRes.ID})
		require.Nil(t, getRes)
		require.EqualError(t, err, "event not found in storage")
	})

	t.Run("delete event with not existed ID", func(t *testing.T) {
		createdReq := &ep.CreateEventRequest{Event: generateCreateEvent()}
		createdRes, err := app.CreateEvent(ctx, createdReq)
		require.NotNil(t, createdRes)
		require.NoError(t, err)

		delRes, err := app.DeleteEvent(ctx, &ep.DeleteEventRequest{ID: "not existed ID"})
		require.Nil(t, delRes)
		require.EqualError(t, err, "event not found in storage")
	})
}

func generateCreateEvent() *ep.CreateEvent {
	UUID, _ := uuid.NewUUID()
	start := timestamppb.New(time.Now().Add(time.Hour * -24))
	finish := timestamppb.New(time.Now())

	return &ep.CreateEvent{
		Title:       "create event title",
		Start:       start,
		Finish:      finish,
		Description: "create event description",
		UserID:      UUID.String(),
		Remind:      int32(1),
	}
}

func generateUpdateEvent() *ep.UpdateEvent {
	finish := timestamppb.New(time.Now().Add(time.Hour * 24))
	return &ep.UpdateEvent{
		Title:       "update event title",
		Finish:      finish,
		Description: "update event description",
		Remind:      int32(2),
	}
}

func initApp() (*App, context.Context) {
	file, err := filepath.Abs("main.go")
	if err != nil {
		panic(err)
	}
	configFile := "/../../configs/config-testing.toml"
	conf := config.New(path.Dir(file) + configFile)
	logg := logger.New(conf.Logger.Level)
	store := storage.New(conf)
	ctx := context.Background()
	err = store.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return New(logg, store), ctx
}
