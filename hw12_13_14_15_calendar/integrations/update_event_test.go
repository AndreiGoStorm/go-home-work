package integrations

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UpdateEventRequest struct {
	ID    string       `json:"ID"`
	Event *UpdateEvent `json:"event"`
}

type UpdateEvent struct {
	Title       string `json:"title"`
	Finish      string `json:"finish"`
	Description string `json:"description"`
	Remind      int32  `json:"remind"`
}

type UpdateEventResponse struct {
	ID string `json:"ID"`
}

type UpdateEventSuite struct {
	suite.Suite
	store       storage.Storage
	pool        *pgxpool.Pool
	port        string
	host        string
	ctx         context.Context
	client      http.Client
	event       *model.Event
	updateEvent *UpdateEvent
}

func NewUpdateEventSuite() *UpdateEventSuite {
	return &UpdateEventSuite{}
}

func (s *UpdateEventSuite) SetupSuite() {
	conf := SetupSuite()
	s.pool = PoolConnect(conf)
	s.store = StorageConnect(conf)
	s.host = conf.HTTPServer.Host
	s.port = strconv.Itoa(conf.HTTPServer.Port)
	s.ctx = context.Background()
	s.client = http.Client{
		Timeout: 30 * time.Second,
	}
}

func (s *UpdateEventSuite) SetupTest() {
	s.event = &model.Event{
		Title:       "title update event",
		Start:       time.Now(),
		Finish:      time.Now().Add(time.Hour * 24).Add(time.Second * 5),
		Description: sql.NullString{String: "description update event", Valid: true},
		UserID:      "8fd5288b-b7fb-4ec1-b8d1-67f017c98704",
		Remind:      3600,
		RemindDate:  time.Now(),
	}
	finish := time.Now().Add(time.Hour * 24 * 2)
	s.updateEvent = &UpdateEvent{
		Title:       "title for update event",
		Finish:      finish.Format(time.RFC3339),
		Description: "description for create event",
		Remind:      3,
	}
}

func (s *UpdateEventSuite) TearDownTest() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE events")
}

func (s *UpdateEventSuite) TestUpdateEvent() {
	ID, err := s.store.Create(s.event)
	s.Require().NoError(err)

	updateEventRequest := &UpdateEventRequest{
		ID:    ID,
		Event: s.updateEvent,
	}
	reqBody, err := json.Marshal(updateEventRequest)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		fmt.Sprintf("http://%s:%s/v1/events", s.host, s.port),
		bytes.NewReader(reqBody),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Require().Equal(http.StatusOK, response.StatusCode)
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)

	updateEventResponse := &UpdateEventResponse{}
	err = json.Unmarshal(respBody, &updateEventResponse)
	s.Require().NoError(err)
	_, err = uuid.Parse(updateEventResponse.ID)
	s.Require().NoError(err)

	actual, err := s.store.GetByID(ID)
	s.Require().NoError(err)
	s.Require().Equal(s.updateEvent.Title, actual.Title)
	finish, _ := time.Parse(time.RFC3339, s.updateEvent.Finish)
	s.Require().Equal(timestamppb.New(finish).GetSeconds(), timestamppb.New(actual.Finish).GetSeconds())
	s.Require().Equal(s.updateEvent.Description, actual.Description.String)
	s.Require().Equal(s.updateEvent.Remind, actual.Remind)
}

func (s *UpdateEventSuite) TestUpdateEventWithWrongId() {
	finish := time.Now().Add(time.Hour * 24 * 2)
	updateEvent := &UpdateEvent{
		Title:       "title for update event",
		Finish:      finish.Format(time.RFC3339),
		Description: "description for create event",
		Remind:      3,
	}
	updateEventRequest := &UpdateEventRequest{
		ID:    "not existed ID",
		Event: updateEvent,
	}

	reqBody, err := json.Marshal(updateEventRequest)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		fmt.Sprintf("http://%s:%s/v1/events", s.host, s.port),
		bytes.NewReader(reqBody),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Require().Equal(http.StatusInternalServerError, response.StatusCode)

	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	is := strings.Contains(string(respBody), "event not found in storage")
	s.Require().True(is)
}

func TestUpdateEventSuite(t *testing.T) {
	suite.Run(t, NewUpdateEventSuite())
}
