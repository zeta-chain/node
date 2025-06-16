package client

import (
	"context"
	"encoding/json"
	"time"

	types "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	hash "github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

// client represents interface version of Client.
// It's unexported on purpose ONLY for mock generation.
//
//go:generate mockery --name client --structname BitcoinClient --filename bitcoin_client.go --output ../../../testutils/mocks
type client interface {
	Ping(ctx context.Context) error
	Healthcheck(ctx context.Context) (time.Time, error)
	GetNetworkInfo(ctx context.Context) (*types.GetNetworkInfoResult, error)

	GetBlockCount(ctx context.Context) (int64, error)
	GetBlockHash(ctx context.Context, blockHeight int64) (*hash.Hash, error)
	GetBlockHeader(ctx context.Context, hash *hash.Hash) (*wire.BlockHeader, error)
	GetBlockVerbose(ctx context.Context, hash *hash.Hash) (*types.GetBlockVerboseTxResult, error)

	GetTransaction(ctx context.Context, hash *hash.Hash) (*types.GetTransactionResult, error)
	GetRawTransaction(ctx context.Context, hash *hash.Hash) (*btcutil.Tx, error)
	GetRawTransactionVerbose(ctx context.Context, hash *hash.Hash) (*types.TxRawResult, error)
	GetMempoolEntry(ctx context.Context, txHash string) (*types.GetMempoolEntryResult, error)
	GetRawMempool(ctx context.Context) ([]*hash.Hash, error)
	GetMempoolTxsAndFees(ctx context.Context, childHash string) (MempoolTxsAndFees, error)

	GetRawTransactionResult(
		ctx context.Context,
		hash *hash.Hash,
		res *types.GetTransactionResult,
	) (types.TxRawResult, error)

	SendRawTransaction(ctx context.Context, tx *wire.MsgTx, allowHighFees bool) (*hash.Hash, error)

	GetEstimatedFeeRate(ctx context.Context, confTarget int64) (uint64, error)

	IsTxStuckInMempool(
		ctx context.Context,
		txHash string,
		maxWaitBlocks int64,
	) (stuck bool, pendingFor time.Duration, err error)

	GetTransactionFeeAndRate(ctx context.Context, tx *types.TxRawResult) (int64, int64, error)
	EstimateSmartFee(
		ctx context.Context,
		confTarget int64,
		mode *types.EstimateSmartFeeMode,
	) (*types.EstimateSmartFeeResult, error)

	GetBlockVerboseByStr(ctx context.Context, blockHash string) (*types.GetBlockVerboseTxResult, error)
	GetBlockHeightByStr(ctx context.Context, blockHash string) (int64, error)
	GetTransactionByStr(ctx context.Context, hash string) (*hash.Hash, *types.GetTransactionResult, error)
	GetRawTransactionByStr(ctx context.Context, hash string) (*btcutil.Tx, error)

	GetTransactionInputSpender(ctx context.Context, txid string, vout uint32) (string, error)
	GetTransactionInitiator(ctx context.Context, txid string) (string, error)

	ListUnspent(ctx context.Context) ([]types.ListUnspentResult, error)
	ListUnspentMinMaxAddresses(
		ctx context.Context,
		minConf, maxConf int,
		addresses []btcutil.Address,
	) ([]types.ListUnspentResult, error)

	CreateWallet(ctx context.Context, name string, opts ...rpcclient.CreateWalletOpt) (*types.CreateWalletResult, error)
	GetNewAddress(ctx context.Context, account string) (btcutil.Address, error)
	ImportAddress(ctx context.Context, address string) error
	GetBalance(ctx context.Context, account string) (btcutil.Amount, error)
	GenerateToAddress(
		ctx context.Context,
		numBlocks int64,
		address btcutil.Address,
		maxTries *int64,
	) ([]*hash.Hash, error)

	RawRequest(ctx context.Context, method string, params []json.RawMessage) (json.RawMessage, error)
}
