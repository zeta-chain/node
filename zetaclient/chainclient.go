package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/common"
	"math/big"
)

// general chain client

type ChainClient interface {
	GetLastBlockHeight() uint64 // 0 means error
	SetLastBlockHeight(uint64)
	Start()
	Stop()
	GetBaseGasPrice() *big.Int
	IsSendOutTxProcessed(sendHash string, nonce int, coinType common.CoinType) (bool, bool, error)
	PostNonceIfNotRecorded() error
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
}
