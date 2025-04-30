// Package fanout provides a fan-out pattern implementation.
// It allows one channel to stream data to multiple independent channels.
// Note that context handling is out of the scope of this package.
package fanout

import (
	"slices"
	"sync"
	"sync/atomic"
)

const DefaultBuffer = 8

// FanOut is a fan-out pattern implementation.
// It is NOT a worker pool, so use it wisely.
type FanOut[T any] struct {
	input   <-chan T
	outputs []*output[T]

	// outputBuffer chan buffer size for outputs channels.
	// This helps with writing to chan in case of slow consumers.
	outputBuffer int

	mu sync.RWMutex
}

type output[T any] struct {
	ch     chan T
	status atomic.Int32
	wg     sync.WaitGroup
}

const (
	statusRunning = int32(0)
	statusClosing = int32(1)
	statusClosed  = int32(2)
)

// New constructs FanOut
func New[T any](source <-chan T, buf int) *FanOut[T] {
	return &FanOut[T]{
		input:        source,
		outputs:      make([]*output[T], 0),
		outputBuffer: buf,
	}
}

// Add adds a new output channel to the fan-out.
// Returns the output channel and a close function.
func (f *FanOut[T]) Add() (<-chan T, func()) {
	out := &output[T]{
		ch: make(chan T, f.outputBuffer),
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.outputs = append(f.outputs, out)

	return out.ch, func() { f.remove(out) }
}

func (f *FanOut[T]) remove(out *output[T]) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := range f.outputs {
		if f.outputs[i].equal(out) {
			// cut item from the slice
			f.outputs = slices.Delete(f.outputs, i, i+1)
			out.close()
			return
		}
	}
}

// Start starts the fan-out process
func (f *FanOut[T]) Start() {
	go func() {
		// loop for new data
		for data := range f.input {
			f.mu.RLock()
			for i := range f.outputs {
				// It's a naive approach but should be more than enough for our use cases.
				//
				// note (a): this might spawn lots of goroutines.
				//
				// note (b): it does NOT guarantee the order of messages:
				// imagine f.input receives 5 msgs/sec *at peak*,
				// but the output processes only 1 msg/sec, thus +4 goroutines will be spawned
				// => no control over the order.
				f.outputs[i].write(data)
			}
			f.mu.RUnlock()
		}

		// at this point, the input was closed
		f.mu.Lock()
		defer f.mu.Unlock()
		for _, out := range f.outputs {
			out.close()
		}

		f.outputs = nil
	}()
}

func (o *output[T]) write(data T) {
	o.wg.Add(1)

	go func() {
		defer o.wg.Done()
		if o.isRunning() {
			o.ch <- data
		}
	}()
}

func (o *output[T]) equal(item *output[T]) bool {
	// channels are equal if they refer to the same instance
	return o.ch == item.ch
}

func (o *output[T]) close() {
	// noop
	if !o.isRunning() {
		return
	}

	o.status.Store(statusClosing)
	o.wg.Wait()

	o.status.Store(statusClosed)
	close(o.ch)
}

func (o *output[T]) isRunning() bool {
	return o.status.Load() == statusRunning
}
