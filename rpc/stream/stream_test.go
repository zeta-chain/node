package stream

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStreamAdd(t *testing.T) {
	testCases := []struct {
		segmentSize, capacity int
	}{
		{128, 1280},
		{1024, 2048},
		{1024, 2148},
		{1024, 1024},
		{2048, 100},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("segmentSize=%d,capacity=%d", tc.segmentSize, tc.capacity)
		t.Run(name, func(t *testing.T) {
			stream := NewStream[int](tc.segmentSize, tc.capacity)

			amount := 100000
			for i := 0; i < amount; i++ {
				require.Equal(t, i+1, stream.Add(i))
			}

			all, _ := stream.ReadAllNonBlocking(0)
			maxSegments := (tc.capacity + tc.segmentSize - 1) / tc.segmentSize
			require.Equal(t, maxSegments*tc.segmentSize+amount%tc.segmentSize, len(all))
			require.Equal(t, 100000-1, all[len(all)-1])
			for i, n := range all[:len(all)-1] {
				require.Equal(t, n+1, all[i+1])
			}
		})
	}
}

func TestStreamReadNonBlocking(t *testing.T) {
	stream := NewStream[int](16, 31)

	for i := 0; i < 32; i++ {
		require.Equal(t, i+1, stream.Add(i))
	}

	items, offset := stream.ReadNonBlocking(0)
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, items)
	require.Equal(t, 16, offset)
}

func TestStreamReadBlocking(t *testing.T) {
	stream := NewStream[int](16, 31)

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	// subscriber
	subscribers := 10
	result := make([][]int, subscribers)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			require.NoError(t, stream.Subscribe(ctx, func(items []int, offset int) error {
				result[i] = append(result[i], items...)
				return nil
			}))
		}(i)
	}

	// wait for subscribers to setup
	time.Sleep(100 * time.Millisecond)

	// publisher
	for i := 0; i < 32; i++ {
		require.Equal(t, i+1, stream.Add(i))
	}

	// wait for subscribers to finish
	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	// check result
	for i := 0; i < subscribers; i++ {
		require.Equal(t, 32, len(result[i]))
		require.Equal(t, 31, result[i][len(result[i])-1])
		for j, n := range result[i][:len(result[i])-1] {
			require.Equal(t, n+1, result[i][j+1])
		}
	}
}

func TestStreamReadFromEnd(t *testing.T) {
	stream := NewStream[int](16, 31)

	items, offset := stream.ReadNonBlocking(-1)
	require.Empty(t, items)
	require.Equal(t, 0, offset)

	stream.Add(1)

	items, offset = stream.ReadNonBlocking(-1)
	require.Empty(t, items)
	require.Equal(t, 1, offset)

	stream.Add(2)

	items, offset = stream.ReadNonBlocking(offset)
	require.Equal(t, []int{2}, items)
	require.Equal(t, 2, offset)
}
