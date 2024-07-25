// Package bg provides primitives for the background tasks
package bg

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog"
)

type config struct {
	name     string
	logger   zerolog.Logger
	callback context.CancelFunc
}

type Opt func(*config)

func WithName(name string) Opt {
	return func(cfg *config) { cfg.name = name }
}

func WithCallback(cancel context.CancelFunc) Opt {
	return func(cfg *config) { cfg.callback = cancel }
}

func WithLogger(logger zerolog.Logger) Opt {
	return func(cfg *config) { cfg.logger = logger }
}

// Work emits a new task in the background
func Work(ctx context.Context, f func(context.Context) error, opts ...Opt) {
	cfg := config{
		name:     "",
		logger:   zerolog.Nop(),
		callback: nil,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("recovered from PANIC in background task: %v", r)
				logError(err, cfg, true)
			}
		}()

		if err := f(ctx); err != nil {
			logError(err, cfg, false)
		}
		// Use cancel function if it is provided.
		// This is used for stopping the main thread based on the outcome of the background task.
		if cfg.callback != nil {
			cfg.logger.Info().Msgf("background task completed for %s", cfg.name)
			cfg.callback()
		}
	}()
}

func logError(err error, cfg config, isPanic bool) {
	if err == nil {
		return
	}

	evt := cfg.logger.Error().Err(err)

	// print stack trace when a panic occurs
	if isPanic {
		buf := make([]byte, 1024)
		for {
			n := runtime.Stack(buf, false)
			if n < len(buf) {
				buf = buf[:n]
				break
			}
			buf = make([]byte, 2*len(buf))
		}

		evt.Bytes("stack_trace", buf)
	}

	name := cfg.name
	if name == "" {
		name = "unknown"
	}

	evt.Str("worker.name", name).Msg("Background task failed")
}
