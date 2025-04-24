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
	r.Release(1, 100)
	r.Release(2, 200)
	r.Release(3, 300)

	// Should be allowed
	require.NoError(t, r.Acquire(4, 401))
	r.Release(4, 401)

	// noop
	r.Release(4, 402)

	require.Equal(t, uint64(0), r.Pending())
}
