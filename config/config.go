package config

import (
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	DatabaseURL string
	ServerPort  int
	Environment string
}

// Load loads configuration from environment variables.
func Load() *Config {
	port := 8080
	if p := os.Getenv("SERVER_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "user=sqlc dbname=sqlc_db sslmode=disable host=localhost"),
		ServerPort:  port,
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv retrieves an environment variable with a default fallback.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
