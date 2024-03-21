package common

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ClientLogger struct {
	Std        zerolog.Logger
	Compliance zerolog.Logger
}

func DefaultLoggers() ClientLogger {
	return ClientLogger{
		Std:        log.Logger,
		Compliance: log.Logger,
	}
}
