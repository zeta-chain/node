// Package ticker provides a dynamic ticker that can change its interval at runtime.
// The ticker can be stopped gracefully and handles context-based termination.
//
// This package is useful for scenarios where periodic execution of a function is needed
// and the interval might need to change dynamically based on runtime conditions.
//
// It also invokes a first tick immediately after the ticker starts. It's safe to use it concurrently.
//
// It also terminates gracefully when the context is done (return ctx.Err()) or when the stop signal is received.
//
// Example usage:
//
//	ticker := New(time.Second, func(ctx context.Context, t *Ticker) error {
//	    resp, err := client.GetPrice(ctx)
//	    if err != nil {
//	        logger.Err(err).Error().Msg("failed to get price")
//	        return nil
//	    }
//
//	    observer.SetPrice(resp.GasPrice)
//	    t.SetInterval(resp.GasPriceInterval)
//
//	    return nil
//	})
//
//	err := ticker.Run(ctx)
package ticker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/errors"
)

// Ticker represents a ticker that will run a function periodically.
// It also invokes BEFORE ticker starts.
type Ticker struct {
	interval   time.Duration
	ticker     *time.Ticker
	runner     Runner
	signalChan chan struct{}

	// runnerMu is a mutex to prevent double run
	runnerMu sync.Mutex

	// stateMu is a mutex to prevent concurrent SetInterval calls
	stateMu sync.Mutex

	stopped bool
}

// Runner is a function that will be called by the Ticker
type Runner func(ctx context.Context, t *Ticker) error

// New creates a new Ticker.
func New(interval time.Duration, runner Runner) *Ticker {
	return &Ticker{interval: interval, runner: runner}
}

// Run creates and runs a new Ticker.
func Run(ctx context.Context, interval time.Duration, runner Runner) error {
	return New(interval, runner).Run(ctx)
}

// Run runs the ticker by blocking current goroutine. It also invokes BEFORE ticker starts.
// Stops when (if any):
// - context is done (returns ctx.Err())
// - runner returns an error or panics
// - shutdown signal is received
func (t *Ticker) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during ticker run: %v", r)
		}
	}()

	// prevent concurrent runs
	t.runnerMu.Lock()
	defer t.runnerMu.Unlock()

	// setup
	t.ticker = time.NewTicker(t.interval)
	t.signalChan = make(chan struct{})
	t.stopped = false

	// initial run
	if err := t.runner(ctx, t); err != nil {
		return errors.Wrap(err, "ticker runner failed")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.ticker.C:
			if err := t.runner(ctx, t); err != nil {
				return errors.Wrap(err, "ticker runner failed")
			}
		case <-t.signalChan:
			return nil
		}
	}
}

// SetInterval updates the interval of the ticker.
func (t *Ticker) SetInterval(interval time.Duration) {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	// noop
	if t.interval == interval || t.ticker == nil {
		return
	}

	t.interval = interval
	t.ticker.Reset(interval)
}

// Stop stops the ticker. Safe to call concurrently or multiple times.
func (t *Ticker) Stop() {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	// noop
	if t.stopped || t.signalChan == nil {
		return
	}

	close(t.signalChan)
	t.stopped = true
}
