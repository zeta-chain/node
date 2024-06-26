package common

import "time"

const (
	// ZetaBlockTime is the block time of the ZetaChain
	ZetaBlockTime = 6000 * time.Millisecond

	// EVMOutboundGasPriceMultiplier is the default gas price multiplier for EVM-chain outbond txs
	EVMOutboundGasPriceMultiplier = 1.2

	// BTCOutboundGasPriceMultiplier is the default gas price multiplier for BTC outbond txs
	BTCOutboundGasPriceMultiplier = 2.0
)
