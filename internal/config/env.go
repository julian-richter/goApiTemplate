package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads the env file if it exists.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] no env file found — using system environment variables")
	}
}
