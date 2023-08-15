package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// general chain client

type ChainClient interface {
	GetLastBlockHeight() int64 // 0 means error
	SetLastBlockHeight(int64)
	Start()
	Stop()
	//GetBaseGasPrice() *big.Int
	IsSendOutTxProcessed(sendHash string, nonce uint64, cointype common.CoinType, logger zerolog.Logger) (bool, bool, error)
	//PostNonceIfNotRecorded(logger zerolog.Logger) error
	GetCoreParameters() config.CoreParams
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
	GetTxID(nonce uint64) string
}
