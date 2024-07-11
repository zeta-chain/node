package base

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	ComplianceLogFile = "compliance.log"
)

// Logger contains the base loggers
type Logger struct {
	Std        zerolog.Logger
	Compliance zerolog.Logger
}

// DefaultLoggers creates default base loggers for tests
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

	// the logger for the chain's gas price
	GasPrice zerolog.Logger

	// the logger for block headers
	Headers zerolog.Logger

	// the logger for the compliance check
	Compliance zerolog.Logger
}

// InitLogger initializes the base loggers
func InitLogger(cfg config.Config) (Logger, error) {
	// open compliance log file
	file, err := openComplianceLogFile(cfg)
	if err != nil {
		return DefaultLogger(), err
	}

	level := zerolog.Level(cfg.LogLevel)

	// create loggers based on configured level and format
	var std zerolog.Logger
	var compliance zerolog.Logger
	switch cfg.LogFormat {
	case "json":
		std = zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
		compliance = zerolog.New(file).Level(level).With().Timestamp().Logger()
	case "text":
		std = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(zerolog.Level(cfg.LogLevel)).
			With().
			Timestamp().
			Logger()
		compliance = zerolog.New(file).Level(level).With().Timestamp().Logger()
	default:
		std = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		compliance = zerolog.New(file).With().Timestamp().Logger()
	}

	if cfg.LogSampler {
		std = std.Sample(&zerolog.BasicSampler{N: 5})
	}
	log.Logger = std // set global logger

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
	name := filepath.Join(logPath, ComplianceLogFile)
	name, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	name = filepath.Clean(name)

	// open (or create) compliance log file
	return os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
