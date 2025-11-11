package database

import (
	"strconv"

	"github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables with fallbacks to default values.
func Load() Config {
	port, _ := strconv.Atoi(pkg.GetEnv("DB_PORT", "5432"))

	return Config{
		Host:     pkg.GetEnv("DB_HOST", "127.0.0.1"),
		Port:     port,
		User:     pkg.GetEnv("DB_USER", "postgres"),
		Password: pkg.GetEnv("DB_PASSWORD", "password"),
		Name:     pkg.GetEnv("DB_NAME", "postgres"),
		SSLMode:  pkg.GetEnv("DB_SSL_MODE", "disable"),
	}
}
