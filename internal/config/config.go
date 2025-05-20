package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config holds the application configuration
type Config struct {
	DatabaseUrl string `yaml:"db_url"`
	Port        string `yaml:"port"`
}

var appConfig Config
var onceAppConfig sync.Once

// GetEnv returns the value of the environment variable if set, otherwise returns the default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetAppConfig loads and returns the application configuration
func GetAppConfig(filepath string) *Config {
	// Singleton pattern to ensure config is loaded only once
	onceAppConfig.Do(func() {
		var err error
		rawConfig, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(rawConfig, &appConfig)
		if err != nil {
			panic(err)
		}

		// Override config entries with environment variables if they are set
		appConfig.DatabaseUrl = GetEnv("VINYLVAULT_DATABASE_URL", appConfig.DatabaseUrl)
		appConfig.Port = GetEnv("VINYLVAULT_PORT", appConfig.Port)
	})

	return &appConfig
}
