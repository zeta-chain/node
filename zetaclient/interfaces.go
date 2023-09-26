package zetaclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
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
		send *crosschaintypes.CrossChainTx,
		outTxMan *OutTxProcessorManager,
		outTxID string,
		evmClient ChainClient,
		zetaBridge ZetaCoreBridger,
		height uint64,
	)
}

// ZetaCoreBridger is the interface to interact with ZetaCore
type ZetaCoreBridger interface {
	AddTxHashToOutTxTracker(
		chainID int64,
		nonce uint64,
		txHash string,
		proof *common.Proof,
		blockHash string,
		txIndex int64,
	) (string, error)
	GetKeys() *Keys
	GetZetaBlockHeight() (int64, error)
	GetAllPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, error)
	GetAllOutTxTrackerByChain(chain common.Chain, order Order) ([]crosschaintypes.OutTxTracker, error)
	GetCrosschainFlags() (observertypes.CrosschainFlags, error)
	GetObserverList(chain common.Chain) ([]string, error)
	Pause()
	Unpause()
}
