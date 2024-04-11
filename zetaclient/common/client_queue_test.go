package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func setupClientQueue() *ClientQueue {
	clientQ := NewClientQueue()
	clientQ.Append(1)
	clientQ.Append(2)
	clientQ.Append(3)
	return clientQ
}

func TestClientQueue_Append(t *testing.T) {
	clientQ := setupClientQueue()
	clientQ.Append(2343)
	require.Equal(t, 4, clientQ.Length())
}

func TestClientQueue_First(t *testing.T) {
	clientQ := setupClientQueue()
	first := clientQ.First()
	require.Equal(t, 1, first)
}

func TestClientQueue_Length(t *testing.T) {
	clientQ := setupClientQueue()
	clientQ.Append(4)
	clientQ.Append(5)
	clientQ.Append(6)
	length := clientQ.Length()
	require.Equal(t, 6, length)
}

func TestClientQueue_Next(t *testing.T) {
	clientQ := setupClientQueue()
	require.Equal(t, 1, clientQ.First())
	clientQ.Next()
	require.Equal(t, 2, clientQ.First())
	clientQ.Next()
	require.Equal(t, 3, clientQ.First())
}

func TestClientQueue_Empty(t *testing.T) {
	clientQ := NewClientQueue()
	first := clientQ.First()
	require.Equal(t, nil, first)
	length := clientQ.Length()
	require.Equal(t, 0, length)

	//This shouldn't panic
	clientQ.Next()
	clientQ.Next()
	clientQ.Next()
}
