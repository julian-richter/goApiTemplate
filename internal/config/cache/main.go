package cache

import (
	"strconv"

	"github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() Config {
	port, _ := strconv.Atoi(pkg.GetEnv("CACHE_PORT", "6379"))
	db, _ := strconv.Atoi(pkg.GetEnv("CACHE_DB", "0"))

	return Config{
		Host:     pkg.GetEnv("CACHE_HOST", "localhost"),
		Port:     port,
		Password: pkg.GetEnv("CACHE_PASSWORD", ""),
		DB:       db,
	}
}
