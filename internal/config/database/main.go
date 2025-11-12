package database

import (
	"fmt"
	"strconv"

	"github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() (Config, error) {
	port, err := strconv.Atoi(env.GetEnv("DB_PORT", "5432"))
	if err != nil {
		return Config{}, err
	}

	// Check if port is valid
	if port < 0 {
		return Config{}, fmt.Errorf("DB_PORT must be non-negative, got %d", port)
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
