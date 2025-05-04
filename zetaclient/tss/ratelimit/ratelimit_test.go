package ratelimit

import (
	"testing"

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
