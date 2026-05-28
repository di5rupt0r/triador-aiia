package config

import (
	"errors"
	"os"
)

type Config struct {
	OpenAIAPIKey string
	DatabasePath string
	ServerPort   string
}

func Load() (*Config, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is required")
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./triador-aiia.db"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		OpenAIAPIKey: apiKey,
		DatabasePath: dbPath,
		ServerPort:   port,
	}, nil
}
