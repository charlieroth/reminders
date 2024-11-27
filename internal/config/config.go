package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Environment          string
	ServerPort           string
	LogLevel             string
	DatabaseURL          string
	JwtSecret            string
	Salt                 string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewConfig() (*Config, error) {
	environment := os.Getenv("ENVIRONMENT")
	serverPort := os.Getenv("SERVER_PORT")
	logLevel := os.Getenv("LOG_LEVEL")
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	salt := os.Getenv("SALT")
	accessTokenDurationMins := os.Getenv("ACCESS_TOKEN_DURATION_MINS")
	refreshTokenDurationMins := os.Getenv("REFRESH_TOKEN_DURATION_MINS")

	if environment == "" {
		return nil, errors.New("ENVIRONMENT must be set")
	}

	if serverPort == "" || databaseURL == "" || jwtSecret == "" || salt == "" {
		return nil, errors.New("SERVER_PORT, DATABASE_URL, JWT_SECRET, and SALT must be set")
	}

	var accessTokenDuration time.Duration
	if accessTokenDurationMins == "" {
		accessTokenDuration = time.Minute * 15
	} else {
		accessTokenDurationParsed, err := time.ParseDuration(accessTokenDurationMins)
		if err != nil {
			return nil, fmt.Errorf("invalid time duration for access token: %w", err)
		}
		accessTokenDuration = accessTokenDurationParsed
	}

	var refreshTokenDuration time.Duration
	if refreshTokenDurationMins == "" {
		refreshTokenDuration = time.Minute * 60
	} else {
		refreshTokenDurationParsed, err := time.ParseDuration(refreshTokenDurationMins)
		if err != nil {
			return nil, fmt.Errorf("invalid time duration for refresh token: %w", err)
		}
		refreshTokenDuration = refreshTokenDurationParsed
	}

	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		Environment:          environment,
		ServerPort:           serverPort,
		LogLevel:             logLevel,
		DatabaseURL:          databaseURL,
		JwtSecret:            jwtSecret,
		Salt:                 salt,
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	}, nil
}
