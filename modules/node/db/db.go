package db

import (
	"database/sql"
	"fmt"

	"github.com/DALDA-IITJ/libr/modules/node/utils"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitDB() {
	var err error

	EnsureDatabaseSetup()

	dbConfig := utils.GetDbConfig()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig["host"], dbConfig["port"], dbConfig["user"], dbConfig["password"], dbConfig["dbname"], dbConfig["sslmode"])

	logger.Debug("Connecting to database with connection string: " + connStr)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}

	_, err = DB.Exec(`
CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    sender TEXT NOT NULL,
    content TEXT NOT NULL,
    bucket_id BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS bucket_index ON message (bucket_id);
	`)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Database initialization failed: %v", err))
	}

	logger.Info("Database initialized successfully")
}

func EnsureDatabaseSetup() {
	logger.Info("Ensuring database setup...")

	dbConfig := utils.GetDbConfig()

	// Use "postgres" as the initial connection DB
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		dbConfig["host"], dbConfig["port"], dbConfig["user"], dbConfig["password"], dbConfig["sslmode"])

	logger.Debug("Connecting to PostgreSQL with connection string: " + connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}
	defer db.Close()

	// Ensure PostgreSQL is running
	err = db.Ping()
	if err != nil {
		logger.Fatal(fmt.Sprintf("PostgreSQL is not running or cannot be reached: %v", err))
	}

	// Check if the target database exists
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)`, dbConfig["dbname"]).Scan(&exists)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to check if database exists: %v", err))
	}

	if !exists {
		logger.Info(fmt.Sprintf("Database '%s' does not exist. Creating database...", dbConfig["dbname"]))
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbConfig["dbname"]))
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to create database: %v", err))
		}
		logger.Info("Database created successfully.")
	} else {
		logger.Info("Database already exists.")
	}
}

func SaveToDB(sender string, content string, bucketId int64) (int64, error) {
	logger.Debug(fmt.Sprintf("Saving message to database: sender=%s, content=%s, bucketId=%d", sender, content, bucketId))

	var id int64
	err := DB.QueryRow("INSERT INTO message (sender, content, bucket_id) VALUES ($1, $2, $3) RETURNING id", sender, content, bucketId).Scan(&id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to save message: %v", err))
		return 0, err
	}

	logger.Info(fmt.Sprintf("Message saved with ID: %d", id))
	return id, nil
}

func GetFromDB(bucketId int64) (interface{}, error) {
	logger.Debug(fmt.Sprintf("Retrieving messages from database for bucketId=%d", bucketId))

	rows, err := DB.Query("SELECT sender, content FROM message WHERE bucket_id = $1", bucketId)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to retrieve messages: %v", err))
		return nil, err
	}
	defer rows.Close()

	messages := make([]map[string]string, 0)
	for rows.Next() {
		var sender, content string
		if err := rows.Scan(&sender, &content); err != nil {
			logger.Error(fmt.Sprintf("Failed to scan row: %v", err))
			return nil, err
		}
		messages = append(messages, map[string]string{"sender": sender, "content": content})
	}

	if err := rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Row iteration error: %v", err))
		return nil, err
	}

	logger.Info(fmt.Sprintf("Retrieved %d messages for bucketId=%d", len(messages), bucketId))
	return messages, nil
}
