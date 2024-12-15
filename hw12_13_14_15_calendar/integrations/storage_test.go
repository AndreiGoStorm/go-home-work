package integrations

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	store storage.Storage
	pool  *pgxpool.Pool
	event model.Event
}

func NewStorageSuite() *StorageSuite {
	return &StorageSuite{}
}

func (s *StorageSuite) SetupSuite() {
	conf := SetupSuite()
	s.pool = PoolConnect(conf)
	s.store = StorageConnect(conf)
}

func (s *StorageSuite) SetupTest() {
	s.event = model.Event{
		Title:       "title sql event",
		Start:       time.Now(),
		Finish:      time.Now().Add(time.Hour * 24).Add(time.Second * 5),
		Description: sql.NullString{String: "description", Valid: true},
		UserID:      "8fd5288b-b7fb-4ec1-b8d1-67f017c98704",
		Remind:      3600,
		RemindDate:  time.Now(),
	}
}

func (s *StorageSuite) TearDownTest() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE events")
}

func (s *StorageSuite) TestStorageCreate() {
	ID, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	createdEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)
	s.Require().Equal(ID, createdEvent.ID)
}

func (s *StorageSuite) TestStorageUpdate() {
	ID, err := s.store.Create(&s.event)
	s.Require().NoError(err)
	createdEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)

	newTitle := "new title"
	newDescription := sql.NullString{String: "new description", Valid: false}
	newRemind := int32(2)
	newRemindDate := time.Now().Add(24 * time.Hour)
	createdEvent.Title = newTitle
	createdEvent.Description = newDescription
	createdEvent.Remind = newRemind
	createdEvent.RemindDate = newRemindDate

	err = s.store.Update(createdEvent)
	s.Require().NoError(err)

	updatedEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)
	s.Require().Equal(newTitle, updatedEvent.Title)
	s.Require().Equal(false, updatedEvent.Description.Valid)
	s.Require().Equal(newRemind, updatedEvent.Remind)
	s.Require().Equal(newRemindDate.Format(time.DateTime), updatedEvent.RemindDate.Format(time.DateTime))
}

func (s *StorageSuite) TestStorageUpdateWithWrongId() {
	ID, err := s.store.Create(&s.event)
	s.Require().NoError(err)
	createdEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)

	createdEvent.ID = "not existed ID"
	err = s.store.Update(createdEvent)
	s.Require().EqualError(err, "event does not exist")
}

func (s *StorageSuite) TestStorageDelete() {
	ID, err := s.store.Create(&s.event)
	s.Require().NoError(err)
	createdEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)

	err = s.store.Delete(createdEvent)
	s.Require().NoError(err)

	createdEvent, err = s.store.GetByID(ID)
	s.Require().Nil(createdEvent)
	s.Require().EqualError(err, "sql: no rows in result set")
}

func (s *StorageSuite) TestStorageDeleteWithWrongId() {
	ID, err := s.store.Create(&s.event)
	s.Require().NoError(err)
	createdEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)

	createdEvent.ID = "not existed ID"
	err = s.store.Delete(createdEvent)
	s.Require().EqualError(err, "event does not exist")
}

func (s *StorageSuite) TestStorageDeleteOldEvents() {
	_, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	eventForDelete := s.event
	eventForDelete.RemindDate = time.Now().AddDate(-1, -1, -1)
	ID, err := s.store.Create(&eventForDelete)
	s.Require().NoError(err)

	err = s.store.DeleteOldEvents()
	s.Require().NoError(err)

	deletedEvent, err := s.store.GetByID(ID)
	s.Require().Nil(deletedEvent)
	s.Require().EqualError(err, "sql: no rows in result set")
}

func (s *StorageSuite) TestStorageGetById() {
	ID, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	createdEvent, err := s.store.GetByID(ID)
	s.Require().NoError(err)
	s.Require().Equal(ID, createdEvent.ID)
	s.Require().Equal(s.event.Title, createdEvent.Title)
	s.Require().Equal(s.event.Start.Format(time.DateTime), createdEvent.Start.Format(time.DateTime))
	s.Require().Equal(s.event.Finish.Format(time.DateTime), createdEvent.Finish.Format(time.DateTime))
	s.Require().Equal(s.event.Description, createdEvent.Description)
	s.Require().Equal(s.event.UserID, createdEvent.UserID)
	s.Require().Equal(s.event.Remind, createdEvent.Remind)
	s.Require().Equal(s.event.RemindDate.Format(time.DateTime), createdEvent.RemindDate.Format(time.DateTime))
}

func (s *StorageSuite) TestStorageGetRemindEvents() {
	_, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	startEvent := time.Now().AddDate(1, 0, 0)
	s.event.RemindDate = startEvent
	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	remindEvents, err := s.store.GetRemindEvents(startEvent)
	s.Require().NoError(err)
	s.Require().Len(remindEvents, 1)
}

func (s *StorageSuite) TestStorageGetOnDay() {
	start := time.Now().AddDate(0, 5, 0)
	s.event.Start = start

	_, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	s.event.Start = start.Add(-24 * time.Hour)
	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	startOfDay := start.Truncate(24 * time.Hour)
	eventsOnDay, err := s.store.GetEventsByDates(startOfDay, startOfDay.Add(24*time.Hour))
	s.Require().NoError(err)
	s.Require().Len(eventsOnDay, 2)
}

func (s *StorageSuite) TestStorageGetOnWeek() {
	start := time.Now().AddDate(0, 5, 0)
	s.event.Start = start

	_, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	wrongStart := time.Now().AddDate(0, 5, 0)
	s.event.Start = wrongStart.Add(7 * (-24) * time.Hour)
	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	startOfWeek := start.Truncate(7 * 24 * time.Hour)
	eventsOnWeek, err := s.store.GetEventsByDates(startOfWeek, startOfWeek.AddDate(0, 0, 7))
	s.Require().NoError(err)
	s.Require().Len(eventsOnWeek, 2)
}

func (s *StorageSuite) TestStorageGetOnMonth() {
	start := time.Now().AddDate(0, 5, 0)
	s.event.Start = start

	_, err := s.store.Create(&s.event)
	s.Require().NoError(err)

	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	wrongStart := time.Now().AddDate(0, 5, 0)
	s.event.Start = wrongStart.AddDate(0, -1, 1)
	_, err = s.store.Create(&s.event)
	s.Require().NoError(err)

	startOfMonth := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	eventsOnMonth, err := s.store.GetEventsByDates(startOfMonth, startOfMonth.AddDate(0, 1, 0))
	s.Require().NoError(err)
	s.Require().Len(eventsOnMonth, 2)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, NewStorageSuite())
}
