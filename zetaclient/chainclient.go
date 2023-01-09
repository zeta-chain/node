package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
)

// general chain client

type ChainClient interface {
	GetLastBlockHeight() uint64 // 0 means error
	SetLastBlockHeight(uint64)
	Start()
	Stop()
	GetBaseGasPrice() *big.Int
	PostNonceIfNotRecorded() error
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
	IsSendOutTxProcessed(send *cctxtypes.CrossChainTx) (bool, bool, error)
}
