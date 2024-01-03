package runner

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	loggerSeparator = " | "
	padding         = 10
)

// Logger is a wrapper around log.Logger that adds verbosity
type Logger struct {
	verbose bool
	logger  *color.Color
	prefix  string
}

// NewLogger creates a new Logger
func NewLogger(verbose bool, printColor color.Attribute, prefix string) *Logger {
	// trim prefix to padding
	if len(prefix) > padding {
		prefix = prefix[:padding]
	}

	return &Logger{
		verbose: verbose,
		logger:  color.New(printColor),
		prefix:  prefix,
	}
}

// Print prints a message to the logger
func (l *Logger) Print(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + text + "\n")
}

// Info prints a message to the logger if verbose is true
func (l *Logger) Info(message string, args ...interface{}) {
	if l.verbose {
		text := fmt.Sprintf(message, args...)
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO]" + text + "\n")
	}
}

// InfoLoud prints a message to the logger if verbose is true
func (l *Logger) InfoLoud(message string, args ...interface{}) {
	if l.verbose {
		text := fmt.Sprintf(message, args...)
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO] =======================================")
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO]" + text + "\n")
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO] =======================================")
	}
}

// Error prints an error message to the logger
func (l *Logger) Error(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[ERROR]" + text + "\n")
}

func (l *Logger) getPrefixWithPadding() string {
	// add padding to prefix
	prefix := l.prefix
	for i := len(l.prefix); i < padding; i++ {
		prefix += " "
	}
	return prefix
}
