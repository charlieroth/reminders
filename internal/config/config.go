package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerPort  string
	LogLevel    string
	DatabaseURL string
	JwtSecret   string
	Salt        string
}

func NewConfig() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	logLevel := os.Getenv("LOG_LEVEL")
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	salt := os.Getenv("SALT")

	if serverPort == "" || databaseURL == "" || jwtSecret == "" || salt == "" {
		return nil, errors.New("SERVER_PORT, DATABASE_URL, JWT_SECRET, and SALT must be set")
	}

	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		ServerPort:  serverPort,
		LogLevel:    logLevel,
		DatabaseURL: databaseURL,
		JwtSecret:   jwtSecret,
		Salt:        salt,
	}, nil
}
