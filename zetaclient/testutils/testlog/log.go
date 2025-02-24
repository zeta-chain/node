package testlog

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

type Log struct {
	zerolog.Logger
	t   *testing.T
	buf *bytes.Buffer
	mu  sync.Mutex
}

// New creates a new Log instance with a buffer and a test writer.
func New(t *testing.T) *Log {
	log := &Log{
		t:   t,
		buf: &bytes.Buffer{},
	}

	log.Logger = zerolog.New(log)

	return log
}

func (log *Log) String() string {
	return log.buf.String()
}

func (log *Log) Write(p []byte) (n int, err error) {
	log.mu.Lock()
	defer log.mu.Unlock()

	// silence panics in case this log line is written AFTER test termination.
	const silencePanicSubstring = "Log in goroutine"
	defer func() { silencePanic(recover(), silencePanicSubstring) }()

	// write to the buffer first
	n, err = log.buf.Write(p)
	if err != nil {
		return n, fmt.Errorf("failed to write to buffer: %w", err)
	}

	// Strip trailing newline because t.Log always adds one.
	// (copied from zerolog NewTestWriter)
	p = bytes.TrimRight(p, "\n")

	// Then write to test output
	log.t.Log(string(p))

	return len(p), nil
}

func silencePanic(r any, substr string) {
	// noop
	if r == nil {
		return
	}

	panicStr := fmt.Sprintf("%v", r)
	if strings.Contains(panicStr, substr) {
		return
	}

	panic(panicStr)
}
