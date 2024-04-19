package types

import "time"

var DefaultGasPriceIncreaseFlags = GasPriceIncreaseFlags{
	// EpochLength is the number of blocks in an epoch before triggering a gas price increase

	EpochLength: 100,
	// RetryInterval is the number of blocks to wait before incrementing the gas price again
	RetryInterval: time.Minute * 10,

	// GasPriceIncreasePercent is the percentage of median gas price by which to increase the gas price during an increment
	// 100 means the gas price is increased by the median gas price
	GasPriceIncreasePercent: 100,

	// Maximum gas price increase in percent of the median gas price
	// 500 means the gas price can be increased by 5 times the median gas price at most
	GasPriceIncreaseMax: 500,

	// Maximum pending CCTXs to iterate for gas price increase
	MaxPendingCctxs: 500,
}

// DefaultBlockHeaderVerificationFlags returns the default block header verification flags used when not defined
// Deprecated(v16): VerificationFlags are now read in the `lightclient` module
var DefaultBlockHeaderVerificationFlags = BlockHeaderVerificationFlags{
	IsEthTypeChainEnabled: true,
	IsBtcTypeChainEnabled: true,
}

// DefaultCrosschainFlags returns the default crosschain flags used when not defined
func DefaultCrosschainFlags() *CrosschainFlags {
	return &CrosschainFlags{
		IsInboundEnabled:      true,
		IsOutboundEnabled:     true,
		GasPriceIncreaseFlags: &DefaultGasPriceIncreaseFlags,

		// Deprecated(v16): VerificationFlags are now read in the `lightclient` module
		BlockHeaderVerificationFlags: &DefaultBlockHeaderVerificationFlags,
	}
}
