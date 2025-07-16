package bg

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

func TestWork(t *testing.T) {
	ctx := context.Background()

	t.Run("basic case", func(t *testing.T) {
		// ARRANGE
		signal := make(chan struct{})

		// ACT
		Work(ctx, func(ctx context.Context) error {
			// simulate some work
			time.Sleep(100 * time.Millisecond)
			close(signal)
			return nil
		})

		// ASSERT
		<-signal
		assertChanClosed(t, signal)
	})

	t.Run("with name and logger", func(t *testing.T) {
		// ARRANGE
		// Given a logger
		out := &bytes.Buffer{}
		logger := zerolog.New(out)

		// And a call returning an error
		call := func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return fmt.Errorf("oopsie")
		}

		// ACT
		Work(ctx, call, WithName("hello"), WithLogger(logger))
		time.Sleep(200 * time.Millisecond)

		// Check the log output
		const expected = `{"level":"error","error":"oopsie","worker.name":"hello","message":"Background task failed"}`
		assert.JSONEq(t, expected, out.String())
	})

	t.Run("with name and logger and onComplete", func(t *testing.T) {
		// ARRANGE
		// Given a logger
		out := &bytes.Buffer{}
		logger := zerolog.New(out)
		check := int64(0)

		// And a call returning an error
		call := func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}

		complete := func() {
			atomic.AddInt64(&check, 1)
		}

		// ACT
		Work(ctx, call, WithName("hello"), WithLogger(logger), OnComplete(complete))
		time.Sleep(200 * time.Millisecond)

		// Check the log output
		const expected = `{"level":"debug", "message":"Background task completed", "worker.name":"hello"}`
		assert.JSONEq(t, expected, out.String())

		// Check onComplete
		assert.Equal(t, int64(1), check)
	})

	t.Run("panic recovery", func(t *testing.T) {
		// ARRANGE
		// Given a logger
		out := &bytes.Buffer{}
		logger := zerolog.New(out)

		// And a call that has panic
		call := func(ctx context.Context) error {
			panic("press F")
			return nil
		}

		// ACT
		Work(ctx, call, WithLogger(logger))
		time.Sleep(100 * time.Millisecond)

		// Check the log output
		const expectedError = "recovered from PANIC in background task: press F"
		const expectedWorker = "unknown"
		const expectedMessage = "Background task failed"
		require.Contains(t, out.String(), expectedError)
		require.Contains(t, out.String(), expectedWorker)
		require.Contains(t, out.String(), expectedMessage)
	})
}

func assertChanClosed(t *testing.T, ch <-chan struct{}) {
	_, ok := <-ch
	assert.False(t, ok, "channel is not closed")
}
