package sample

import (
	"bytes"

	"cosmossdk.io/log"
	log2 "github.com/cometbft/cometbft/libs/log"
	"github.com/rs/zerolog"
)

type TestLogger struct {
	buf bytes.Buffer
	log.Logger
}

func NewTestLogger() *TestLogger {
	tl := &TestLogger{}
	// TODO: simplify?
	tl.Logger = log.NewLogger(zerolog.New(log2.NewSyncWriter(&tl.buf)))
	return tl
}

func (t *TestLogger) Write(p []byte) (n int, err error) {
	return t.buf.Write(p)
}

func (t *TestLogger) String() string {
	return t.buf.String()
}
