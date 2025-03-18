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
	IsRegnet() bool
	Healthcheck(ctx context.Context, tssAddress btcutil.Address) (time.Time, error)
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
	GetTotalMempoolParentsSizeNFees(
		ctx context.Context,
		childHash string,
		timeout time.Duration,
	) (int64, float64, int64, int64, error)

	GetRawTransactionResult(
		ctx context.Context,
		hash *hash.Hash,
		res *types.GetTransactionResult,
	) (types.TxRawResult, error)

	CreateRawTransaction(
		ctx context.Context,
		inputs []types.TransactionInput,
		amounts map[btcutil.Address]btcutil.Amount,
		lockTime *int64,
	) (*wire.MsgTx, error)

	SendRawTransaction(ctx context.Context, tx *wire.MsgTx, allowHighFees bool) (*hash.Hash, error)

	GetEstimatedFeeRate(ctx context.Context, confTarget int64) (int64, error)
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

	ListUnspent(ctx context.Context) ([]types.ListUnspentResult, error)
	ListUnspentMinMaxAddresses(
		ctx context.Context,
		minConf, maxConf int,
		addresses []btcutil.Address,
	) ([]types.ListUnspentResult, error)

	CreateWallet(ctx context.Context, name string, opts ...rpcclient.CreateWalletOpt) (*types.CreateWalletResult, error)
	GetNewAddress(ctx context.Context, account string) (btcutil.Address, error)
	ImportAddress(ctx context.Context, address string) error
	ImportPrivKeyRescan(ctx context.Context, privKeyWIF *btcutil.WIF, label string, rescan bool) error
	GetBalance(ctx context.Context, account string) (btcutil.Amount, error)
	GenerateToAddress(
		ctx context.Context,
		numBlocks int64,
		address btcutil.Address,
		maxTries *int64,
	) ([]*hash.Hash, error)

	SignRawTransactionWithWallet2(
		ctx context.Context,
		tx *wire.MsgTx,
		inputs []types.RawTxWitnessInput,
	) (*wire.MsgTx, bool, error)

	RawRequest(ctx context.Context, method string, params []json.RawMessage) (json.RawMessage, error)
}
