package config

import (
	"log"
	"os"

	"github.com/DALDA-IITJ/libr/modules/core/crypto"
	"github.com/joho/godotenv"
)

// LoadEnv initializes environment variables from a .env file.
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	publicKey := os.Getenv("PUBLIC_KEY")

	if privateKey == "" || publicKey == "" {
		log.Println("Private/Public key not found in .env, generating new keys...")
		err := crypto.GenerateKeys()
		if err != nil {
			log.Fatalf("Failed to generate keys: %v", err)
		}
	}

}

// GetEnv retrieves an environment variable with an optional fallback.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
