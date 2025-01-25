package fanout

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFanOut(t *testing.T) {
	// ARRANGE
	// Given an input
	input := make(chan int)

	// Given a fanout
	f := New(input, DefaultBuffer)

	// That has 3 outputs
	out1, _ := f.Add()
	out2, _ := f.Add()
	out3, _ := f.Add()

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

func TestFanOutClose(t *testing.T) {
	// ARRANGE
	// Given input
	input := make(chan int64)

	f := New(input, DefaultBuffer)

	// Given 2 channels
	out4, _ := f.Add()
	out5, close5 := f.Add()

	// Given total counter
	var total int64

	// ACT
	f.Start()

	// Write to the input
	go func() {
		for i := int64(0); i < 10; i++ {
			input <- i
			time.Sleep(15 * time.Millisecond)
		}
		close(input)
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// Read from chan4 unless input is closed
	go func() {
		defer wg.Done()

		for i := range out4 {
			t.Logf("out4: received %d", i)
			atomic.AddInt64(&total, 1)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Read 5 items from chan5 and then close it
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			v := <-out5
			t.Logf("out5: received %d", v)
			atomic.AddInt64(&total, 1)
			time.Sleep(10 * time.Millisecond)
		}

		// after 5 iterations, close out5
		t.Logf("out5: closing")
		close5()
	}()

	wg.Wait()

	// ASSERT
	require.Equal(t, int64(10+5), total)
}

// go test -v -bench=BenchmarkFanOutMemoryUsage -benchmem -benchtime=30s
func BenchmarkFanOutMemoryUsage(b *testing.B) {
	b.ReportAllocs()

	const (
		outputBufferSize = 16
		numOutputs       = 3
		numMessages      = 100
	)

	runtime.GC()

	for iter := 0; iter < b.N; iter++ {
		// ARRANGE
		// Given an input channel & fanout
		input := make(chan int)
		f := New(input, outputBufferSize)

		// Given counter
		var counter uint64

		// Given outputs what simply consume data
		var wg sync.WaitGroup
		for i := 0; i < numOutputs; i++ {
			out, _ := f.Add()

			wg.Add(1)
			go func(out <-chan int, i int) {
				defer wg.Done()
				for range out {
					atomic.AddUint64(&counter, 1)
				}
			}(out, i)
		}

		// Start fanout
		f.Start()

		// Given mem stats
		var memStatsBefore runtime.MemStats
		runtime.ReadMemStats(&memStatsBefore)

		// ACT
		// Send messages to the input channel
		for i := 0; i < numMessages; i++ {
			input <- i
		}

		// Then close input channel to stop the fan-out process
		close(input)

		// Wait for consumers to finish
		wg.Wait()

		// ASSERT
		// Check counter. We don't have guarantees that after (input.close)
		assert.NotZero(b, counter)

		// Track memory usage after processing messages
		var memStatsAfter runtime.MemStats
		runtime.ReadMemStats(&memStatsAfter)

		logMem(b, &memStatsBefore, &memStatsAfter)
	}
}

func logMem(t testing.TB, before, after *runtime.MemStats) {
	t.Logf(
		"Mem Before: Alloc = %d KB, TotalAlloc = %d KB, HeapAlloc = %d KB, HeapInUse = %d, HeapSys = %d KB",
		before.Alloc/1024,
		before.TotalAlloc/1024,
		before.HeapAlloc/1024,
		before.HeapInuse/1024,
		before.HeapSys/1024,
	)

	t.Logf(
		"Mem After: Alloc = %d KB, TotalAlloc = %d KB, HeapAlloc = %d KB, HeapInUse = %d KB, HeapSys = %d KB",
		after.Alloc/1024,
		after.TotalAlloc/1024,
		after.HeapAlloc/1024,
		before.HeapInuse/1024,
		after.HeapSys/1024,
	)
}
