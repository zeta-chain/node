// Package retry provides a generic retry mechanism with exponential backoff.
//
// Example:
//
//	ctx := context.Background()
//	client := foobar.NewClient()
//
//	err := retry.Do(func() error {
//		err := client.UpdateConfig(ctx, map[string]any{"key": "value"})
//
//		// will be retied
//		if errors.Is(err, foobar.ErrTxConflict) {
//			return retry.Retryable(err)
//		}
//
//		return err
//	})
package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

type Backoff = backoff.BackOff

type Callback func() error

type TypedCallback[T any] func() (T, error)

type errRetryable struct {
	error
}

// DefaultBackoff returns a default backoff strategy with 5 exponential retries.
func DefaultBackoff() Backoff {
	bo := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(50*time.Millisecond),
		backoff.WithMaxInterval(500*time.Millisecond),
		backoff.WithMultiplier(1.8),
	)

	return backoff.WithMaxRetries(bo, 5)
}

// Do executes the callback function with the default backoff config.
// It will retry a callback ONLY if error is retryable.
func Do(cb Callback) error {
	return DoWithBackoff(cb, DefaultBackoff())
}

// DoWithBackoff executes the callback function with provided backoff config.
// It will retry a callback ONLY if error is retryable.
func DoWithBackoff(cb Callback, bo Backoff) error {
	for {
		err := cb()
		if err == nil {
			return nil
		}

		var errRetry errRetryable
		isRetryable := errors.As(err, &errRetry)
		if !isRetryable {
			return err
		}

		sleepFor := bo.NextBackOff()
		if sleepFor == backoff.Stop {
			return errors.Wrap(err, "retry limit exceeded")
		}

		time.Sleep(sleepFor)
	}
}

// DoTyped is typed version of Do that returns a value along with an error.
// It will retry a callback ONLY if error is retryable.
func DoTyped[T any](cb TypedCallback[T]) (T, error) {
	return DoTypedWithBackoff(cb, DefaultBackoff())
}

// DoTypedWithBackoff is typed version of DoWithBackoff that returns a value along with an error.
// It will retry a callback ONLY if error is retryable.
func DoTypedWithBackoff[T any](cb TypedCallback[T], bo Backoff) (T, error) {
	var (
		result T
		err    error
	)

	// #nosec G703 error is propagated
	_ = DoWithBackoff(func() error {
		result, err = cb()
		return err
	}, bo)

	return result, err
}

// DoTypedWithRetry is DoTyped but ANY error is retried.
func DoTypedWithRetry[T any](cb TypedCallback[T]) (T, error) {
	wrapper := func() (T, error) {
		return RetryTyped(cb())
	}

	return DoTypedWithBackoffAndRetry(wrapper, DefaultBackoff())
}

// DoTypedWithBackoffAndRetry is DoTypedWithBackoff but ANY error is retried.
func DoTypedWithBackoffAndRetry[T any](cb TypedCallback[T], bo Backoff) (T, error) {
	wrapper := func() (T, error) {
		return RetryTyped(cb())
	}

	return DoTypedWithBackoff(wrapper, bo)
}

// Retry wraps error to mark it as retryable. Skips retry for context errors.
func Retry(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		// do not retry context errors
		return err
	default:
		return errRetryable{error: err}
	}
}

// RetryTyped wraps error to mark it as retryable
//
//goland:noinspection GoNameStartsWithPackageName
//nolint:revive
func RetryTyped[T any](result T, err error) (T, error) {
	if err == nil {
		return result, nil
	}

	return result, Retry(err)
}
