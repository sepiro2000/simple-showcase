package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port          string
	WriteDBHost   string
	ReadDBHost    string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisHost     string
	RedisPassword string
	RedisDB       int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port: getEnvOrDefault("PORT", "8080"),
	}

	// Required environment variables for database
	if cfg.WriteDBHost = os.Getenv("WRITE_DB_HOST"); cfg.WriteDBHost == "" {
		return nil, fmt.Errorf("WRITE_DB_HOST environment variable is required")
	}
	if cfg.ReadDBHost = os.Getenv("READ_DB_HOST"); cfg.ReadDBHost == "" {
		return nil, fmt.Errorf("READ_DB_HOST environment variable is required")
	}
	if cfg.DBPort = os.Getenv("DB_PORT"); cfg.DBPort == "" {
		return nil, fmt.Errorf("DB_PORT environment variable is required")
	}
	if cfg.DBUser = os.Getenv("DB_USER"); cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER environment variable is required")
	}
	if cfg.DBPassword = os.Getenv("DB_PASSWORD"); cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if cfg.DBName = os.Getenv("DB_NAME"); cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME environment variable is required")
	}

	// Optional Redis configuration
	cfg.RedisHost = os.Getenv("REDIS_HOST")
	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		var err error
		cfg.RedisDB, err = strconv.Atoi(redisDB)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB value: %v", err)
		}
	}

	return cfg, nil
}

// getEnvOrDefault returns the value of the environment variable or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
