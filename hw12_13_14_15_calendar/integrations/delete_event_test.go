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
)

type DeleteEventRequest struct {
	ID string `json:"id"`
}

type DeleteEventResponse struct {
	ID string `json:"id"`
}

type DeleteEventSuite struct {
	suite.Suite
	store  storage.Storage
	pool   *pgxpool.Pool
	port   string
	host   string
	ctx    context.Context
	client http.Client
	event  *model.Event
}

func NewDeleteEventSuite() *DeleteEventSuite {
	return &DeleteEventSuite{}
}

func (s *DeleteEventSuite) SetupSuite() {
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

func (s *DeleteEventSuite) SetupTest() {
	s.event = &model.Event{
		Title:       "title delete event",
		Start:       time.Now(),
		Finish:      time.Now().Add(time.Hour * 24).Add(time.Second * 5),
		Description: sql.NullString{String: "description delete event", Valid: true},
		UserID:      "8fd5288b-b7fb-4ec1-b8d1-67f017c98704",
		Remind:      0,
		RemindDate:  time.Now(),
	}
}

func (s *DeleteEventSuite) TearDownTest() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE events")
}

func (s *DeleteEventSuite) TestDeleteEvent() {
	ID, err := s.store.Create(s.event)
	s.Require().NoError(err)

	deleteEventRequest := &DeleteEventRequest{ID}
	reqBody, err := json.Marshal(deleteEventRequest)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodDelete,
		fmt.Sprintf("http://%s:%s/v1/events/%s", s.host, s.port, deleteEventRequest.ID),
		bytes.NewReader(reqBody),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Require().Equal(http.StatusOK, response.StatusCode)
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)

	deleteEventResponse := &DeleteEventResponse{}
	err = json.Unmarshal(respBody, &deleteEventResponse)
	s.Require().NoError(err)
	_, err = uuid.Parse(deleteEventResponse.ID)
	s.Require().NoError(err)

	actual, err := s.store.GetByID(ID)
	s.Require().Nil(actual)
	s.Require().EqualError(err, "sql: no rows in result set")
}

func (s *DeleteEventSuite) TestDeleteEventWithWrongId() {
	deleteEventRequest := &DeleteEventRequest{"ID"}
	reqBody, err := json.Marshal(deleteEventRequest)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodDelete,
		fmt.Sprintf("http://%s:%s/v1/events/%s", s.host, s.port, deleteEventRequest.ID),
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

func TestDeleteEventSuite(t *testing.T) {
	suite.Run(t, NewDeleteEventSuite())
}
