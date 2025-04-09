package utils

import (
	"fmt"
	"os"

	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
	"gopkg.in/yaml.v2"

	"github.com/DALDA-IITJ/libr/modules/core/config"
)

var Config map[string]interface{}

// LoadConfig loads the configuration from a YAML file into the global Config variable
func LoadConfigAndEnv() {
	file, err := os.Open("config.yaml")
	logger.Debug("hello")
	if err != nil {
		logger.Error("Failed to open config file: " + err.Error())
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		logger.Fatal("Failed to decode config file: " + err.Error())
	}

	config.LoadEnv()

}

func GetConfig() map[interface{}]interface{} {
	if Config == nil {
		LoadConfigAndEnv()
	}
	return Config["dbnode"].(map[interface{}]interface{})
}

func GetDbConfig() map[interface{}]interface{} {
	config := GetConfig()
	logger.Debug(fmt.Sprintf("Config: %v", config))

	dbConfig, ok := config["database"].(map[interface{}]interface{})
	if !ok {
		logger.Fatal("Invalid database configuration")
	}
	return dbConfig
}

func GetTransactionConfig() map[interface{}]interface{} {
	config := GetConfig()
	// logger.Debug(fmt.Sprintf("Config: %v", config))

	transactionConfig, ok := config["transaction"].(map[interface{}]interface{})
	if !ok {
		logger.Fatal("Invalid transaction configuration")
	}
	return transactionConfig
}

func GetNodeConfig() map[interface{}]interface{} {
	config := GetConfig()
	// logger.Debug(fmt.Sprintf("Config: %v", config))

	nodeConfig, ok := config["node"].(map[interface{}]interface{})
	if !ok {
		logger.Fatal("Invalid node configuration")
	}
	return nodeConfig
}
