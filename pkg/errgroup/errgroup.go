// Package errgroup provides synchronization, error propagation, and Context
// cancellation for groups of goroutines working on subtasks of a common task.
//
// It wraps, and exposes a similar API to, the upstream package
// golang.org/x/sync/errgroup.  Our version additionally recovers from panics,
// converting them into errors.
//
// Copyright 2009 The Go Authors and 2021 Steve Coffman
package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	// Sadly, we have to copy the whole implementation, because:
	// - we want a zero errgroup to work, which means we'd need to embed the
	//   upstream errgroup by value
	// - we can't copy an errgroup, which means we can't embed by value
	// (We could get around this with our own initialization-Once, but that
	// seems even more convoluted.) So we just copy -- it's not that much
	// code. The only change below is to add catchPanics(), in Go().
	cancel  func()
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or panics, or the first time Wait returns,
// whichever occurs first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will
// be returned by Wait.
//
// If the function panics, this is treated as if it returned an error.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		// here's the only change from upstream: this was
		//  err := f(); ...
		if err := catchPanics(f)(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

// fromPanicValue takes a value recovered from a panic and converts it into an
// error, for logging purposes. If the value is nil, it returns nil instead of
// an error.
//
// Use like:
//
//	 defer func() {
//			err := fromPanicValue(recover())
//			// log or otherwise use err
//		}()
func fromPanicValue(i interface{}) error {
	switch value := i.(type) {
	case nil:
		return nil
	case string:
		return fmt.Errorf("panic: %v\n%s", value, collectStack())
	case error:
		return fmt.Errorf("panic in errgroup goroutine %w\n%s", value, collectStack())
	default:
		return fmt.Errorf("unknown panic: %+v\n%s", value, collectStack())
	}
}

func collectStack() []byte {
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	return buf
}

func catchPanics(f func() error) func() error {
	return func() (err error) {
		defer func() {
			// modified from log.PanicHandler, except instead of log.Panic we
			// set `err`, which is the named-return from our closure to
			// `g.Group.Go`, to an error based on the panic value.
			// We do not log here -- we are effectively returning the (panic)
			// error to our caller which suffices.
			if r := recover(); r != nil {
				err = fromPanicValue(r)
			}
		}()

		return f()
	}
}
