package config

import "os"

// EnvOrDefault returns the configured environment value or fallback when empty.
func EnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
