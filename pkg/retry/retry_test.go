package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		err := Do(func() error { return nil })

		assert.NoError(t, err)
	})

	t.Run("non-retryable error", func(t *testing.T) {
		var counter int
		err := Do(func() error {
			counter++
			return errors.New("something went wrong")
		})

		assert.Equal(t, 1, counter)
		assert.ErrorContains(t, err, "something went wrong")
	})

	t.Run("retryable error suddenly became non-retryable", func(t *testing.T) {
		var counter int
		err := Do(func() error {
			err := errors.New("something went wrong")

			counter++
			if counter < 3 {
				return Retry(err)
			}

			return err
		})

		assert.Equal(t, 3, counter)
		assert.ErrorContains(t, err, "something went wrong")
	})

	t.Run("retryable code eventually works", func(t *testing.T) {
		var counter int
		err := Do(func() error {
			err := errors.New("something went wrong")

			counter++
			if counter < 3 {
				return Retry(err)
			}

			return nil
		})

		assert.Equal(t, 3, counter)
		assert.NoError(t, err)
	})

	t.Run("retry limit exceeded", func(t *testing.T) {
		start := time.Now()

		var counter int
		err := Do(func() error {
			trackTime(t, start)
			err := errors.New("something went wrong")

			counter++
			return Retry(err)
		})

		assert.ErrorContains(t, err, "retry limit exceeded")
	})

	t.Run("context errors are non-retryable", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		var counter int
		err := Do(func() error {
			time.Sleep(100 * time.Millisecond)

			if err := ctx.Err(); err != nil {
				return err
			}

			counter++

			return nil
		})

		assert.Equal(t, 0, counter)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

func TestDoTyped(t *testing.T) {
	t.Parallel()

	type myType struct {
		Value string
	}

	t.Run("no error", func(t *testing.T) {
		result, err := DoTyped(func() (myType, error) {
			return myType{Value: "abc"}, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, "abc", result.Value)
	})

	t.Run("fails", func(t *testing.T) {
		var counter int

		result, err := DoTyped(func() (myType, error) {
			counter++
			return myType{}, errors.New("something went wrong")
		})

		assert.ErrorContains(t, err, "something went wrong")
		assert.Empty(t, result)
		assert.Equal(t, 1, counter)
	})

	t.Run("recovers", func(t *testing.T) {
		var counter int

		result, err := DoTyped(func() (myType, error) {
			counter++
			if counter == 4 {
				return myType{Value: "abc"}, nil
			}
			return myType{}, Retry(errors.New("something went wrong"))
		})

		assert.NoError(t, err)
		assert.Equal(t, "abc", result.Value)
		assert.Equal(t, 4, counter)
	})
}

func trackTime(t *testing.T, from time.Time) {
	duration := time.Since(from)

	t.Logf("Retrier invokation: t = %dms", duration.Milliseconds())
}
