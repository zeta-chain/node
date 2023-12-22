package runner

import (
	"fmt"
	"github.com/fatih/color"
)

// Logger is a wrapper around log.Logger that adds verbosity
type Logger struct {
	verbose bool
	logger  *color.Color
	prefix  string
}

// NewLogger creates a new Logger
func NewLogger(verbose bool, printColor color.Attribute, prefix string) *Logger {
	return &Logger{
		verbose: verbose,
		logger:  color.New(printColor),
		prefix:  prefix,
	}
}

// Print prints a message to the logger
func (l *Logger) Print(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	l.logger.Printf(l.prefix + " - " + text + "\n")
}

// Info prints a message to the logger if verbose is true
func (l *Logger) Info(message string, args ...interface{}) {
	if l.verbose {
		text := fmt.Sprintf(message, args...)
		l.logger.Printf(l.prefix + " - " + "[INFO]" + text + "\n")
	}
}

// InfoLoud prints a message to the logger if verbose is true
func (l *Logger) InfoLoud(message string, args ...interface{}) {
	if l.verbose {
		text := fmt.Sprintf(message, args...)
		l.logger.Printf(l.prefix + " - " + "[INFO] =======================================")
		l.logger.Printf(l.prefix + " - " + "[INFO]" + text + "\n")
		l.logger.Printf(l.prefix + " - " + "[INFO] =======================================")
	}
}

// Error prints an error message to the logger
func (l *Logger) Error(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	l.logger.Printf(l.prefix + " - " + "[ERROR]" + text + "\n")
}
