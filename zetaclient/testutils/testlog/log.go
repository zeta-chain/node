package testlog

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

type Log struct {
	zerolog.Logger
	buf *concurrentBytesBuffer
}

type concurrentBytesBuffer struct {
	buf *bytes.Buffer
	mu  sync.RWMutex
}

// New creates a new Log instance with a buffer and a test writer.
func New(t *testing.T) *Log {
	buf := &concurrentBytesBuffer{
		buf: &bytes.Buffer{},
		mu:  sync.RWMutex{},
	}

	log := zerolog.New(io.MultiWriter(zerolog.NewTestWriter(t), buf))

	return &Log{Logger: log, buf: buf}
}

func (log *Log) String() string {
	return log.buf.string()
}

func (b *concurrentBytesBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *concurrentBytesBuffer) string() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.buf.String()
}
