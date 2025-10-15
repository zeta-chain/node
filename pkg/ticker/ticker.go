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

	"github.com/rs/zerolog"
)

// Task is a function that will be called by the Ticker
type Task func(ctx context.Context, t *Ticker) error

// Ticker represents a ticker that will run a function periodically.
// It also invokes BEFORE ticker starts.
type Ticker struct {
	interval time.Duration
	ticker   *time.Ticker
	task     Task

	// runnerMu is a mutex to prevent double run
	runnerMu sync.Mutex

	// stateMu is a mutex to prevent concurrent SetInterval calls
	stateMu sync.Mutex

	stopped   bool
	ctxCancel context.CancelFunc

	runCompleteChan  chan struct{}
	externalStopChan <-chan struct{}
	logger           zerolog.Logger
}

// Opt is a configuration option for the Ticker.
type Opt func(*Ticker)

// WithLogger sets the logger for the Ticker.
func WithLogger(log zerolog.Logger, name string) Opt {
	return func(t *Ticker) {
		t.logger = log.With().Str("ticker_name", name).Logger()
	}
}

// WithStopChan sets the stop channel for the Ticker.
// Please note that stopChan is NOT signalChan.
// Stop channel is a trigger for invoking ticker.Stop();
func WithStopChan(stopChan <-chan struct{}) Opt {
	return func(cfg *Ticker) { cfg.externalStopChan = stopChan }
}

// New creates a new Ticker.
func New(interval time.Duration, task Task, opts ...Opt) *Ticker {
	t := &Ticker{
		interval: interval,
		task:     task,
		logger:   zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Run creates and runs a new Ticker.
func Run(ctx context.Context, interval time.Duration, task Task, opts ...Opt) error {
	return New(interval, task, opts...).Start(ctx)
}

// Start runs the ticker by blocking current goroutine. It also invokes BEFORE ticker starts.
// Stops when (if any):
// - context is done (returns ctx.Err())
// - task returns an error or panics
// - shutdown signal is received
func (t *Ticker) Start(ctx context.Context) (err error) {
	// prevent concurrent runs
	t.runnerMu.Lock()
	defer t.runnerMu.Unlock()

	ctx = t.setStartState(ctx)

	defer func() {
		// used in StopBlocking()
		close(t.runCompleteChan)

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

	// initial run
	if err = t.task(ctx, t); err != nil {
		t.Stop()
		return fmt.Errorf("ticker task failed (initial run): %w", err)
	}

	defer t.setStopState()

	for {
		select {
		case <-ctx.Done():
			// if task is finished (i.e. last tick completed BEFORE ticker.Stop(),
			// then we need to return nil)
			if t.isStopped() {
				return nil
			}
			return ctx.Err()
		case <-t.ticker.C:
			// If another goroutine calls ticker.Stop() while the current tick is running,
			// Then it's okay to return ctx error
			if err := t.task(ctx, t); err != nil {
				return fmt.Errorf("ticker task failed: %w", err)
			}
		case <-t.externalStopChan:
			t.Stop()
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

	t.logger.Info().
		Float64("ticker.old_interval", t.interval.Seconds()).
		Float64("ticker.new_interval", interval.Seconds()).
		Msg("Changing interval")

	t.interval = interval
	t.ticker.Reset(interval)
}

// Stop stops the ticker in a NON-blocking way. If the task is running in a separate goroutine,
// this call *might* not wait for it to finish. To wait for task finish, use StopBlocking().
// It's safe to call Stop() multiple times / concurrently / within the task.
func (t *Ticker) Stop() {
	t.setStopState()
}

// StopBlocking stops the ticker in a blocking way i.e. it waits for the task to finish.
// DO NOT call this within the task.
func (t *Ticker) StopBlocking() {
	t.setStopState()
	<-t.runCompleteChan
}

func (t *Ticker) setStartState(ctx context.Context) context.Context {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	ctx, t.ctxCancel = context.WithCancel(ctx)
	t.ticker = time.NewTicker(t.interval)
	t.stopped = false

	// this signals that Run() is about to return
	t.runCompleteChan = make(chan struct{})

	return ctx
}

func (t *Ticker) setStopState() {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	// noop
	if t.stopped {
		return
	}

	t.ctxCancel()
	t.stopped = true
	if t.ticker != nil {
		t.ticker.Stop()
	}

	t.logger.Info().Msgf("Ticker stopped")
}

func (t *Ticker) isStopped() bool {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	return t.stopped
}

// DurationFromUint64Seconds converts uint64 of seconds to time.Duration.
func DurationFromUint64Seconds(seconds uint64) time.Duration {
	// #nosec G115 seconds should be in range and is not user controlled
	return time.Duration(seconds) * time.Second
}
