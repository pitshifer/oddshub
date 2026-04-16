package config

import (
	"fmt"
	"os"
)

type Config struct {
	TheOddsApiKey string
}

func LoadConfig() (*Config, error) {
	theOddsApiKey := os.Getenv("THE_ODDS_API_KEY")
	if theOddsApiKey == "" {
		return nil, fmt.Errorf("THE_ODDS_API_KEY environment variable is required")
	}

	return &Config{
		TheOddsApiKey: theOddsApiKey,
	}, nil
}
