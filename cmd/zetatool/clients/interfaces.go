package clients

import (
	"context"
	"math/big"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// EVMClient defines the interface for EVM chain interactions
type EVMClient interface {
	TransactionByHash(ctx context.Context, hash string) (*ethtypes.Transaction, bool, error)
	TransactionReceipt(ctx context.Context, hash string) (*ethtypes.Receipt, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

// BitcoinClient defines the interface for Bitcoin chain interactions
type BitcoinClient interface {
	Ping(ctx context.Context) error
	GetRawTransactionVerbose(ctx context.Context, txHash *chainhash.Hash) (*btcjson.TxRawResult, error)
	GetBlockVerbose(ctx context.Context, blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)
}

// SolanaClient defines the interface for Solana chain interactions
type SolanaClient interface {
	GetTransaction(ctx context.Context, signature solana.Signature) (*solrpc.GetTransactionResult, error)
}

// ZetacoreReader defines the read-only interface for querying ZetaCore.
// This interface contains only the methods that zetatool needs for CCTX tracking.
type ZetacoreReader interface {
	// CCTX queries
	GetCctxByHash(ctx context.Context, hash string) (*crosschaintypes.CrossChainTx, error)
	InboundHashToCctxData(ctx context.Context, hash string) (*crosschaintypes.QueryInboundHashToCctxDataResponse, error)

	// Tracker queries
	GetOutboundTracker(ctx context.Context, chainID int64, nonce uint64) (*crosschaintypes.OutboundTracker, error)

	// Chain params
	GetChainParamsForChainID(ctx context.Context, chainID int64) (*observertypes.ChainParams, error)

	// TSS queries
	GetTssAddress(ctx context.Context, btcChainID int64) (*observertypes.QueryGetTssAddressResponse, error)
	GetEVMTSSAddress(ctx context.Context) (string, error)
	GetBTCTSSAddress(ctx context.Context, chainID int64) (string, error)

	// Ballot queries
	GetBallotByID(ctx context.Context, id string) (*observertypes.QueryBallotByIdentifierResponse, error)
}
