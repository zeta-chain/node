package types

const (
	// CCTXIndexLength is the length of a crosschain transaction index
	CCTXIndexLength          = 66
	MaxOutboundTrackerHashes = 5

	// UsableRemainingFeesPercentage UnableToPayFeesPercentage is the percentage of fees that is considered as unable.
	// A portion of the fees is reserved for refunds, and the rest is used for funding the stability pool
	UsableRemainingFeesPercentage = uint64(95)
)
