package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/big"
)

// general chain client

type ChainClient interface {
	GetLastBlockHeight() uint64 // 0 means error
	SetLastBlockHeight(uint64)
	Start()
	Stop()
	GetBaseGasPrice() *big.Int
	IsSendOutTxProcessed(sendHash string, nonce int, fromOrToZeta bool) (bool, bool, error)
	PostNonceIfNotRecorded() error
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
}
