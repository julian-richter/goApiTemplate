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
func Load() (Config, error) {
	// Load environment variables (optional env file)
	LoadEnv()

	cacheCfg, err := cache.Load()
	if err != nil {
		return Config{}, err
	}

	dbCfg, err := database.Load()
	if err != nil {
		return Config{}, err
	}

	appCfg, err := app.Load()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Cache:    cacheCfg,
		Database: dbCfg,
		App:      appCfg,
	}, nil
}
