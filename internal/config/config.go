package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	LogLevel    string
}

func NewConfig() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	logLevel := os.Getenv("LOG_LEVEL")

	if serverPort == "" || databaseURL == "" {
		return nil, errors.New("SERVER_PORT and DATABASE_URL must be set")
	}

	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		ServerPort:  serverPort,
		DatabaseURL: databaseURL,
		LogLevel:    logLevel,
	}, nil
}
