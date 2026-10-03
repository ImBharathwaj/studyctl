package database

import (
	"database/sql"
	"time"
)

func SetGoal(date time.Time, seconds int64) error {
	dateString := date.Format("2006-01-02")

	_, err := DB.Exec(`
		INSERT INTO goals (date, target_seconds)
		VALUES (?, ?)
		ON CONFLICT(date)
		DO UPDATE SET target_seconds = excluded.target_seconds
	`,
		dateString,
		seconds,
	)

	return err
}

func GetGoal(date time.Time) (int64, error) {
	dateString := date.Format("2006-01-02")

	var seconds int64

	err := DB.QueryRow(`
		SELECT target_seconds
		FROM goals
		WHERE date = ?
	`, dateString).Scan(&seconds)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return seconds, err
}