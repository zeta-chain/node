package common

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ClientLogger is a struct that contains the logger for a chain client
type ClientLogger struct {
	Std        zerolog.Logger
	Compliance zerolog.Logger
}

// DefaultLoggers returns the default loggers for a chain client
func DefaultLoggers() ClientLogger {
	return ClientLogger{
		Std:        log.Logger,
		Compliance: log.Logger,
	}
}
