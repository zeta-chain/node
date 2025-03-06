package errgroup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPanicWithString(t *testing.T) {
	g, ctx := WithContext(context.Background())
	g.Go(func() error { panic("oh noes") })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	require.Error(t, err)
	require.Contains(t, err.Error(), "oh noes")

	// ctx should now be canceled.
	require.Error(t, ctx.Err())
}

func TestPanicWithError(t *testing.T) {
	g, ctx := WithContext(context.Background())

	panicErr := errors.New("oh noes")
	g.Go(func() error { panic(panicErr) })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	require.Error(t, err)
	require.Contains(t, err.Error(), "oh noes")
	require.True(t, errors.Is(err, panicErr))

	// ctx should now be canceled.
	require.Error(t, ctx.Err())
}

func TestPanicWithOtherValue(t *testing.T) {
	g, ctx := WithContext(context.Background())

	panicVal := struct {
		int
		string
	}{1234567890, "oh noes"}
	g.Go(func() error { panic(panicVal) })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	require.Error(t, err)
	require.Contains(t, err.Error(), "oh noes")
	require.Contains(t, err.Error(), "1234567890")

	// ctx should now be canceled.
	require.Error(t, ctx.Err())
}

func TestError(t *testing.T) {
	g, ctx := WithContext(context.Background())

	goroutineErr := errors.New("oh noes")
	g.Go(func() error { return goroutineErr })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	require.Error(t, err)
	require.Contains(t, err.Error(), "oh noes")
	require.True(t, errors.Is(err, goroutineErr))

	// ctx should now be canceled.
	require.Error(t, ctx.Err())
}

func TestSuccess(t *testing.T) {
	g, ctx := WithContext(context.Background())

	g.Go(func() error { return nil })
	// since no goroutine errored, ctx.Err() should be nil
	// (until all goroutines are done)
	g.Go(ctx.Err)

	err := g.Wait()
	require.NoError(t, err)

	// ctx should now still be canceled.
	require.Error(t, ctx.Err())
}

func TestManyGoroutines(t *testing.T) {
	n := 100
	g, ctx := WithContext(context.Background())

	for i := 0; i < n; i++ {
		// put in a bunch of goroutines that just return right away
		g.Go(func() error { return nil })
		// and also a bunch that wait for the error
		g.Go(func() error {
			<-ctx.Done()
			return ctx.Err()
		})
	}

	// finally, put in a panic
	g.Go(func() error { panic("oh noes") })

	// as before, Wait() will finish only once all goroutines do, but returns
	// the first error (namely the panic)
	err := g.Wait()
	require.Error(t, err)
	require.Contains(t, err.Error(), "oh noes")

	// ctx should now be canceled.
	require.Error(t, ctx.Err())
}

func TestZeroGroupPanic(t *testing.T) {
	var g Group

	// either of these could happen first, since a zero group does not cancel
	g.Go(func() error { panic("oh noes") })
	g.Go(func() error { return nil })

	// Wait() still returns the error.
	err := g.Wait()
	require.Error(t, err)
	require.Contains(t, err.Error(), "oh noes")
}

func TestZeroGroupSuccess(t *testing.T) {
	var g Group

	g.Go(func() error { return nil })
	g.Go(func() error { return nil })

	err := g.Wait()
	require.NoError(t, err)
}
