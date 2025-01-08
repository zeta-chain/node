// Package fanout provides a fan-out pattern implementation.
// It allows one channel to stream data to multiple independent channels.
// Note that context handling is out of the scope of this package.
package fanout

import "sync"

const DefaultBuffer = 8

// FanOut is a fan-out pattern implementation.
// It is NOT a worker pool, so use it wisely.
type FanOut[T any] struct {
	input   <-chan T
	outputs []chan T

	// outputBuffer chan buffer size for outputs channels.
	// This helps with writing to chan in case of slow consumers.
	outputBuffer int

	mu sync.RWMutex
}

// New constructs FanOut
func New[T any](source <-chan T, buf int) *FanOut[T] {
	return &FanOut[T]{
		input:        source,
		outputs:      make([]chan T, 0),
		outputBuffer: buf,
	}
}

func (f *FanOut[T]) Add() <-chan T {
	out := make(chan T, f.outputBuffer)

	f.mu.Lock()
	defer f.mu.Unlock()

	f.outputs = append(f.outputs, out)

	return out
}

// Start starts the fan-out process
func (f *FanOut[T]) Start() {
	go func() {
		// loop for new data
		for data := range f.input {
			f.mu.RLock()
			for _, output := range f.outputs {
				// note that this might spawn lots of goroutines.
				// it is a naive approach, but should be more than enough for our use cases.
				go func(output chan<- T) { output <- data }(output)
			}
			f.mu.RUnlock()
		}

		// at this point, the input was closed
		f.mu.Lock()
		defer f.mu.Unlock()
		for _, out := range f.outputs {
			close(out)
		}

		f.outputs = nil
	}()
}
