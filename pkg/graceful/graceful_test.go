package graceful

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

func TestProcess(t *testing.T) {
	const defaultTimeout = 2 * time.Second

	ctx := context.Background()

	t.Run("Service sync", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)

		// ACT
		// Run service
		ts.process.AddService(ctx, ts.mockService)

		start := time.Now()

		// And after 1 second someone presses ctrl+c
		go func() {
			time.Sleep(time.Second)
			ts.mockSignal <- os.Interrupt
		}()

		ts.process.WaitForShutdown()

		// ASSERT
		// Check that service was stopped in a timely manner
		assert.Less(t, time.Since(start), defaultTimeout)
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
		assert.Contains(t, ts.logger.String(), "mock is running in blocking mode")
	})

	t.Run("Service async", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, true)

		// Run service
		ts.process.AddService(ctx, ts.mockService)

		// ACT
		start := time.Now()

		// And after 700ms someone presses ctrl+c
		go func() {
			time.Sleep(700 * time.Millisecond)
			ts.mockSignal <- os.Interrupt
		}()

		ts.process.WaitForShutdown()

		// ASSERT
		// Check that service was stopped in a timely manner
		assert.Less(t, time.Since(start), defaultTimeout)
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
		assert.Contains(t, ts.logger.String(), "mock is running in non-blocking mode")
	})

	t.Run("Manual starters and stoppers", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)

		// Given one starter
		ts.process.AddStarter(ctx, func(ctx context.Context) error {
			ts.logger.Info().Msg("Hello world")
			return nil
		})

		// And two stoppers
		ts.process.AddStopper(func() {
			time.Sleep(200 * time.Millisecond)
			ts.logger.Info().Msg("Stopper 1")
		})

		ts.process.AddStopper(func() {
			time.Sleep(300 * time.Millisecond)
			ts.logger.Info().Msg("Stopper 2")
		})

		// ACT
		start := time.Now()

		// And after 1s someone presses ctrl+c
		go func() {
			time.Sleep(time.Second)
			ts.mockSignal <- os.Interrupt
		}()

		ts.process.WaitForShutdown()

		// ASSERT
		// Check that service was stopped in a timely manner
		assert.Less(t, time.Since(start), defaultTimeout)
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
		assert.Contains(t, ts.logger.String(), "Stopper 1")
		assert.Contains(t, ts.logger.String(), "Stopper 2")
	})

	t.Run("Starter error", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)

		ts.mockService.errStart = fmt.Errorf("failed to start service")

		ts.process.AddService(ctx, ts.mockService)

		// ACT
		start := time.Now()
		ts.process.WaitForShutdown()

		// ASSERT
		// Check that service had errors and was stopped
		assert.Less(t, time.Since(start), defaultTimeout)
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
		assert.Contains(t, ts.logger.String(), "failed to start service")
	})

	t.Run("Panic handling during startup", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)

		ts.process.AddStarter(ctx, func(ctx context.Context) error {
			panic("oopsie")
			return nil
		})

		// ACT
		ts.process.WaitForShutdown()

		// ASSERT
		// Check that service had errors and was stopped
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
		assert.Contains(t, ts.logger.String(), "panic in service")

		// Check that error contains exact line of panic
		assert.Contains(t, ts.logger.String(), "graceful_test.go:144")
	})

	t.Run("Panic handling during shutdown", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)

		ts.process.AddStopper(func() {
			panic("bombarda maxima")
		})

		// ACT
		ts.process.ShutdownNow()

		// ASSERT
		// Check that service had errors and was stopped
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
		assert.Contains(t, ts.logger.String(), "panic during shutdown")

		// Check that error contains exact line of panic
		assert.Contains(t, ts.logger.String(), "graceful_test.go:167")
	})

	t.Run("WaitForShutdown noop", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)
		ts.process.AddService(ctx, ts.mockService)

		// ACT
		ts.process.ShutdownNow()
		ts.process.WaitForShutdown()

		// ASSERT
		assert.Contains(t, ts.logger.String(), "Shutdown completed")
	})

	t.Run("Shutdown timeout", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		ts := newTestSuite(t, defaultTimeout, false)

		// Given some slow stopper
		const workDuration = defaultTimeout + 5*time.Second

		ts.process.AddStopper(func() {
			ts.logger.Info().Msg("Stopping something")
			time.Sleep(workDuration)
			ts.logger.Info().Msg("Stopped something")
		})

		// ACT
		ts.process.ShutdownNow()

		// ASSERT
		assert.Contains(t, ts.logger.String(), "Stopping something")
		assert.Contains(t, ts.logger.String(), "Shutdown interrupted by timeout")

		// log doesn't contain this line because it was interrupted
		assert.NotContains(t, ts.logger.String(), "Stopped something")
	})
}

type testSuite struct {
	process     *Process
	mockService *mockService
	mockSignal  chan os.Signal

	logger *testlog.Log
}

func newTestSuite(t *testing.T, timeout time.Duration, async bool) *testSuite {
	logger := testlog.New(t)

	stop := NewSigChan(os.Interrupt)
	process := New(timeout, logger.Logger, stop)

	return &testSuite{
		mockSignal: stop,
		process:    process,
		logger:     logger,
		mockService: &mockService{
			async:  async,
			Logger: logger.Logger,
		},
	}
}

type mockService struct {
	errStart error
	async    bool
	running  bool
	zerolog.Logger
}

func (m *mockService) Start(_ context.Context) error {
	const interval = 300 * time.Millisecond

	m.running = true

	// emulate async started
	if m.async {
		go func() {
			for {
				if m.errStart != nil || !m.running {
					return
				}

				m.Info().Msg("mock is running in non-blocking mode")
				time.Sleep(interval)
			}
		}()

		return nil
	}

	for {
		switch {
		case m.errStart != nil:
			m.running = false
			return m.errStart
		case !m.running:
			return nil
		default:
			m.Info().Msg("mock is running in blocking mode")
			time.Sleep(interval)
		}
	}
}

func (m *mockService) Stop() {
	m.running = false
	m.Info().Msg("Stopping mock service")
}
