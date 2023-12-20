package runner

import (
	"log"
	"os"
)

// Logger is a wrapper around log.Logger that adds verbosity
type Logger struct {
	verbose bool
	logger  *log.Logger
}

// NewLogger creates a new Logger
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		logger:  log.New(os.Stdout, "", 0),
	}
}

// Print prints a message to the logger
func (ml *Logger) Print(message string, args ...interface{}) {
	ml.logger.Printf(message+"\n", args...)
}

// Info prints a message to the logger if verbose is true
func (ml *Logger) Info(message string, args ...interface{}) {
	if ml.verbose {
		ml.logger.Printf("[INFO]"+message+"\n", args)
	}
}

// InfoLoud prints a message to the logger if verbose is true
func (ml *Logger) InfoLoud(message string, args ...interface{}) {
	if ml.verbose {
		ml.logger.Printf("[INFO] =======================================")
		ml.logger.Printf("[INFO]"+message+"\n", message, args)
		ml.logger.Printf("[INFO] =======================================")
	}
}

// Error prints an error message to the logger
func (ml *Logger) Error(message string, args ...interface{}) {
	ml.logger.Printf("[ERROR]"+message+"\n", args)
}
