package testlog

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Run("PanicHandling", func(t *testing.T) {
		// ARRANGE
		tl := New(t)

		// ACT
		// Log indefinitely, even after parent test is done
		go func() {
			for {
				tl.Info().Msg("hello from goroutine")
				time.Sleep(10 * time.Millisecond)
			}
		}()

		time.Sleep(500 * time.Millisecond)
	})

	// Let parent test run a bit longer than PanicHandling
	time.Sleep(time.Second)
}
