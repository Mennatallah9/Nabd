package utils

import (
	"database/sql"
	"log"
	"nabd/models"

	_ "modernc.org/sqlite"
)

// InitDatabase initializes the SQLite database and creates tables
func InitDatabase(dbPath string) error {
	var err error
	models.DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err = models.DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to SQLite database")

	// Create tables
	if err = createTables(); err != nil {
		return err
	}

	return nil
}

// createTables creates the necessary database tables
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS container_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			container_id TEXT NOT NULL,
			name TEXT NOT NULL,
			cpu_percent REAL NOT NULL,
			memory_usage INTEGER NOT NULL,
			memory_limit INTEGER NOT NULL,
			network_rx INTEGER NOT NULL,
			network_tx INTEGER NOT NULL,
			status TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS autoheal_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			container_id TEXT NOT NULL,
			name TEXT NOT NULL,
			action TEXT NOT NULL,
			reason TEXT NOT NULL,
			success BOOLEAN NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			container_id TEXT NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			severity TEXT NOT NULL,
			active BOOLEAN NOT NULL DEFAULT 1,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := models.DB.Exec(query); err != nil {
			return err
		}
	}

	log.Println("Database tables created successfully")
	return nil
}