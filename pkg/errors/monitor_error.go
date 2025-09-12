package errors

import "fmt"

// ErrTxMonitor represents an error from the monitoring goroutine.
type ErrTxMonitor struct {
	Err                error
	InboundBlockHeight uint64
	ZetaTxHash         string
	BallotIndex        string
}

func (m ErrTxMonitor) Error() string {
	if m.Err == nil {
		return "monitoring completed without error"
	}
	return fmt.Sprintf("monitoring error: %v", m.Err)
}
