package config

import (
	"github.com/joeshaw/envdecode"
)

// Config holds all configuration for the application.
// It is populated from environment variables using struct tags.
type Config struct {
	AppBaseURL    string `env:"APP_BASE_URL,default=http://localhost:8080"`
	ServerAddress string `env:"SERVER_ADDRESS,default=:8080"`
	DBPath        string `env:"DB_PATH,default=./data/urls.db"`
}

// New creates and loads a new Config object from environment variables.
// This function acts as a constructor for the Config struct.
func New() (*Config, error) {
	var c Config
	if err := envdecode.Decode(&c); err != nil {
		// If decoding fails, return the zero value for Config and the error.
		return nil, err
	}
	// Return the populated Config struct.
	return &c, nil
}
