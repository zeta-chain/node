package common

import "time"

const (
	// DefaultGasPriceMultiplier is the default gas price multiplier for all chains
	DefaultGasPriceMultiplier = 1.0

	// EVMOutboundGasPriceMultiplier is the default gas price multiplier for EVM-chain outbond txs
	EVMOutboundGasPriceMultiplier = 1.2

	// BTCOutboundGasPriceMultiplier is the default gas price multiplier for BTC outbond txs
	BTCOutboundGasPriceMultiplier = 2.0

	// RPCStatusCheckInterval is the interval to check RPC status, 1 minute
	RPCStatusCheckInterval = time.Minute

	// BTCMempoolStuckTxCheckInterval is the interval to check for Bitcoin stuck transactions in the mempool
	BTCMempoolStuckTxCheckInterval = 30 * time.Second
)
