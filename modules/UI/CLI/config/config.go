package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the configuration keys.
type Config struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// LoadConfig reads the configuration file.
func LoadConfig() (Config, error) {
	configFilePath := getConfigFilePath()
	file, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config file does not exist")
		}
		return Config{}, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

// SaveConfig writes the configuration to a file.
func SaveConfig(config Config) error {
	configFilePath := getConfigFilePath()
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

// GetKeys retrieves configuration or prompts the user for input if not found.
func GetKeys() Config {
	config, err := LoadConfig()
	if err == nil {
		return config
	}

	// Prompt user for input if config file does not exist
	var privateKey, publicKey string
	fmt.Print("Enter your PRIVATE_KEY: ")
	privateKey = readInput()
	fmt.Print("Enter your PUBLIC_KEY: ")
	publicKey = readInput()

	// Save keys to the config file
	config = Config{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		os.Exit(1)
	}
	return config
}

// Helper to get the configuration file path
func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".libr")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}
	}

	return filepath.Join(configDir, "config.json")
}

// Helper to read input from the user
func readInput() string {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}
