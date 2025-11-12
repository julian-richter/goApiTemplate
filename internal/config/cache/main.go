package cache

import (
	"fmt"
	"strconv"

	"github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() (Config, error) {
	port, err := strconv.Atoi(env.GetEnv("CACHE_PORT", "6379"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid CACHE_PORT: %w", err)
	}

	db, err := strconv.Atoi(env.GetEnv("CACHE_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid CACHE_DB: %w", err)
	}

	if db < 0 {
		return Config{}, fmt.Errorf("[config] CACHE_DB must be non-negative, got %d", db)
	}

	return Config{
		Host:     env.GetEnv("CACHE_HOST", "localhost"),
		Port:     port,
		Password: env.GetEnv("CACHE_PASSWORD", ""),
		DB:       db,
	}, nil
}
