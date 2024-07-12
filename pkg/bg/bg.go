// Package bg provides primitives for the background tasks
package bg

import (
	"context"
	"fmt"

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
