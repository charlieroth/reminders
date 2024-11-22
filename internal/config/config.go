package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerPort  string
	LogLevel    string
	DatabaseURL string
}

func NewConfig() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	logLevel := os.Getenv("LOG_LEVEL")
	databaseURL := os.Getenv("DATABASE_URL")

	if serverPort == "" || databaseURL == "" {
		return nil, errors.New("SERVER_PORT and DATABASE_URL must be set")
	}

	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		ServerPort:  serverPort,
		LogLevel:    logLevel,
		DatabaseURL: databaseURL,
	}, nil
}
