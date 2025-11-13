package cache

import (
	"fmt"
	"strconv"
	"strings"

	env "github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() (Config, error) {
	port, err := strconv.Atoi(env.GetEnv("VALKEY_PORT", "6379"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid VALKEY_PORT: %w", err)
	}

	db, err := strconv.Atoi(env.GetEnv("VALKEY_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid VALKEY_DB: %w", err)
	}

	if db < 0 || db > 15 {
		return Config{}, fmt.Errorf("VALKEY_DB must be between 0 and 15 (standard Redis range), got %d", db)
	}

	return Config{
		User:     strings.TrimSpace(env.GetEnv("VALKEY_USER", "")),
		Host:     strings.TrimSpace(env.GetEnv("VALKEY_HOST", "localhost")),
		Port:     port,
		Password: strings.TrimSpace(env.GetEnv("VALKEY_PASSWORD", "")),
		DB:       db,
	}, nil
}
