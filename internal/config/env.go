package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file if appropriate.
// In production mode, missing .env files are ignored.
// In development mode, missing .env files trigger a warning.
func LoadEnv() {
	appEnv := os.Getenv("APP_ENV")

	err := godotenv.Load()
	if err == nil {
		log.Println("[config] environment variables loaded from .env file")
		return
	}

	if errors.Is(err, os.ErrNotExist) {
		if appEnv == "production" || appEnv == "prod" {
			log.Println("[config] production mode — skipping .env loading (using system environment variables)")
		} else {
			log.Println("[config] warning: no .env file found — using system environment variables")
		}
		return
	}

	// Any other error (like permission denied, parse failure, etc.)
	log.Fatalf("[config] failed to load .env file: %v", err)
}
