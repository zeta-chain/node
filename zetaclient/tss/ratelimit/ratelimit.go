package ratelimit

import (
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

// RateLimiter is a simple semaphore for limiting the number of concurrent signatures.
//
// This is a naive implementation that probably requires
// improvements to handle real-world scenarios.
//
// Pros:
// - Has a simple interface that hides the underlying implementation details;
//
// Cons:
// - Doesn't take number of signatures per chain into account;
// - Doesn't take nonce ordering into account;
// - Doesn't take chain "fairness" into account;
//
// TODO TBD:
// https://github.com/zeta-chain/node/issues/3830
//
//  1. We could get/adjust the value from an on-chain param instead of the config.
//  2. How to ensure that each O+S throttles the same CCTX at a given point in time?
//     Otherwise, different nodes might throttle different cctx => no party formed => error
type RateLimiter struct {
	sem     *semaphore.Weighted
	pending *atomic.Int32
}

var ErrThrottled = errors.New("action is throttled")

// New RateLimiter constructor.
func New(maxPending int64) *RateLimiter {
	return &RateLimiter{
		sem:     semaphore.NewWeighted(maxPending),
		pending: &atomic.Int32{},
	}
}

// Acquire acquires a signature for a given chain and nonce.
// Returns ErrThrottled if the rate limit is exceeded.
func (r *RateLimiter) Acquire(chainID, nonce uint64) error {
	if !r.sem.TryAcquire(1) {
		return errors.Wrapf(ErrThrottled, "chain: %d, nonce: %d", chainID, nonce)
	}

	r.pending.Add(1)

	return nil
}

// Release releases a signature for a given chain and nonce.
func (r *RateLimiter) Release(_, _ uint64) {
	// noop
	if r.pending.Load() == 0 {
		return
	}

	r.sem.Release(1)
	r.pending.Add(-1)
}

// Pending returns the number of pending signatures.
func (r *RateLimiter) Pending() uint64 {
	// #nosec G115 always in range
	return uint64(r.pending.Load())
}
