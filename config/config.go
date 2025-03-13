package config

import (
	"os"
	"strings"

	"loan-billing-system/internal/db"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DB     db.Config
	Server struct {
		Port string
	}
}

// Load loads the application configuration from environment variables
func Load() (*Config, error) {
	// Load from .env file if it exists
	_ = godotenv.Load()

	config := &Config{}

	// Database configuration
	config.DB.Host = getEnv("DB_HOST", "localhost")
	config.DB.Port = getEnv("DB_PORT", "5431")
	config.DB.User = getEnv("DB_USER", "postgres")
	config.DB.Password = getEnv("DB_PASSWORD", "")
	config.DB.DBName = getEnv("DB_NAME", "loan_billing_system")
	config.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Server configuration
	config.Server.Port = getEnv("SERVER_PORT", "8080")

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}
