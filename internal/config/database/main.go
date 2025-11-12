package database

import (
	"fmt"
	"strconv"

	env "github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() (Config, error) {
	portStr := env.GetEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_PORT value %q: %w", portStr, err)
	}

	// Check if the port is valid
	if port < 0 || port > 65535 {
		return Config{}, fmt.Errorf("DB_PORT must be non-negative or above 65535, got %d", port)
	}

	return Config{
		Host:     env.GetEnv("DB_HOST", "127.0.0.1"),
		Port:     port,
		User:     env.GetEnv("DB_USER", "postgres"),
		Password: env.GetEnv("DB_PASSWORD", "password"),
		Name:     env.GetEnv("DB_NAME", "postgres"),
		SSLMode:  env.GetEnv("DB_SSL_MODE", "disable"),
	}, nil
}
