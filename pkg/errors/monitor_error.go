package errors

import "fmt"

// MonitorError represents an error from the monitoring goroutine.
type MonitorError struct {
	Err                error
	InboundBlockHeight uint64
	ZetaTxHash         string
	BallotIndex        string
}

func (m MonitorError) Error() string {
	if m.Err == nil {
		return "monitoring completed without error"
	}
	return fmt.Sprintf("monitoring error: %v", m.Err)
}
