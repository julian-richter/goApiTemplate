package app

import (
	"fmt"
	"strconv"

	"github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() (Config, error) {
	port := env.GetEnv("APP_PORT", "8080")

	portNumber, err := strconv.Atoi(env.GetEnv("APP_PORT", ""))
	if err != nil {
		return Config{}, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	if portNumber < 0 || portNumber > 65535 {
		return Config{}, fmt.Errorf("APP_PORT must be non-negative or above 65535, got %d", portNumber)
	}
	return Config{
		Port: port,
	}, nil
}
