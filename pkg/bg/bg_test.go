package bg

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
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
		const expected = `{"level":"error","error":"recovered from PANIC in background task: press F",` +
			`"worker.name":"unknown","message":"Background task failed"}`
		assert.JSONEq(t, expected, out.String())
	})
}

func assertChanClosed(t *testing.T, ch <-chan struct{}) {
	_, ok := <-ch
	assert.False(t, ok, "channel is not closed")
}
