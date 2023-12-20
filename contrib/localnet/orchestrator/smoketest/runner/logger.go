package runner

import (
	"fmt"
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
	text := fmt.Sprintf(message, args...)
	ml.logger.Printf(text + "\n")
}

// Info prints a message to the logger if verbose is true
func (ml *Logger) Info(message string, args ...interface{}) {
	if ml.verbose {
		text := fmt.Sprintf(message, args...)
		ml.logger.Printf("[INFO]" + text + "\n")
	}
}

// InfoLoud prints a message to the logger if verbose is true
func (ml *Logger) InfoLoud(message string, args ...interface{}) {
	if ml.verbose {
		text := fmt.Sprintf(message, args...)
		ml.logger.Printf("[INFO] =======================================")
		ml.logger.Printf("[INFO]" + text + "\n")
		ml.logger.Printf("[INFO] =======================================")
	}
}

// Error prints an error message to the logger
func (ml *Logger) Error(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	ml.logger.Printf("[ERROR]" + text + "\n")
}
