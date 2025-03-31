package db

import (
	"database/sql"
	"fmt"
	logger "libr/core"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "host=localhost user=admin dbname=librdb sslmode=disable")
	if err != nil {
		logger.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		sender TEXT NOT NULL,
		content TEXT NOT NULL,
		bucket_id BIGINT NOT NULL,
		INDEX bucket_index (bucket_id)
	)
	`)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Database initialization failed: %v", err))
	}

	logger.Info("Database initialized successfully")
}
