package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv initializes environment variables from a .env file.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// debugging, show all env variables
	log.Println("Loaded Environment Variables:")
	for _, env := range os.Environ() {
		log.Println(env)
	}
}

// GetEnv retrieves an environment variable with an optional fallback.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
