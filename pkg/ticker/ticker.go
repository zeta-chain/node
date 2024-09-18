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
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"cosmossdk.io/errors"
)

// Ticker represents a ticker that will run a function periodically.
// It also invokes BEFORE ticker starts.
type Ticker struct {
	interval   time.Duration
	ticker     *time.Ticker
	task       Task
	signalChan chan struct{}

	// runnerMu is a mutex to prevent double run
	runnerMu sync.Mutex

	// stateMu is a mutex to prevent concurrent SetInterval calls
	stateMu sync.Mutex

	stopped bool
}

// Task is a function that will be called by the Ticker
type Task func(ctx context.Context, t *Ticker) error

// New creates a new Ticker.
func New(interval time.Duration, runner Task) *Ticker {
	return &Ticker{interval: interval, task: runner}
}

// Run creates and runs a new Ticker.
func Run(ctx context.Context, interval time.Duration, task Task) error {
	return New(interval, task).Run(ctx)
}

// SecondsFromUint64 converts uint64 to time.Duration in seconds.
func SecondsFromUint64(d uint64) time.Duration {
	return time.Duration(d) * time.Second
}

// Run runs the ticker by blocking current goroutine. It also invokes BEFORE ticker starts.
// Stops when (if any):
// - context is done (returns ctx.Err())
// - task returns an error or panics
// - shutdown signal is received
func (t *Ticker) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			lines := strings.Split(stack, "\n")
			line := ""
			// 8th line should be the actual line, see the unit tests
			if len(lines) > 8 {
				line = strings.TrimSpace(lines[8])
			}
			err = fmt.Errorf("panic during ticker run: %v at %s", r, line)
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
	if err := t.task(ctx, t); err != nil {
		return errors.Wrap(err, "ticker task failed")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.ticker.C:
			if err := t.task(ctx, t); err != nil {
				return errors.Wrap(err, "ticker task failed")
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
	t.ticker.Stop()
}
