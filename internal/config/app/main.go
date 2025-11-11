package app

import (
	"strconv"

	"github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() Config {
	port, _ := strconv.Atoi(pkg.GetEnv("APP_PORT", "3000"))

	return Config{
		Port: port,
	}
}
