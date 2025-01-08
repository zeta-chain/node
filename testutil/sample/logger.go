package sample

import (
	"bytes"

	"github.com/cometbft/cometbft/libs/log"
)

type TestLogger struct {
	buf bytes.Buffer
	log.Logger
}

func NewTestLogger() *TestLogger {
	tl := &TestLogger{}
	tl.Logger = log.NewTMLogger(log.NewSyncWriter(&tl.buf))
	return tl
}

func (t *TestLogger) Write(p []byte) (n int, err error) {
	return t.buf.Write(p)
}

func (t *TestLogger) String() string {
	return t.buf.String()
}
