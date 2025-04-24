package ratelimit

import (
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

// RateLimiter is a simple semaphore for limiting the number of concurrent signatures.
//
// This is a naive implementation that probably requires more complex logic to
// handle real-world scenarios.
//
// Pros:
// - Has simple interface that hides the underlying implementation details.
//
// Cons:
// - Doesn't take into account number of signatures per chain,
// - Doesn't take nonce ordering into account
// - Doesn't take chain-fairness into account
//
// TBD:
// How to ensure that each O+S throttles the same CCTX at a given point in time?
// Otherwise, different nodes might throttle different cctx => no party formed => error
type RateLimiter struct {
	sem     *semaphore.Weighted
	pending *atomic.Int32
}

var ErrThrottled = errors.New("action is throttled")

func New(maxPending int64) *RateLimiter {
	return &RateLimiter{
		sem:     semaphore.NewWeighted(maxPending),
		pending: &atomic.Int32{},
	}
}

func (r *RateLimiter) Acquire(chainID, nonce uint64) error {
	if !r.sem.TryAcquire(1) {
		return errors.Wrapf(ErrThrottled, "chain: %d, nonce: %d", chainID, nonce)
	}

	r.pending.Add(1)

	return nil
}

func (r *RateLimiter) Release(_, _ uint64) {
	// noop
	if r.pending.Load() == 0 {
		return
	}

	r.sem.Release(1)
	r.pending.Add(-1)
}

func (r *RateLimiter) Pending() uint64 {
	return uint64(r.pending.Load())
}
