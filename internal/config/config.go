package config

import (
	"fmt"
	"os"
)

type Config struct {
	TheOddsApiKey string
	DatabaseURL   string
}

func LoadConfig() (*Config, error) {
	theOddsApiKey := os.Getenv("THE_ODDS_API_KEY")
	if theOddsApiKey == "" {
		return nil, fmt.Errorf("THE_ODDS_API_KEY environment variable is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	return &Config{
		TheOddsApiKey: theOddsApiKey,
		DatabaseURL:   databaseURL,
	}, nil
}
