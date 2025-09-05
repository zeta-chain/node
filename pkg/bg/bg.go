// Package bg provides primitives for the background tasks
package bg

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog"
)

type config struct {
	name       string
	logger     zerolog.Logger
	onComplete func()
}

type Opt func(*config)

func WithName(name string) Opt {
	return func(cfg *config) { cfg.name = name }
}

// OnComplete is a callback function that is called
// when the background task is completed without an error
func OnComplete(fn func()) Opt {
	return func(cfg *config) { cfg.onComplete = fn }
}

func WithLogger(logger zerolog.Logger) Opt {
	return func(cfg *config) { cfg.logger = logger }
}

// Work emits a new task in the background
func Work(ctx context.Context, f func(context.Context) error, opts ...Opt) {
	cfg := config{
		name:       "",
		logger:     zerolog.Nop(),
		onComplete: nil,
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
			return
		}

		if cfg.onComplete != nil {
			cfg.onComplete()
		}

		cfg.logger.Debug().Str("worker.name", cfg.getName()).Msg("Background task completed")
	}()
}

func (c config) getName() string {
	if c.name != "" {
		return c.name
	}

	return "unknown"
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

	evt.Str("worker.name", cfg.getName()).Msg("Background task failed")
}
