package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/common"
	"math/big"
)

// general chain client

type ChainClient interface {
	GetLastBlockHeight() int64 // 0 means error
	SetLastBlockHeight(int64)
	Start()
	Stop()
	GetBaseGasPrice() *big.Int
	IsSendOutTxProcessed(sendHash string, nonce int, cointype common.CoinType) (bool, bool, error)
	PostNonceIfNotRecorded() error
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
}
