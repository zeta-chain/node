package ratelimit

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimiter(t *testing.T) {
	// Given rate limiter
	r := New(3)

	// Acquire 3 requests
	require.NoError(t, r.Acquire(1, 100))
	require.NoError(t, r.Acquire(2, 200))
	require.NoError(t, r.Acquire(3, 300))

	require.Equal(t, uint64(3), r.Pending())

	// Should be throttled
	require.ErrorIs(t, r.Acquire(4, 400), ErrThrottled)

	// Release 3 requests
	r.Release()
	r.Release()
	r.Release()

	// Should be allowed
	require.NoError(t, r.Acquire(4, 401))
	r.Release()

	// noop
	r.Release()

	require.Equal(t, uint64(0), r.Pending())
}

// TestRateLimiterRobustness tests that the rate limiter handles edge cases robustly
func TestRateLimiterRobustness(t *testing.T) {
	// This test verifies that the rate limiter handles edge cases correctly
	// The implementation ensures proper ordering of operations and handles
	// potential edge cases gracefully.

	t.Run("Handles Edge Cases Correctly", func(t *testing.T) {
		r := New(5)

		// Acquire 1 permit
		require.NoError(t, r.Acquire(1, 100))
		require.Equal(t, uint64(1), r.Pending())

		// Release 1 time (normal)
		r.Release()
		require.Equal(t, uint64(0), r.Pending())

		// Release 5 more times (excessive) - this should be safe now
		for i := 0; i < 5; i++ {
			r.Release() // Fixed code prevents over-release
		}

		// Should still be able to acquire normally
		require.NoError(t, r.Acquire(2, 200))
		require.Equal(t, uint64(1), r.Pending())
	})
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	r := New(10)
	var wg sync.WaitGroup
	concurrentCount := 50

	// Test concurrent acquires and releases
	for i := 0; i < concurrentCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Try to acquire
			err := r.Acquire(uint64(id), uint64(id))
			if err == nil {
				// Simulate some work
				time.Sleep(1 * time.Millisecond)
				r.Release()
			}
		}(i)
	}

	wg.Wait()

	// After all operations, pending should be 0
	require.Equal(t, uint64(0), r.Pending())
}

func TestRateLimiterExcessiveReleases(t *testing.T) {
	r := New(5)

	// Acquire 3 permits
	require.NoError(t, r.Acquire(1, 100))
	require.NoError(t, r.Acquire(2, 200))
	require.NoError(t, r.Acquire(3, 300))

	require.Equal(t, uint64(3), r.Pending())

	// Release 3 times (normal)
	r.Release()
	r.Release()
	r.Release()

	require.Equal(t, uint64(0), r.Pending())

	// Release 5 more times (excessive)
	for i := 0; i < 5; i++ {
		r.Release() // Should not cause panic or negative counter
	}

	// Pending should still be 0
	require.Equal(t, uint64(0), r.Pending())

	// Should still be able to acquire
	require.NoError(t, r.Acquire(4, 400))
	require.Equal(t, uint64(1), r.Pending())
}

func TestRateLimiterRaceCondition(t *testing.T) {
	r := New(20)
	var wg sync.WaitGroup
	releaseCount := 100

	// Acquire some permits first
	for i := 0; i < 10; i++ {
		require.NoError(t, r.Acquire(uint64(i), uint64(i)))
	}

	require.Equal(t, uint64(10), r.Pending())

	// Concurrently release all permits
	for i := 0; i < releaseCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Release()
		}()
	}

	wg.Wait()

	// Pending should be 0 (not negative)
	require.Equal(t, uint64(0), r.Pending())

	// Should still be able to acquire new permits
	require.NoError(t, r.Acquire(100, 100))
	require.Equal(t, uint64(1), r.Pending())
}

func TestRateLimiterUnderflowProtection(t *testing.T) {
	r := New(5)

	// Don't acquire any permits
	require.Equal(t, uint64(0), r.Pending())

	// Try to release multiple times without acquiring
	for i := 0; i < 10; i++ {
		r.Release() // Should not cause underflow or panic
	}

	// Pending should still be 0
	require.Equal(t, uint64(0), r.Pending())

	// Should still be able to acquire after excessive releases
	require.NoError(t, r.Acquire(1, 100))
	require.Equal(t, uint64(1), r.Pending())

	// Release normally
	r.Release()
	require.Equal(t, uint64(0), r.Pending())
}

// TestPendingFunctionUnderflowProtection tests that Pending() never returns negative values
func TestPendingFunctionUnderflowProtection(t *testing.T) {
	r := New(5)

	// Test that Pending() returns 0 when no permits are acquired
	require.Equal(t, uint64(0), r.Pending())

	// Acquire and release to potentially cause underflow
	require.NoError(t, r.Acquire(1, 100))
	r.Release()

	// Release multiple times to potentially cause negative internal state
	for i := 0; i < 10; i++ {
		r.Release()
	}

	// Pending() should never return a negative value, even if internal state is negative
	require.Equal(t, uint64(0), r.Pending())

	// Should still be able to acquire normally
	require.NoError(t, r.Acquire(2, 200))
	require.Equal(t, uint64(1), r.Pending())
}
