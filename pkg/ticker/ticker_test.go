package ticker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicker(t *testing.T) {
	const (
		dur      = time.Millisecond * 100
		durSmall = dur / 10
	)

	t.Run("Basic case with context", func(t *testing.T) {
		// ARRANGE
		// Given a counter
		var counter int

		// And a context
		ctx, cancel := context.WithTimeout(context.Background(), dur+durSmall)
		defer cancel()

		// And a ticker
		ticker := New(dur, func(_ context.Context, t *Ticker) error {
			counter++

			return nil
		})

		// ACT
		err := ticker.Run(ctx)

		// ASSERT
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		// two runs: start run + 1 tick
		assert.Equal(t, 2, counter)
	})

	t.Run("Halts when error occurred", func(t *testing.T) {
		// ARRANGE
		// Given a counter
		var counter int

		ctx := context.Background()

		// And a ticker func that returns an error after 10 runs
		ticker := New(durSmall, func(_ context.Context, t *Ticker) error {
			counter++
			if counter > 9 {
				return fmt.Errorf("oops")
			}

			return nil
		})

		// ACT
		err := ticker.Run(ctx)

		// ASSERT
		assert.ErrorContains(t, err, "oops")
		assert.Equal(t, 10, counter)
	})

	t.Run("Dynamic interval update", func(t *testing.T) {
		// ARRANGE
		// Given a counter
		var counter int

		// Given duration
		duration := dur * 10

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		// And a ticker what decreases the interval by 2 each time
		ticker := New(durSmall, func(_ context.Context, ticker *Ticker) error {
			t.Logf("Counter: %d, Duration: %s", counter, duration.String())

			counter++
			duration /= 2

			ticker.SetInterval(duration)

			return nil
		})

		// ACT
		err := ticker.Run(ctx)

		// ASSERT
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		// It should have run at 2 times with ctxTimeout = tickerDuration (start + 1 tick),
		// But it should have run more than that because of the interval decrease
		assert.GreaterOrEqual(t, counter, 2)
	})

	t.Run("Stop ticker", func(t *testing.T) {
		// ARRANGE
		// Given a counter
		var counter int

		// And a context
		ctx := context.Background()

		// And a ticker
		ticker := New(durSmall, func(_ context.Context, _ *Ticker) error {
			counter++
			return nil
		})

		// And a function with a stop signal
		go func() {
			time.Sleep(dur)
			ticker.Stop()
		}()

		// ACT
		err := ticker.Run(ctx)

		// ASSERT
		assert.NoError(t, err)
		assert.Greater(t, counter, 8)

		t.Run("Stop ticker for the second time", func(t *testing.T) {
			ticker.Stop()
		})
	})

	t.Run("Panic", func(t *testing.T) {
		// ARRANGE
		// Given a context
		ctx := context.Background()

		// And a ticker
		ticker := New(durSmall, func(_ context.Context, _ *Ticker) error {
			panic("oops")
		})

		// ACT
		err := ticker.Run(ctx)

		// ASSERT
		assert.ErrorContains(t, err, "panic during ticker run: oops")
		// assert that we get error with the correct line number
		assert.ErrorContains(t, err, "ticker_test.go:142")
	})

	t.Run("Nil panic", func(t *testing.T) {
		// ARRANGE
		// Given a context
		ctx := context.Background()

		// And a ticker
		ticker := New(durSmall, func(_ context.Context, _ *Ticker) error {
			var a func()
			a()
			return nil
		})

		// ACT
		err := ticker.Run(ctx)

		// ASSERT
		assert.ErrorContains(
			t,
			err,
			"panic during ticker run: runtime error: invalid memory address or nil pointer dereference",
		)
		// assert that we get error with the correct line number
		assert.ErrorContains(t, err, "ticker_test.go:162")
	})

	t.Run("Run as a single call", func(t *testing.T) {
		// ARRANGE
		// Given a counter
		var counter int

		// Given a context
		ctx, cancel := context.WithTimeout(context.Background(), dur+durSmall)
		defer cancel()

		tick := func(ctx context.Context, t *Ticker) error {
			counter++
			return nil
		}

		// ACT
		err := Run(ctx, dur, tick)

		// ASSERT
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, 2, counter)
	})
}
