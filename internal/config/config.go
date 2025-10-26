// Package config handles application configuration management.
// It loads settings from environment variables with sensible defaults.
package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}

	return &Config{
		Port: port,
	}
}
