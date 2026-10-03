package session

import (
	"database/sql"
	"time"
	"fmt"
	"github.com/yourname/studyctl/internal/database"
)

type DailyStats struct {
	TotalSeconds int64
	StudySeconds int64
	WorkSeconds  int64
}

type HourStats struct {
	Hour         int
	StudySeconds int64
	WorkSeconds  int64
}

func GetDailyStats(date time.Time) (DailyStats, error) {
	dateString := date.Format("2006-01-02")

	var stats DailyStats

	err := database.DB.QueryRow(`
		SELECT
			COALESCE(SUM(duration_seconds), 0),
			COALESCE(SUM(
				CASE WHEN type = 'study'
				THEN duration_seconds ELSE 0 END
			), 0),
			COALESCE(SUM(
				CASE WHEN type = 'work'
				THEN duration_seconds ELSE 0 END
			), 0)
		FROM sessions
		WHERE date(start_time, 'localtime') = ?
		AND end_time IS NOT NULL
	`, dateString).Scan(
		&stats.TotalSeconds,
		&stats.StudySeconds,
		&stats.WorkSeconds,
	)

	return stats, err
}

func GetHourlyStats(date time.Time) ([]HourStats, error) {
	dateString := date.Format("2006-01-02")

	rows, err := database.DB.Query(`
		SELECT
			strftime('%H', start_time, 'localtime') AS hour,

			COALESCE(SUM(
				CASE WHEN type = 'study'
				THEN duration_seconds ELSE 0 END
			), 0),

			COALESCE(SUM(
				CASE WHEN type = 'work'
				THEN duration_seconds ELSE 0 END
			), 0)

		FROM sessions

		WHERE date(start_time, 'localtime') = ?
		AND end_time IS NOT NULL

		GROUP BY hour
		ORDER BY hour
	`, dateString)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []HourStats

	for rows.Next() {
		var hourString string
		var stats HourStats

		err := rows.Scan(
			&hourString,
			&stats.StudySeconds,
			&stats.WorkSeconds,
		)

		if err != nil {
			return nil, err
		}

		var hour int
		_, err = fmt.Sscanf(hourString, "%d", &hour)

		if err != nil {
			return nil, err
		}

		stats.Hour = hour

		result = append(result, stats)
	}

	return result, rows.Err()
}

type TaskStats struct {
	Task    string
	Type    string
	Seconds int64
}

func GetTaskStats(date time.Time) ([]TaskStats, error) {
	dateString := date.Format("2006-01-02")

	rows, err := database.DB.Query(`
		SELECT
			task,
			type,
			SUM(duration_seconds)
		FROM sessions
		WHERE date(start_time, 'localtime') = ?
		AND end_time IS NOT NULL
		GROUP BY task, type
		ORDER BY SUM(duration_seconds) DESC
	`, dateString)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []TaskStats

	for rows.Next() {
		var stats TaskStats

		err := rows.Scan(
			&stats.Task,
			&stats.Type,
			&stats.Seconds,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, stats)
	}

	return result, rows.Err()
}

func GetSessions(date time.Time) (*sql.Rows, error) {
	dateString := date.Format("2006-01-02")

	return database.DB.Query(`
		SELECT
			id,
			task,
			type,
			start_time,
			end_time,
			duration_seconds
		FROM sessions
		WHERE date(start_time, 'localtime') = ?
		ORDER BY start_time
	`, dateString)
}