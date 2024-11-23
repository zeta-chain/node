package types

func (o *OutboundTracker) IsMaxed() bool {
	return len(o.HashList) >= MaxOutboundTrackerHashes
}
