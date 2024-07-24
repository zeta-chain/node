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
	cancel context.CancelCauseFunc
}

type Opt func(*config)

func WithName(name string) Opt {
	return func(cfg *config) { cfg.name = name }
}

func WithCancel(cancel context.CancelCauseFunc) Opt {
	return func(cfg *config) { cfg.cancel = cancel }
}

func WithLogger(logger zerolog.Logger) Opt {
	return func(cfg *config) { cfg.logger = logger }
}

// Work emits a new task in the background
func Work(ctx context.Context, f func(context.Context) error, opts ...Opt) {
	cfg := config{
		name:   "",
		logger: zerolog.Nop(),
		cancel: nil,
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

		err := f(ctx)
		if err != nil {
			logError(err, cfg)
		}
		// Use cancel function if it is provided.
		// This is used for stopping the main thread based on the outcome of the background task.
		if cfg.cancel != nil {
			cfg.cancel(fmt.Errorf("cancel function triggered by %s", cfg.name))
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
