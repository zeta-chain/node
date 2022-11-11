package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/big"
)

// general chain client

type ChainClient interface {
	GetBlockHeight() uint64 // 0 means error
	Start()
	Stop()
	GetBaseGasPrice() *big.Int
	IsSendOutTxProcessed(sendHash string, nonce int, fromOrToZeta bool) (bool, bool, error)
	PostNonceIfNotRecorded() error
	GetPromGauge(name string) (prometheus.Gauge, error)
	RegisterPromGauge(name string, help string) error
}
