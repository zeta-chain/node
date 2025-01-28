package common

import "time"

const (
	// DefaultGasPriceMultiplierFeeCharge is the default gas price multiplier to charge fees from users
	DefaultGasPriceMultiplierFeeCharge = 1.0

	// EVMGasPriceMultiplierFeeCharge is the default gas price multiplier to charge fees from users
	EVMGasPriceMultiplierFeeCharge = 1.2

	// BTCGasPriceMultiplierFeeCharge is the default gas price multiplier to charge fees from users
	BTCGasPriceMultiplierFeeCharge = 2.0

	// BTCGasPriceMultiplierSendTx is the default gas price multiplier to send out BTC TSS txs
	BTCGasPriceMultiplierSendTx = 1.5

	// RPCStatusCheckInterval is the interval to check RPC status, 1 minute
	RPCStatusCheckInterval = time.Minute

	// BTCMempoolStuckTxCheckInterval is the interval to check for Bitcoin stuck transactions in the mempool
	BTCMempoolStuckTxCheckInterval = 30 * time.Second
)
