package ticker

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicker(t *testing.T) {
	const (
		dur      = time.Millisecond * 100
		durSmall = dur / 10
	)

	t.Run("Basic case with context", func(t *testing.T) {
		t.Parallel()

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
		err := ticker.Start(ctx)

		// ASSERT
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		// two runs: start run + 1 tick
		assert.Equal(t, 2, counter)
	})

	t.Run("Halts when error occurred", func(t *testing.T) {
		t.Parallel()

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
		err := ticker.Start(ctx)

		// ASSERT
		assert.ErrorContains(t, err, "oops")
		assert.Equal(t, 10, counter)
	})

	t.Run("Dynamic interval update", func(t *testing.T) {
		t.Parallel()

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
		err := ticker.Start(ctx)

		// ASSERT
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		// It should have run at 2 times with ctxTimeout = tickerDuration (start + 1 tick),
		// But it should have run more than that because of the interval decrease
		assert.GreaterOrEqual(t, counter, 2)
	})

	t.Run("Stop ticker", func(t *testing.T) {
		t.Parallel()

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
		err := ticker.Start(ctx)

		// ASSERT
		assert.NoError(t, err)
		assert.Greater(t, counter, 8)

		t.Run("Stop ticker for the second time", func(t *testing.T) {
			ticker.Stop()
		})
	})

	t.Run("Stop ticker in a blocking fashion", func(t *testing.T) {
		t.Parallel()

		const (
			tickerInterval = 100 * time.Millisecond
			workDuration   = 600 * time.Millisecond
			stopAfterStart = workDuration + tickerInterval/2
		)

		newLogger := func(t *testing.T) zerolog.Logger {
			return zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()
		}

		// test task that imitates some work
		newTask := func(counter *int32, logger zerolog.Logger) Task {
			return func(ctx context.Context, _ *Ticker) error {
				logger.Info().Msg("Tick start")
				atomic.AddInt32(counter, 1)

				time.Sleep(workDuration)

				logger.Info().Msgf("Tick end")
				atomic.AddInt32(counter, -1)

				return nil
			}
		}

		t.Run("Non-blocking stop fails do finish the work", func(t *testing.T) {
			t.Parallel()

			// ARRANGE
			// Given some test task that imitates some work
			testLogger := newLogger(t)
			counter := int32(0)
			task := newTask(&counter, testLogger)

			// Given a ticker
			ticker := New(tickerInterval, task, WithLogger(testLogger, "test-non-blocking-ticker"))

			// ACT
			// Imitate the ticker run in the background
			go func() {
				err := ticker.Start(context.Background())
				require.NoError(t, err)
			}()

			// Then stop the ticker after some delay
			time.Sleep(stopAfterStart)
			testLogger.Info().Msg("Stopping ticker")
			ticker.Stop()
			testLogger.Info().Msg("Stopped ticker")

			// ASSERT
			// If ticker is stopped BEFORE the work is done i.e. "in the middle of work",
			// thus the counter would be `1. You can also check the logs
			assert.Equal(t, int32(1), counter)
		})

		t.Run("Blocking stop works as expected", func(t *testing.T) {
			t.Parallel()

			// ARRANGE
			// Now if we have the SAME test but with blocking stop, it should work
			testLogger := newLogger(t)
			counter := int32(0)
			task := newTask(&counter, testLogger)

			ticker := New(tickerInterval, task, WithLogger(testLogger, "test-non-blocking-ticker"))

			// ACT
			go func() {
				err := ticker.Start(context.Background())
				require.NoError(t, err)
			}()

			time.Sleep(stopAfterStart)
			testLogger.Info().Msg("Stopping ticker")
			ticker.StopBlocking()
			testLogger.Info().Msg("Stopped ticker")

			// ASSERT
			// If ticker is stopped AFTER the work is done
			assert.Equal(t, int32(0), counter)
		})
	})

	t.Run("Panic", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		// Given a context
		ctx := context.Background()

		// And a ticker
		ticker := New(durSmall, func(_ context.Context, _ *Ticker) error {
			panic("oops")
		})

		// ACT
		err := ticker.Start(ctx)

		// ASSERT
		assert.ErrorContains(t, err, "panic during ticker run: oops")
		// assert that we get error with the correct line number
		assert.ErrorContains(t, err, "ticker_test.go:243")
	})

	t.Run("Nil panic", func(t *testing.T) {
		t.Parallel()

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
		err := ticker.Start(ctx)

		// ASSERT
		assert.ErrorContains(
			t,
			err,
			"panic during ticker run: runtime error: invalid memory address or nil pointer dereference",
		)
		// assert that we get error with the correct line number
		assert.ErrorContains(t, err, "ticker_test.go:265")
	})

	t.Run("Run as a single call", func(t *testing.T) {
		t.Parallel()

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

	t.Run("With stop channel", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		var (
			tickerInterval = 100 * time.Millisecond
			counter        = 0

			stopChan        = make(chan struct{})
			sleepBeforeStop = 5*tickerInterval + (10 * time.Millisecond)
		)

		task := func(ctx context.Context, _ *Ticker) error {
			t.Logf("Tick %d", counter)
			counter++

			return nil
		}

		// ACT
		go func() {
			time.Sleep(sleepBeforeStop)
			close(stopChan)
		}()

		err := Run(context.Background(), tickerInterval, task, WithStopChan(stopChan))

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, 6, counter) // initial tick + 5 more ticks
	})

	t.Run("With logger", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		out := &bytes.Buffer{}
		logger := zerolog.New(out)

		// ACT
		task := func(ctx context.Context, _ *Ticker) error {
			return fmt.Errorf("hey")
		}

		err := Run(context.Background(), time.Second, task, WithLogger(logger, "my-task"))

		// ARRANGE
		require.ErrorContains(t, err, "hey")
		require.Contains(t, out.String(), `{"level":"info","ticker_name":"my-task","message":"Ticker stopped"}`)
	})
}
