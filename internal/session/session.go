package session

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yourname/studyctl/internal/database"
)

type Session struct {
	ID              int64
	Task            string
	Type            string
	StartTime       time.Time
	EndTime         sql.NullTime
	DurationSeconds int64
	Notes           sql.NullString
}

func Start(task string, sessionType string) error {
	if task == "" {
		return errors.New("task cannot be empty")
	}

	if sessionType != "study" && sessionType != "work" {
		return errors.New("type must be 'study' or 'work'")
	}

	active, err := Active()
	if err != nil {
		return err
	}

	if active != nil {
		return errors.New("a session is already running")
	}

	_, err = database.DB.Exec(`
		INSERT INTO sessions
			(task, type, start_time)
		VALUES (?, ?, ?)
	`, task, sessionType, time.Now())

	return err
}

func Stop(note string) (*Session, error) {
	active, err := Active()

	if err != nil {
		return nil, err
	}

	if active == nil {
		return nil, errors.New("no active session")
	}

	end := time.Now()
	duration := int64(end.Sub(active.StartTime).Seconds())

	_, err = database.DB.Exec(`
		UPDATE sessions
		SET
			end_time = ?,
			duration_seconds = ?,
			notes = ?
		WHERE id = ?
	`,
		end,
		duration,
		note,
		active.ID,
	)

	if err != nil {
		return nil, err
	}

	active.EndTime = sql.NullTime{
		Time:  end,
		Valid: true,
	}

	active.DurationSeconds = duration

	active.Notes = sql.NullString{
		String: note,
		Valid: note != "",
	}

	return active, nil
}

func Active() (*Session, error) {
	row := database.DB.QueryRow(`
		SELECT
			id,
			task,
			type,
			start_time,
			end_time,
			duration_seconds,
			notes
		FROM sessions
		WHERE end_time IS NULL
		ORDER BY start_time DESC
		LIMIT 1
	`)

	var s Session

	err := row.Scan(
		&s.ID,
		&s.Task,
		&s.Type,
		&s.StartTime,
		&s.EndTime,
		&s.DurationSeconds,
		&s.Notes,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func Delete(id int64) error {
	result, err := database.DB.Exec(`
		DELETE FROM sessions
		WHERE id = ?
	`, id)

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("session not found")
	}

	return nil
}