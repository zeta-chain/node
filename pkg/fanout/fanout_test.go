package fanout

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFanOut(t *testing.T) {
	// ARRANGE
	// Given an input
	input := make(chan int)

	// Given a fanout
	f := New(input, DefaultBuffer)

	// That has 3 outputs
	out1 := f.Add()
	out2 := f.Add()
	out3 := f.Add()

	// Given a wait group
	wg := sync.WaitGroup{}
	wg.Add(3)

	// Given a sample number
	var total int32

	// Given a consumer
	consumer := func(out <-chan int, name string, lag time.Duration) {
		defer wg.Done()
		var local int32
		for i := range out {
			// simulate some work
			time.Sleep(lag)

			local += int32(i)
			t.Logf("%s: received %d", name, i)
		}

		// add only if input was closed
		atomic.AddInt32(&total, local)
	}

	// ACT
	f.Start()

	// Write to the channel
	go func() {
		for i := 1; i <= 10; i++ {
			input <- i
			t.Logf("fan-out: sent %d", i)
			time.Sleep(50 * time.Millisecond)
		}

		close(input)
	}()

	go consumer(out1, "out1: fast consumer", 10*time.Millisecond)
	go consumer(out2, "out2: average consumer", 60*time.Millisecond)
	go consumer(out3, "out3: slow consumer", 150*time.Millisecond)

	wg.Wait()

	// ASSERT
	// Check that total is valid
	// total == sum(1...10) * 3 = n(n+1)/2 * 3 = 55 * 3 = 165
	require.Equal(t, int32(165), total)
}
