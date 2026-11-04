package utils

import (
	"fmt"
	"os"
	"regexp"

	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
	"gopkg.in/yaml.v2"

	"github.com/DALDA-IITJ/libr/modules/core/config"
)

var Config map[string]interface{}

// LoadConfig loads the configuration from a YAML file into the global Config variable
func LoadConfigAndEnv() {
	config.LoadEnv()

	content, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Fatal("Failed to open config file: " + err.Error())
	}

	re := regexp.MustCompile(`\$\{([^}:]+)(?::-([^}]*))?\}`)
	expanded := re.ReplaceAllStringFunc(string(content), func(match string) string {
		sub := re.FindStringSubmatch(match)
		varName := sub[1]
		defaultVal := ""
		if len(sub) > 2 {
			defaultVal = sub[2]
		}
		if val, exists := os.LookupEnv(varName); exists && val != "" {
			return val
		}
		return defaultVal
	})

	err = yaml.Unmarshal([]byte(expanded), &Config)
	if err != nil {
		logger.Fatal("Failed to decode config file: " + err.Error())
	}
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
