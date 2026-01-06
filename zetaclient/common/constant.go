package common

import "time"

const (
	// RPCStatusCheckInterval is the interval to check RPC status, 1 minute
	RPCStatusCheckInterval = time.Minute

	// BTCMempoolStuckTxCheckInterval is the interval to check for Bitcoin stuck transactions in the mempool
	BTCMempoolStuckTxCheckInterval = 30 * time.Second
)
