package clients

import (
	"context"
	"math/big"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
)

// EVMClient defines the interface for EVM chain interactions
type EVMClient interface {
	TransactionByHash(ctx context.Context, hash string) (*ethtypes.Transaction, bool, error)
	TransactionReceipt(ctx context.Context, hash string) (*ethtypes.Receipt, error)
	IsTxConfirmed(ctx context.Context, txHash string, confirmations uint64) (bool, error)
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
