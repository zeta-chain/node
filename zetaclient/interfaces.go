package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ChainClient is the interface for chain clients
type ChainClient interface {
	Start()
	Stop()
	IsSendOutTxProcessed(sendHash string, nonce uint64, cointype common.CoinType, logger zerolog.Logger) (bool, bool, error)
	SetCoreParams(observertypes.CoreParams)
	GetCoreParams() observertypes.CoreParams
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
	GetTxID(nonce uint64) string
}

// ChainSigner is the interface to sign transactions for a chain
type ChainSigner interface {
	TryProcessOutTx(
		send *types.CrossChainTx,
		outTxMan *OutTxProcessorManager,
		outTxID string,
		evmClient ChainClient,
		zetaBridge *ZetaCoreBridge,
		height uint64,
	)
}
