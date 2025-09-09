package base

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zeta-chain/node/zetaclient/config"
)

const complianceLogFile = "compliance.log"

// Logger contains the base loggers
type Logger struct {
	Std        zerolog.Logger
	Compliance zerolog.Logger
}

// DefaultLogger creates default base loggers for tests
func DefaultLogger() Logger {
	return Logger{
		Std:        log.Logger,
		Compliance: log.Logger,
	}
}

// ObserverLogger contains the loggers for chain observers
type ObserverLogger struct {
	// the parent logger for the chain observer
	Chain zerolog.Logger

	// the logger for inbound transactions
	Inbound zerolog.Logger

	// the logger for outbound transactions
	Outbound zerolog.Logger

	// the logger for the compliance check
	Compliance zerolog.Logger
}

// NewLogger initializes the base loggers
func NewLogger(cfg config.Config) (Logger, error) {
	// open compliance log file
	complianceFile, err := openComplianceLogFile(cfg)
	if err != nil {
		return Logger{}, err
	}

	augmentLogger := func(logger zerolog.Logger) zerolog.Logger {
		level := zerolog.Level(cfg.LogLevel)

		return logger.Level(level).With().Timestamp().Logger()
	}

	// create loggers based on configured level and format
	var stdWriter io.Writer = os.Stdout
	if cfg.LogFormat != "json" {
		stdWriter = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	std := augmentLogger(zerolog.New(stdWriter))
	compliance := augmentLogger(zerolog.New(complianceFile))

	if cfg.LogSampler {
		std = std.Sample(&zerolog.BasicSampler{N: 5})
	}

	// set global logger
	log.Logger = std

	return Logger{
		Std:        std,
		Compliance: compliance,
	}, nil
}

// openComplianceLogFile opens the compliance log file
func openComplianceLogFile(cfg config.Config) (*os.File, error) {
	// use zetacore home as default
	logPath := cfg.ZetaCoreHome
	if cfg.ComplianceConfig.LogPath != "" {
		logPath = cfg.ComplianceConfig.LogPath
	}

	// clean file name
	name := filepath.Join(logPath, complianceLogFile)
	name, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	name = filepath.Clean(name)

	// open (or create) compliance log file
	return os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
