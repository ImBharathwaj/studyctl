package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".studyctl")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dir, "studyctl.db")

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		type TEXT NOT NULL DEFAULT 'study',
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		duration_seconds INTEGER DEFAULT 0,
		notes TEXT
	);

	CREATE TABLE IF NOT EXISTS goals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL UNIQUE,
		target_seconds INTEGER NOT NULL
	);
	`

	_, err := DB.Exec(query)

	return err
}