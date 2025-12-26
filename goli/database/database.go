package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
	var err error
	// SQLite database file path
	dbPath := "/goli/data/goli.db"

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Create tables if they don't exist
	if err = createTables(); err != nil {
		return err
	}

	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// createTables creates the necessary database tables
func createTables() error {
	// This is a simplified version - you may need to adjust based on your actual schema
	queries := []string{
		`CREATE TABLE IF NOT EXISTS pipelines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			definition TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pipeline_id INTEGER,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			triggered_by TEXT,
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			logs TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (pipeline_id) REFERENCES pipelines(id)
		)`,
		`CREATE TABLE IF NOT EXISTS job_steps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			step_name TEXT NOT NULL,
			step_order INTEGER NOT NULL,
			status TEXT NOT NULL,
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			logs TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (job_id) REFERENCES jobs(id)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT,
			phone TEXT,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			two_fa_email_enabled INTEGER DEFAULT 0,
			two_fa_sms_enabled INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS two_factor_codes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			channel TEXT NOT NULL,
			code TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			consumed INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pipeline_variables (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pipeline_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			value TEXT NOT NULL,
			is_secret INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (pipeline_id) REFERENCES pipelines(id) ON DELETE CASCADE,
			UNIQUE(pipeline_id, name)
		)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
