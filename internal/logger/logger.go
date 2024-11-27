package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(environment string) *zerolog.Logger {
	var logger zerolog.Logger
	if environment == "PRODUCTION" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	return &logger
}
