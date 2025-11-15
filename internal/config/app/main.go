package app

import (
	"fmt"
	"strconv"
	"strings"

	env "github.com/julian-richter/ApiTemplate/pkg"
)

// Load initializes a Config struct by fetching environment variables
// with fallbacks to default values.
func Load() (Config, error) {
	port := strings.TrimSpace(env.GetEnv("APP_PORT", "8080"))
	appName := strings.TrimSpace(env.GetEnv("APP_NAME", "ApiTemplate"))
	environment := Env(strings.TrimSpace(env.GetEnv("APP_ENV", "dev")))

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return Config{}, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	if portNumber < 0 || portNumber > 65535 {
		return Config{}, fmt.Errorf("APP_PORT must be between 0 and 65535, got %d", portNumber)
	}

	// Validate environment
	if !environment.Valid() {
		return Config{}, fmt.Errorf("invalid APP_ENV: %q (valid: %q or %q)", environment, EnvProd, EnvDev)
	}

	return Config{
		Env:             environment,
		ApplicationName: appName,
		Port:            port,
	}, nil
}
