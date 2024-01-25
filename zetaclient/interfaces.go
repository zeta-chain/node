package zetaclient

import (
	"context"
	"math/big"

	sdkmath "cosmossdk.io/math"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ChainClient is the interface for chain clients
type ChainClient interface {
	Start()
	Stop()
	IsSendOutTxProcessed(sendHash string, nonce uint64, cointype common.CoinType, logger zerolog.Logger) (bool, bool, error)
	SetChainParams(observertypes.ChainParams)
	GetChainParams() observertypes.ChainParams
	GetPromGauge(name string) (prometheus.Gauge, error)
	GetPromCounter(name string) (prometheus.Counter, error)
	GetTxID(nonce uint64) string
	ExternalChainWatcherForNewInboundTrackerSuggestions()
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
	PostVoteInbound(gasLimit, retryGasLimit uint64, msg *crosschaintypes.MsgVoteOnObservedInboundTx) (string, string, error)
	PostVoteOutbound(
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
	) (string, string, error)
	PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error)
	PostAddBlockHeader(chainID int64, txhash []byte, height int64, header common.HeaderData) (string, error)
	GetBlockHeaderStateByChain(chainID int64) (observertypes.QueryGetBlockHeaderStateResponse, error)

	PostBlameData(blame *blame.Blame, chainID int64, index string) (string, error)
	AddTxHashToOutTxTracker(
		chainID int64,
		nonce uint64,
		txHash string,
		proof *common.Proof,
		blockHash string,
		txIndex int64,
	) (string, error)
	GetKeys() *Keys
	GetBlockHeight() (int64, error)
	GetZetaBlockHeight() (int64, error)
	GetLastBlockHeightByChain(chain common.Chain) (*crosschaintypes.LastBlockHeight, error)
	ListPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error)
	GetPendingNoncesByChain(chainID int64) (observertypes.PendingNonces, error)
	GetCctxByNonce(chainID int64, nonce uint64) (*crosschaintypes.CrossChainTx, error)
	GetAllOutTxTrackerByChain(chainID int64, order Order) ([]crosschaintypes.OutTxTracker, error)
	GetCrosschainFlags() (observertypes.CrosschainFlags, error)
	GetObserverList() ([]string, error)
	GetKeyGen() (*observertypes.Keygen, error)
	GetBtcTssAddress() (string, error)
	GetInboundTrackersForChain(chainID int64) ([]crosschaintypes.InTxTracker, error)
	GetLogger() *zerolog.Logger
	ZetaChain() common.Chain
	Pause()
	Unpause()
	GetZetaHotKeyBalance() (sdkmath.Int, error)
}

// BTCRPCClient is the interface for BTC RPC client
type BTCRPCClient interface {
	GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error)
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error)
	EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)
	GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error)
	GetRawTransactionVerbose(txHash *chainhash.Hash) (*btcjson.TxRawResult, error)
	GetBlockCount() (int64, error)
	GetBlockHash(blockHeight int64) (*chainhash.Hash, error)
	GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	GetBlockVerboseTx(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)
	GetBlockHeader(blockHash *chainhash.Hash) (*wire.BlockHeader, error)
}

// EVMRPCClient is the interface for EVM RPC client
type EVMRPCClient interface {
	bind.ContractBackend
	SendTransaction(ctx context.Context, tx *ethtypes.Transaction) error
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethtypes.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethtypes.Header, error)
	TransactionByHash(ctx context.Context, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash ethcommon.Hash) (*ethtypes.Receipt, error)
	TransactionSender(ctx context.Context, tx *ethtypes.Transaction, block ethcommon.Hash, index uint) (ethcommon.Address, error)
}

// KlaytnRPCClient is the interface for Klaytn RPC client
type KlaytnRPCClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*RPCBlock, error)
}
