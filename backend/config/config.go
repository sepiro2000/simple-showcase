package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port: getEnvOrDefault("PORT", "8080"),
	}

	// Required environment variables
	if cfg.DBHost = os.Getenv("DB_HOST"); cfg.DBHost == "" {
		return nil, fmt.Errorf("DB_HOST environment variable is required")
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

	return cfg, nil
}

// getEnvOrDefault returns the value of the environment variable or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
