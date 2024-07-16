// Package bg provides primitives for the background tasks
package bg

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog"
)

type config struct {
	name   string
	logger zerolog.Logger
}

type Opt func(*config)

func WithName(name string) Opt {
	return func(cfg *config) { cfg.name = name }
}

func WithLogger(logger zerolog.Logger) Opt {
	return func(cfg *config) { cfg.logger = logger }
}

// Work emits a new task in the background
func Work(ctx context.Context, f func(context.Context) error, opts ...Opt) {
	cfg := config{
		name:   "",
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("recovered from PANIC in background task: %v", r)
				logError(err, cfg)
				printStack()
			}
		}()

		if err := f(ctx); err != nil {
			logError(err, cfg)
		}
	}()
}

func logError(err error, cfg config) {
	if err == nil {
		return
	}

	name := cfg.name
	if name == "" {
		name = "unknown"
	}

	cfg.logger.Error().Err(err).Str("worker.name", name).Msgf("Background task failed")
}

// printStack prints the stack trace when a panic occurs
func printStack() {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	fmt.Printf("Stack trace:\n%s\n", string(buf))
}
