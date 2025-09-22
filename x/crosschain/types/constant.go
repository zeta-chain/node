package types

const (
	// CCTXIndexLength is the length of a crosschain transaction index
	CCTXIndexLength = 66
	// MaxOutboundTrackerHashes is the maximum number of outbound tracker hashes to store in the hashlist for an outbound tracker
	MaxOutboundTrackerHashes = 5

	// UsableRemainingFeesPercentage UnableToPayFeesPercentage is the percentage of fees that is considered as unable.
	// A portion of the fees is reserved for refunds, and the rest is used for funding the stability pool
	UsableRemainingFeesPercentage = uint64(95)

	// DefaultStabilityPoolFundPercentage is the default percentage of fees allocated to the stability pool
	DefaultStabilityPoolFundPercentage = uint64(100)
)
