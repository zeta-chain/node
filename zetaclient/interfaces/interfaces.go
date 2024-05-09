package interfaces

import (
	"context"
	"math/big"

	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"

	"github.com/onrik/ethrpc"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"

	sdkmath "cosmossdk.io/math"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type Order string

const (
	NoOrder    Order = ""
	Ascending  Order = "ASC"
	Descending Order = "DESC"
)

// ChainClient is the interface for chain clients
type ChainClient interface {
	Start()
	Stop()
	IsOutboundProcessed(cctx *crosschaintypes.CrossChainTx, logger zerolog.Logger) (bool, bool, error)
	SetChainParams(observertypes.ChainParams)
	GetChainParams() observertypes.ChainParams
	GetTxID(nonce uint64) string
	WatchInboundTracker()
}

// ChainSigner is the interface to sign transactions for a chain
type ChainSigner interface {
	TryProcessOutbound(
		cctx *crosschaintypes.CrossChainTx,
		outboundManager *outboundprocessor.Processor,
		outboundID string,
		evmClient ChainClient,
		zetaBridge ZetaCoreBridger,
		height uint64,
	)
	SetZetaConnectorAddress(address ethcommon.Address)
	SetERC20CustodyAddress(address ethcommon.Address)
	GetZetaConnectorAddress() ethcommon.Address
	GetERC20CustodyAddress() ethcommon.Address
}

// ZetaCoreBridger is the interface to interact with ZetaCore
type ZetaCoreBridger interface {
	PostVoteInbound(gasLimit, retryGasLimit uint64, msg *crosschaintypes.MsgVoteInbound) (string, string, error)
	PostVoteOutbound(
		sendHash string,
		outboundHash string,
		outBlockHeight uint64,
		outTxGasUsed uint64,
		outTxEffectiveGasPrice *big.Int,
		outTxEffectiveGasLimit uint64,
		amount *big.Int,
		status chains.ReceiveStatus,
		chain chains.Chain,
		nonce uint64,
		coinType coin.CoinType,
	) (string, string, error)
	PostGasPrice(chain chains.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error)
	PostVoteBlockHeader(chainID int64, txhash []byte, height int64, header proofs.HeaderData) (string, error)
	GetBlockHeaderChainState(chainID int64) (lightclienttypes.QueryGetChainStateResponse, error)

	PostBlameData(blame *blame.Blame, chainID int64, index string) (string, error)
	AddTxHashToOutboundTracker(
		chainID int64,
		nonce uint64,
		txHash string,
		proof *proofs.Proof,
		blockHash string,
		txIndex int64,
	) (string, error)
	GetKeys() *keys.Keys
	GetZetaBlockHeight() (int64, error)
	GetLastBlockHeightByChain(chain chains.Chain) (*crosschaintypes.LastBlockHeight, error)
	ListPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error)
	ListPendingCctxWithinRatelimit() ([]*crosschaintypes.CrossChainTx, uint64, int64, string, bool, error)
	GetRateLimiterInput(window int64) (crosschaintypes.QueryRateLimiterInputResponse, error)
	GetPendingNoncesByChain(chainID int64) (observertypes.PendingNonces, error)
	GetCctxByNonce(chainID int64, nonce uint64) (*crosschaintypes.CrossChainTx, error)
	GetOutboundTracker(chain chains.Chain, nonce uint64) (*crosschaintypes.OutboundTracker, error)
	GetAllOutboundTrackerByChainbound(chainID int64, order Order) ([]crosschaintypes.OutboundTracker, error)
	GetCrosschainFlags() (observertypes.CrosschainFlags, error)
	GetRateLimiterFlags() (crosschaintypes.RateLimiterFlags, error)
	GetObserverList() ([]string, error)
	GetKeyGen() (*observertypes.Keygen, error)
	GetBtcTssAddress(chainID int64) (string, error)
	GetInboundTrackersForChain(chainID int64) ([]crosschaintypes.InboundTracker, error)
	GetLogger() *zerolog.Logger
	ZetaChain() chains.Chain
	Pause()
	Unpause()
	GetZetaHotKeyBalance() (sdkmath.Int, error)
}

// BTCRPCClient is the interface for BTC RPC client
type BTCRPCClient interface {
	GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error)
	CreateWallet(name string, opts ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error)
	GetNewAddress(account string) (btcutil.Address, error)
	GenerateToAddress(numBlocks int64, address btcutil.Address, maxTries *int64) ([]*chainhash.Hash, error)
	GetBalance(account string) (btcutil.Amount, error)
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	ListUnspent() ([]btcjson.ListUnspentResult, error)
	ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error)
	EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)
	GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error)
	GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error)
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

// EVMJSONRPCClient is the interface for EVM JSON RPC client
type EVMJSONRPCClient interface {
	EthGetBlockByNumber(number int, withTransactions bool) (*ethrpc.Block, error)
	EthGetTransactionByHash(hash string) (*ethrpc.Transaction, error)
}
