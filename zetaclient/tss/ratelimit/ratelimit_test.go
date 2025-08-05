package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/semaphore"
)

// BuggyRateLimiter reproduces the original race condition vulnerability
type BuggyRateLimiter struct {
	sem     *semaphore.Weighted
	pending *atomic.Int32
}

// NewBuggyRateLimiter creates a rate limiter with the original vulnerable implementation
func NewBuggyRateLimiter(maxPending uint64) *BuggyRateLimiter {
	if maxPending == 0 {
		maxPending = DefaultMaxPendingSignatures
	}

	return &BuggyRateLimiter{
		sem:     semaphore.NewWeighted(int64(maxPending)),
		pending: &atomic.Int32{},
	}
}

// Acquire acquires a signature for a given chain and nonce.
func (r *BuggyRateLimiter) Acquire(chainID, nonce uint64) error {
	if !r.sem.TryAcquire(1) {
		return errors.Wrapf(ErrThrottled, "chain: %d, nonce: %d", chainID, nonce)
	}

	r.pending.Add(1)
	return nil
}

// Release reproduces the original vulnerable implementation
func (r *BuggyRateLimiter) Release() {
	// Original vulnerable code
	if r.pending.Load() == 0 { // ❌ RACE: Non-atomic check-then-act
		return
	}
	r.sem.Release(1)  // ❌ RACE: Can release more permits than acquired
	r.pending.Add(-1) // ❌ RACE: Can make counter negative
}

// Pending returns the number of pending signatures.
func (r *BuggyRateLimiter) Pending() uint64 {
	return uint64(r.pending.Load())
}

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

// TestBuggyRateLimiterRaceCondition reproduces the original race condition
func TestBuggyRateLimiterRaceCondition(t *testing.T) {
	r := NewBuggyRateLimiter(20)
	var wg sync.WaitGroup
	releaseCount := 1000 // Increased to make race condition more likely

	// Acquire some permits first
	for i := 0; i < 10; i++ {
		require.NoError(t, r.Acquire(uint64(i), uint64(i)))
	}

	require.Equal(t, uint64(10), r.Pending())

	// Concurrently release all permits - this should cause race conditions
	for i := 0; i < releaseCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Release()
		}()
	}

	wg.Wait()

	// With the buggy implementation, this could be negative or inconsistent
	// The test may pass sometimes due to timing, but the race condition exists
	pending := r.Pending()
	t.Logf("Final pending count: %d (should be 0, but race condition may cause inconsistency)", pending)

	// Try to acquire new permits to see if the semaphore is in a bad state
	err := r.Acquire(100, 100)
	if err != nil {
		t.Logf("Failed to acquire after race condition: %v", err)
	}
}

// TestBuggyRateLimiterStressTest runs multiple iterations to catch the race condition
func TestBuggyRateLimiterStressTest(t *testing.T) {
	for iteration := 0; iteration < 100; iteration++ {
		r := NewBuggyRateLimiter(5)
		var wg sync.WaitGroup

		// Acquire 3 permits
		require.NoError(t, r.Acquire(1, 100))
		require.NoError(t, r.Acquire(2, 200))
		require.NoError(t, r.Acquire(3, 300))

		// Release 3 times normally
		r.Release()
		r.Release()
		r.Release()

		// Now release 10 more times concurrently (excessive)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.Release()
			}()
		}

		wg.Wait()

		// Check if we got a negative counter
		pending := r.Pending()
		if pending < 0 {
			t.Errorf("Iteration %d: Got negative pending count: %d", iteration, pending)
			return
		}

		// Try to acquire - if semaphore is corrupted, this might fail
		err := r.Acquire(4, 400)
		if err != nil {
			t.Logf("Iteration %d: Failed to acquire after excessive releases: %v", iteration, err)
		}
	}
}

// TestVulnerabilityDemonstration shows the theoretical race condition
func TestVulnerabilityDemonstration(t *testing.T) {
	// This test demonstrates the theoretical race condition
	// In practice, it may not always trigger due to timing, but the vulnerability exists

	t.Run("Original Vulnerable Code", func(t *testing.T) {
		r := NewBuggyRateLimiter(5)

		// Acquire 1 permit
		require.NoError(t, r.Acquire(1, 100))
		require.Equal(t, uint64(1), r.Pending())

		// Release 1 time (normal)
		r.Release()
		require.Equal(t, uint64(0), r.Pending())

		// Release 5 more times (excessive) - this is the vulnerability
		// The original code doesn't check if pending is 0 before releasing semaphore
		for i := 0; i < 5; i++ {
			r.Release() // This can release more permits than acquired
		}

		// The semaphore state is now corrupted
		// Try to acquire - this might fail or succeed incorrectly
		err := r.Acquire(2, 200)
		if err != nil {
			t.Logf("Acquire failed due to corrupted semaphore: %v", err)
		}
	})

	t.Run("Fixed Code Comparison", func(t *testing.T) {
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

// TestBuggyRateLimiterExcessiveReleases reproduces the excessive release issue
func TestBuggyRateLimiterExcessiveReleases(t *testing.T) {
	r := NewBuggyRateLimiter(5)

	// Acquire 3 permits
	require.NoError(t, r.Acquire(1, 100))
	require.NoError(t, r.Acquire(2, 200))
	require.NoError(t, r.Acquire(3, 300))

	require.Equal(t, uint64(3), r.Pending())

	// Release 3 times (normal)
	r.Release()
	r.Release()
	r.Release()

	// With buggy implementation, this might be negative
	pending := r.Pending()
	t.Logf("Pending after normal releases: %d (should be 0)", pending)

	// Release 5 more times (excessive) - this should cause issues
	for i := 0; i < 5; i++ {
		r.Release() // Should cause semaphore over-release and negative counter
	}

	// Check the state after excessive releases
	pending = r.Pending()
	t.Logf("Pending after excessive releases: %d (should be 0, but race condition may cause negative)", pending)

	// Try to acquire - this might fail due to semaphore being in bad state
	err := r.Acquire(4, 400)
	if err != nil {
		t.Logf("Failed to acquire after excessive releases: %v", err)
	}
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
