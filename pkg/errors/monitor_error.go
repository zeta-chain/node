package errors

// MonitorError represents an error from the monitoring goroutine.
type MonitorError struct {
	Err                error
	InboundBlockHeight uint64
	ZetaTxHash         string
	BallotIndex        string
}
