package utils

import (
	"node/utils/logger"
	"os"

	"gopkg.in/yaml.v2"
)

var Config map[string]interface{}

// LoadConfig loads the configuration from a YAML file into the global Config variable
func LoadConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		logger.Fatal("Failed to open config file: " + err.Error())
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		logger.Fatal("Failed to decode config file: " + err.Error())
	}
}

func GetConfig() map[interface{}]interface{} {
	if Config == nil {
		LoadConfig()
	}
	return Config["dbnode"].(map[interface{}]interface{})
}

func GetDbConfig() map[interface{}]interface{} {
	config := GetConfig()
	// logger.Debug(fmt.Sprintf("Config: %v", config))

	dbConfig, ok := config["database"].(map[interface{}]interface{})
	if !ok {
		logger.Fatal("Invalid database configuration")
	}
	return dbConfig
}
