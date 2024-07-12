// Package interfaces provides interfaces for clients and signers for the chain to interact with
package interfaces

import (
	"context"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	keyinterfaces "github.com/zeta-chain/zetacore/zetaclient/keys/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
)

type Order string

const (
	NoOrder    Order = ""
	Ascending  Order = "ASC"
	Descending Order = "DESC"
)

// ChainObserver is the interface for chain observer
type ChainObserver interface {
	Start(ctx context.Context)
	Stop()
	IsOutboundProcessed(
		ctx context.Context,
		cctx *crosschaintypes.CrossChainTx,
		logger zerolog.Logger,
	) (bool, bool, error)
	SetChainParams(observertypes.ChainParams)
	GetChainParams() observertypes.ChainParams
	GetTxID(nonce uint64) string
	WatchInboundTracker(ctx context.Context) error
}

// ChainSigner is the interface to sign transactions for a chain
type ChainSigner interface {
	TryProcessOutbound(
		ctx context.Context,
		cctx *crosschaintypes.CrossChainTx,
		outboundProc *outboundprocessor.Processor,
		outboundID string,
		observer ChainObserver,
		zetacoreClient ZetacoreClient,
		height uint64,
	)
	SetZetaConnectorAddress(address ethcommon.Address)
	SetERC20CustodyAddress(address ethcommon.Address)
	GetZetaConnectorAddress() ethcommon.Address
	GetERC20CustodyAddress() ethcommon.Address
}

// ZetacoreVoter represents voter interface.
type ZetacoreVoter interface {
	PostVoteBlockHeader(
		ctx context.Context,
		chainID int64,
		txhash []byte,
		height int64,
		header proofs.HeaderData,
	) (string, error)
	PostVoteGasPrice(
		ctx context.Context,
		chain chains.Chain,
		gasPrice uint64,
		supply string,
		blockNum uint64,
	) (string, error)
	PostVoteInbound(
		ctx context.Context,
		gasLimit, retryGasLimit uint64,
		msg *crosschaintypes.MsgVoteInbound,
	) (string, string, error)
	PostVoteOutbound(
		ctx context.Context,
		gasLimit, retryGasLimit uint64,
		msg *crosschaintypes.MsgVoteOutbound,
	) (string, string, error)
	PostVoteBlameData(ctx context.Context, blame *blame.Blame, chainID int64, index string) (string, error)
}

// ZetacoreClient is the client interface to interact with zetacore
//
//go:generate mockery --name ZetacoreClient --filename zetacore_client.go --case underscore --output ../../testutils/mocks
type ZetacoreClient interface {
	ZetacoreVoter

	Chain() chains.Chain
	GetLogger() *zerolog.Logger
	GetKeys() keyinterfaces.ObserverKeys

	GetKeyGen(ctx context.Context) (*observertypes.Keygen, error)

	GetBlockHeight(ctx context.Context) (int64, error)
	GetBlockHeaderChainState(ctx context.Context, chainID int64) (*lightclienttypes.ChainState, error)

	ListPendingCCTX(ctx context.Context, chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error)
	ListPendingCCTXWithinRateLimit(
		ctx context.Context,
	) (*crosschaintypes.QueryListPendingCctxWithinRateLimitResponse, error)

	GetRateLimiterInput(ctx context.Context, window int64) (*crosschaintypes.QueryRateLimiterInputResponse, error)
	GetPendingNoncesByChain(ctx context.Context, chainID int64) (observertypes.PendingNonces, error)

	GetCctxByNonce(ctx context.Context, chainID int64, nonce uint64) (*crosschaintypes.CrossChainTx, error)
	GetOutboundTracker(ctx context.Context, chain chains.Chain, nonce uint64) (*crosschaintypes.OutboundTracker, error)
	GetAllOutboundTrackerByChain(
		ctx context.Context,
		chainID int64,
		order Order,
	) ([]crosschaintypes.OutboundTracker, error)
	GetCrosschainFlags(ctx context.Context) (observertypes.CrosschainFlags, error)
	GetRateLimiterFlags(ctx context.Context) (crosschaintypes.RateLimiterFlags, error)
	GetObserverList(ctx context.Context) ([]string, error)
	GetBTCTSSAddress(ctx context.Context, chainID int64) (string, error)
	GetZetaHotKeyBalance(ctx context.Context) (sdkmath.Int, error)
	GetInboundTrackersForChain(ctx context.Context, chainID int64) ([]crosschaintypes.InboundTracker, error)

	// todo(revamp): refactor input to struct
	AddOutboundTracker(
		ctx context.Context,
		chainID int64,
		nonce uint64,
		txHash string,
		proof *proofs.Proof,
		blockHash string,
		txIndex int64,
	) (string, error)

	Stop()
	OnBeforeStop(callback func())
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
	TransactionSender(
		ctx context.Context,
		tx *ethtypes.Transaction,
		block ethcommon.Hash,
		index uint,
	) (ethcommon.Address, error)
}

// EVMJSONRPCClient is the interface for EVM JSON RPC client
type EVMJSONRPCClient interface {
	EthGetBlockByNumber(number int, withTransactions bool) (*ethrpc.Block, error)
	EthGetTransactionByHash(hash string) (*ethrpc.Transaction, error)
}

// TSSSigner is the interface for TSS signer
type TSSSigner interface {
	Pubkey() []byte

	// Sign signs the data
	// Note: it specifies optionalPubkey to use a different pubkey than the current pubkey set during keygen
	// TODO: check if optionalPubkey is needed
	// https://github.com/zeta-chain/node/issues/2085
	Sign(
		ctx context.Context,
		data []byte,
		height uint64,
		nonce uint64,
		chainID int64,
		optionalPubkey string,
	) ([65]byte, error)

	// SignBatch signs the data in batch
	SignBatch(ctx context.Context, digests [][]byte, height uint64, nonce uint64, chainID int64) ([][65]byte, error)

	EVMAddress() ethcommon.Address
	BTCAddress() string
	BTCAddressWitnessPubkeyHash() *btcutil.AddressWitnessPubKeyHash
	PubKeyCompressedBytes() []byte
}
