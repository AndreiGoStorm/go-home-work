package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" //nolint:nolintlint
)

type Storage struct {
	dns string
	db  *sql.DB
	ctx context.Context
}

func New(db config.Database) *Storage {
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.Username, db.Password, db.DBName)
	return &Storage{dns: dns}
}

func (s *Storage) GetEventsByDates(eventStart, eventFinish time.Time) ([]*model.Event, error) {
	query := `select * from events where start >= $1 AND start < $2 order by start`
	stmp, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmp.Close()

	rows, err := s.db.QueryContext(s.ctx, query, eventStart, eventFinish)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*model.Event, 0, 50)

	for rows.Next() {
		var event model.Event
		err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Start,
			&event.Finish,
			&event.Description,
			&event.UserID,
			&event.Remind)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (s *Storage) GetByID(id string) (*model.Event, error) {
	query := `select * from events where id = $1`
	stmp, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmp.Close()

	row := stmp.QueryRowContext(s.ctx, id)
	event := &model.Event{}
	err = row.Scan(
		&event.ID,
		&event.Title,
		&event.Start,
		&event.Finish,
		&event.Description,
		&event.UserID,
		&event.Remind)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Storage) Create(event *model.Event) (string, error) {
	eventUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	event.ID = eventUUID.String()

	query := `INSERT INTO events (
		id,
		title,
		start,
		finish,
		description,
		user_id,
		remind
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	stmp, err := s.db.PrepareContext(s.ctx, query)
	if err != nil {
		return "", err
	}
	defer stmp.Close()

	start := event.Start.Format(time.DateTime)
	finish := event.Finish.Format(time.DateTime)
	_, err = stmp.ExecContext(s.ctx,
		event.ID,
		event.Title,
		start,
		finish,
		event.Description,
		event.UserID,
		event.Remind)
	if err != nil {
		return "", err
	}

	return event.ID, nil
}

func (s *Storage) Update(event *model.Event) error {
	query := `UPDATE events
		SET	title = $1,
			finish = $2,
			description = $3,
			remind = $4
		WHERE id = $5`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		event.Title,
		event.Finish,
		event.Description,
		event.Remind,
		event.ID)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("event does not exist")
	}

	return nil
}

func (s *Storage) Delete(event *model.Event) error {
	stmt, err := s.db.Prepare(`delete from events where id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(event.ID)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("event does not exist")
	}

	return nil
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", s.dns)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if err := s.db.Close(); err != nil {
		return err
	}

	s.ctx = nil
	return nil
}
