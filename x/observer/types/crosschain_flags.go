package types

import "time"

var DefaultGasPriceIncreaseFlags = GasPriceIncreaseFlags{
	// EpochLength is the number of blocks in an epoch before triggering a gas price increase
	EpochLength: 100,

	// RetryInterval is the time to wait before incrementing the gas price again
	RetryInterval: time.Minute * 10,

	// GasPriceIncreasePercent is the percentage of median gas price by which to increase the gas price during an increment
	// 100 means the gas price is increased by the median gas price
	GasPriceIncreasePercent: 100,

	// Maximum gas price increase in percent of the median gas price
	// 500 means the gas price can be increased by 5 times the median gas price at most
	GasPriceIncreaseMax: 500,

	// Maximum pending CCTXs to iterate for gas price increase
	MaxPendingCctxs: 500,

	// RetryIntervalBTC is the time to wait before incrementing the gas price again for Bitcoin chain
	// Bitcoin block time is 10 minutes, so the interval needs to be longer than confirmation time (10 min * 3 blocks).
	// 40 minutes is chosen to be a mild interval to balance the finality and the sensitivity to fee marget changes.
	RetryIntervalBTC: time.Minute * 40,
}

// DefaultCrosschainFlags returns the default crosschain flags used when not defined
func DefaultCrosschainFlags() *CrosschainFlags {
	return &CrosschainFlags{
		IsInboundEnabled:      true,
		IsOutboundEnabled:     true,
		GasPriceIncreaseFlags: &DefaultGasPriceIncreaseFlags,
		IsV2ZetaEnabled:       false,
	}
}
