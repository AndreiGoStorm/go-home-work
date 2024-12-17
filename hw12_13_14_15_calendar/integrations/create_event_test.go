package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateEventRequest struct {
	Event *CreateEvent `json:"event"`
}

type CreateEvent struct {
	Title       string `json:"title"`
	Start       string `json:"start"`
	Finish      string `json:"finish"`
	Description string `json:"description"`
	UserID      string `json:"userID"`
	Remind      int32  `json:"remind"`
}

type CreateEventResponse struct {
	ID string `json:"ID"`
}

type CreateEventSuite struct {
	suite.Suite
	store       storage.Storage
	pool        *pgxpool.Pool
	port        string
	host        string
	ctx         context.Context
	client      http.Client
	createEvent *CreateEvent
}

func NewCreateEventSuite() *CreateEventSuite {
	return &CreateEventSuite{}
}

func (s *CreateEventSuite) SetupSuite() {
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

func (s *CreateEventSuite) SetupTest() {
	start := time.Now()
	finish := time.Now().Add(time.Hour * 24)
	s.createEvent = &CreateEvent{
		Title:       "title for create event",
		Start:       start.Format(time.RFC3339),
		Finish:      finish.Format(time.RFC3339),
		Description: "description for create event",
		UserID:      "8fd5288b-b7fb-4ec1-b8d1-67f017c98704",
		Remind:      0,
	}
}

func (s *CreateEventSuite) TearDownTest() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE events")
}

func (s *CreateEventSuite) TestCreateEvent() {
	createEventRequest := &CreateEventRequest{
		Event: s.createEvent,
	}
	reqBody, err := json.Marshal(createEventRequest)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
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

	createEventResponse := &CreateEventResponse{}
	err = json.Unmarshal(respBody, &createEventResponse)
	s.Require().NoError(err)
	_, err = uuid.Parse(createEventResponse.ID)
	s.Require().NoError(err)

	actual, err := s.store.GetByID(createEventResponse.ID)
	s.Require().NoError(err)
	s.Require().Equal(s.createEvent.Title, actual.Title)
	start, _ := time.Parse(time.RFC3339, s.createEvent.Start)
	s.Require().Equal(timestamppb.New(start).GetSeconds(), timestamppb.New(actual.Start).GetSeconds())
	finish, _ := time.Parse(time.RFC3339, s.createEvent.Finish)
	s.Require().Equal(timestamppb.New(finish).GetSeconds(), timestamppb.New(actual.Finish).GetSeconds())
	s.Require().Equal(s.createEvent.Description, actual.Description.String)
	s.Require().Equal(s.createEvent.UserID, actual.UserID)
	s.Require().Equal(s.createEvent.Remind, actual.Remind)
}

func (s *CreateEventSuite) TestCreateEventWithWrongData() {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s:%s/v1/events", s.host, s.port),
		bytes.NewReader(nil),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusInternalServerError, response.StatusCode)
	is := strings.Contains(string(respBody), "can not convert proto created model to event model")
	s.Require().True(is)
}

func TestCreateEventSuite(t *testing.T) {
	suite.Run(t, NewCreateEventSuite())
}
