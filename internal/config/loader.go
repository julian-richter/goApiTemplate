package config

import (
	"github.com/julian-richter/ApiTemplate/internal/config/app"
	"github.com/julian-richter/ApiTemplate/internal/config/cache"
	"github.com/julian-richter/ApiTemplate/internal/config/database"
)

// Config represents the top-level configuration.
type Config struct {
	Cache    cache.Config
	Database database.Config
	App      app.Config
}

// Load initializes and returns the top-level configuration by aggregating
// cache, database, and application configurations.
func Load() Config {
	LoadEnv()

	return Config{
		Cache:    cache.Load(),
		Database: database.Load(),
		App:      app.Load(),
	}
}
