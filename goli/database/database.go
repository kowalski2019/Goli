package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase initializes the SQLite database and creates necessary tables
func InitDatabase() error {
	dbPath := "/goli/data/goli.db"
	if os.Getenv("GOOS") == "windows" || os.Getenv("OS") == "Windows_NT" {
		dbPath = "C:\\goli\\data\\goli.db"
		os.MkdirAll("C:\\goli\\data", 0755)
	} else {
		os.MkdirAll("/goli/data", 0755)
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	if err != nil {
		return err
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Create tables
	if err = createTables(); err != nil {
		return err
	}

	log.Println("Database initialized successfully at", dbPath)
	return nil
}

func createTables() error {
	queries := []string{
		// Pipelines table
		`CREATE TABLE IF NOT EXISTS pipelines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			definition TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Jobs table
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pipeline_id INTEGER,
			name TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			triggered_by TEXT,
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			logs TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (pipeline_id) REFERENCES pipelines(id) ON DELETE SET NULL
		)`,

		// Job steps table (for tracking individual steps in a job)
		`CREATE TABLE IF NOT EXISTS job_steps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			step_name TEXT NOT NULL,
			step_order INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			logs TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
		)`,

		// Create indexes
		`CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_pipeline_id ON jobs(pipeline_id)`,
		`CREATE INDEX IF NOT EXISTS idx_job_steps_job_id ON job_steps(job_id)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
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
