package pkg

import "os"

// GetEnv returns the value of the environment variable with the given key or fallback if unset.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
