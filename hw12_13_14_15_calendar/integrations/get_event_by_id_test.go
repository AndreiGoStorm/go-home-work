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
	"testing"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type GetEventRequest struct {
	ID string `json:"ID"`
}

type GetEventResponse struct {
	Event *GetEvent `json:"event"`
}

type GetEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Start       string `json:"start"`
	Finish      string `json:"finish"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
	Remind      int32  `json:"remind"`
	RemindDate  string `json:"remindDate"`
}

type GetEventSuite struct {
	suite.Suite
	store  storage.Storage
	pool   *pgxpool.Pool
	port   string
	host   string
	ctx    context.Context
	client http.Client
	event  *model.Event
}

func NewGetEventSuite() *GetEventSuite {
	return &GetEventSuite{}
}

func (s *GetEventSuite) SetupSuite() {
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

func (s *GetEventSuite) SetupTest() {
	s.event = &model.Event{
		Title:       "title get event",
		Start:       time.Now(),
		Finish:      time.Now().Add(time.Hour * 24).Add(time.Second * 5),
		Description: sql.NullString{String: "description get event", Valid: true},
		UserID:      "8fd5288b-b7fb-4ec1-b8d1-67f017c98704",
		Remind:      1000,
		RemindDate:  time.Now(),
	}
}

func (s *GetEventSuite) TearDownTest() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE events")
}

func (s *GetEventSuite) TestGetEventByID() {
	ID, err := s.store.Create(s.event)
	s.Require().NoError(err)

	getEventRequest := &GetEventRequest{ID}
	reqBody, err := json.Marshal(getEventRequest)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s:%s/v1/events/%s", s.host, s.port, getEventRequest.ID),
		bytes.NewReader(reqBody),
	)
	s.Require().NoError(err)
	response, err := s.client.Do(req)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, response.StatusCode)
	respBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)

	eventResponse := &GetEventResponse{}
	err = json.Unmarshal(respBody, &eventResponse)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Require().Equal(s.event.ID, eventResponse.Event.ID)
	s.Require().Equal(s.event.Title, eventResponse.Event.Title)
	start, _ := time.Parse(time.RFC3339, eventResponse.Event.Start)
	s.Require().Equal(s.event.Start.Format(time.DateTime), start.Format(time.DateTime))
	finish, _ := time.Parse(time.RFC3339, eventResponse.Event.Finish)
	s.Require().Equal(s.event.Finish.Format(time.DateTime), finish.Format(time.DateTime))
	s.Require().Equal(s.event.Description.String, eventResponse.Event.Description)
	s.Require().Equal(s.event.UserID, eventResponse.Event.UserID)
	s.Require().Equal(s.event.Remind, eventResponse.Event.Remind)
}

func TestGetEventSuite(t *testing.T) {
	suite.Run(t, NewGetEventSuite())
}
