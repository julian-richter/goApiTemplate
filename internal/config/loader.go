package config

import (
	"fmt"

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
		return Config{}, fmt.Errorf("failed to load cache config: %w", err)
	}

	dbCfg, err := database.Load()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load database config: %w", err)
	}

	appCfg, err := app.Load()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load application config: %w", err)
	}

	return Config{
		Cache:    cacheCfg,
		Database: dbCfg,
		App:      appCfg,
	}, nil
}
