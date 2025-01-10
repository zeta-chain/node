package types

// MaxReached returns true if the OutboundTracker has reached the maximum number of hashes it can store.
func (o *OutboundTracker) MaxReached() bool {
	return len(o.HashList) >= MaxOutboundTrackerHashes
}
