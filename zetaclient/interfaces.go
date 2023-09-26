package zetaclient

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
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
	PostSend(
		sender string,
		senderChain int64,
		txOrigin string,
		receiver string,
		receiverChain int64,
		amount math.Uint,
		message string,
		inTxHash string,
		inBlockHeight uint64,
		gasLimit uint64,
		coinType common.CoinType,
		zetaGasLimit uint64,
		asset string,
	) (string, error)
	PostReceiveConfirmation(
		sendHash string,
		outTxHash string,
		outBlockHeight uint64,
		outTxGasUsed uint64,
		outTxEffectiveGasPrice *big.Int,
		outTxEffectiveGasLimit uint64,
		amount *big.Int,
		status common.ReceiveStatus,
		chain common.Chain,
		nonce uint64,
		coinType common.CoinType,
	) (string, error)
	PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error)
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
	GetPendingNoncesByChain(chainID int64) (crosschaintypes.PendingNonces, error)
	GetCctxByNonce(chainID int64, nonce uint64) (*crosschaintypes.CrossChainTx, error)
	GetAllOutTxTrackerByChain(chain common.Chain, order Order) ([]crosschaintypes.OutTxTracker, error)
	GetCrosschainFlags() (observertypes.CrosschainFlags, error)
	GetObserverList(chain common.Chain) ([]string, error)
	Pause()
	Unpause()
}

// BTCRPCClient is the interface for BTC RPC client
type BTCRPCClient interface {
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error)
	EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)
	GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error)
	GetRawTransactionVerbose(txHash *chainhash.Hash) (*btcjson.TxRawResult, error)
	GetBlockCount() (int64, error)
	GetBlockHash(blockHeight int64) (*chainhash.Hash, error)
	GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	GetBlockVerboseTx(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)
}
