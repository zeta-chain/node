// Package observer implements the Bitcoin chain observer
package observer

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	hash "github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

type RPC interface {
	Healthcheck(ctx context.Context, tssAddress btcutil.Address) (time.Time, error)

	GetBlockCount(ctx context.Context) (int64, error)
	GetBlockHash(ctx context.Context, blockHeight int64) (*hash.Hash, error)
	GetBlockHeader(ctx context.Context, hash *hash.Hash) (*wire.BlockHeader, error)
	GetBlockVerbose(ctx context.Context, hash *hash.Hash) (*btcjson.GetBlockVerboseTxResult, error)

	GetTransaction(ctx context.Context, hash *hash.Hash) (*btcjson.GetTransactionResult, error)
	GetRawTransaction(ctx context.Context, hash *hash.Hash) (*btcutil.Tx, error)
	GetRawTransactionVerbose(ctx context.Context, hash *hash.Hash) (*btcjson.TxRawResult, error)
	GetRawTransactionResult(
		ctx context.Context,
		hash *hash.Hash,
		res *btcjson.GetTransactionResult,
	) (btcjson.TxRawResult, error)

	GetTransactionFeeAndRate(ctx context.Context, tx *btcjson.TxRawResult) (int64, int64, error)

	EstimateSmartFee(
		ctx context.Context,
		confTarget int64,
		mode *btcjson.EstimateSmartFeeMode,
	) (*btcjson.EstimateSmartFeeResult, error)

	ListUnspentMinMaxAddresses(
		ctx context.Context,
		minConf, maxConf int,
		addresses []btcutil.Address,
	) ([]btcjson.ListUnspentResult, error)

	GetBlockVerboseByStr(ctx context.Context, blockHash string) (*btcjson.GetBlockVerboseTxResult, error)
	GetBlockHeightByStr(ctx context.Context, blockHash string) (int64, error)
	GetTransactionByStr(ctx context.Context, hash string) (*hash.Hash, *btcjson.GetTransactionResult, error)
	GetRawTransactionByStr(ctx context.Context, hash string) (*btcutil.Tx, error)
}

const (
	// btcBlocksPerDay represents Bitcoin blocks per days for LRU block cache size
	btcBlocksPerDay = 144

	// RegnetStartBlock is the hardcoded start block for regnet
	RegnetStartBlock = 100

	// BigValueSats contains the threshold to determine a big value in Bitcoin represents 2 BTC
	BigValueSats = 200000000

	// BigValueConfirmationCount represents the number of confirmation necessary for bigger values: 6 confirmations
	BigValueConfirmationCount = 6
)

// Logger contains list of loggers used by Bitcoin chain observer
type Logger struct {
	// base.Logger contains a list of base observer loggers
	base.ObserverLogger

	// UTXOs is the logger for UTXOs management
	UTXOs zerolog.Logger
}

// BTCBlockNHeader contains bitcoin block and the header
type BTCBlockNHeader struct {
	Header *wire.BlockHeader
	Block  *btcjson.GetBlockVerboseTxResult
}

// Observer is the Bitcoin chain observer
type Observer struct {
	// base.Observer implements the base chain observer
	base.Observer

	// netParams contains the Bitcoin network parameters
	netParams *chaincfg.Params

	// btcClient is the Bitcoin RPC client that interacts with the Bitcoin node
	rpc RPC

	// pendingNonce is the outbound artificial pending nonce
	pendingNonce uint64

	// utxos contains the UTXOs owned by the TSS address
	utxos []btcjson.ListUnspentResult

	// includedTxHashes indexes included tx with tx hash
	includedTxHashes map[string]bool

	// includedTxResults indexes tx results with the outbound tx identifier
	includedTxResults map[string]*btcjson.GetTransactionResult

	// broadcastedTx indexes the outbound hash with the outbound tx identifier
	broadcastedTx map[string]string

	// nodeEnabled indicates whether BTC node is enabled (might be disabled during certain E2E tests)
	// We assume it's true by default. The flag is updated on each ObserveInbound call.
	nodeEnabled atomic.Bool

	// logger contains the loggers used by the bitcoin observer
	logger Logger
}

// NewObserver returns a new Bitcoin chain observer
func NewObserver(
	chain chains.Chain,
	rpc RPC,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	database *db.DB,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*Observer, error) {
	// create base observer
	baseObserver, err := base.NewObserver(
		chain,
		chainParams,
		zetacoreClient,
		tss,
		btcBlocksPerDay,
		ts,
		database,
		logger,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create base observer")
	}

	// get the bitcoin network params
	netParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get BTC net params")
	}

	// create bitcoin observer
	ob := &Observer{
		Observer:          *baseObserver,
		netParams:         netParams,
		rpc:               rpc,
		pendingNonce:      0,
		utxos:             []btcjson.ListUnspentResult{},
		includedTxHashes:  make(map[string]bool),
		includedTxResults: make(map[string]*btcjson.GetTransactionResult),
		broadcastedTx:     make(map[string]string),
		logger: Logger{
			ObserverLogger: *baseObserver.Logger(),
			UTXOs:          baseObserver.Logger().Chain.With().Str("module", "utxos").Logger(),
		},
	}

	ob.nodeEnabled.Store(true)

	// load last scanned block
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ob.LoadLastBlockScanned(ctx); err != nil {
		return nil, errors.Wrap(err, "unable to load last scanned block")
	}

	// load broadcasted transactions
	if err = ob.LoadBroadcastedTxMap(); err != nil {
		return nil, errors.Wrap(err, "unable to load broadcasted tx map")
	}

	return ob, nil
}

func (ob *Observer) isNodeEnabled() bool {
	return ob.nodeEnabled.Load()
}

// GetPendingNonce returns the artificial pending nonce
// Note: pending nonce is accessed concurrently
func (ob *Observer) GetPendingNonce() uint64 {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.pendingNonce
}

// ConfirmationsThreshold returns number of required Bitcoin confirmations depending on sent BTC amount.
func (ob *Observer) ConfirmationsThreshold(amount *big.Int) int64 {
	if amount.Cmp(big.NewInt(BigValueSats)) >= 0 {
		return BigValueConfirmationCount
	}
	if BigValueConfirmationCount < ob.ChainParams().ConfirmationCount {
		return BigValueConfirmationCount
	}

	// #nosec G115 always in range
	return int64(ob.ChainParams().ConfirmationCount)
}

// PostGasPrice posts gas price to zetacore
// TODO(revamp): move to gas price file
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	var (
		err              error
		feeRateEstimated uint64
	)

	// special handle regnet and testnet gas rate
	// regnet:  RPC 'EstimateSmartFee' is not available
	// testnet: RPC 'EstimateSmartFee' returns unreasonable high gas rate
	if ob.Chain().NetworkType != chains.NetworkType_mainnet {
		feeRateEstimated, err = ob.specialHandleFeeRate(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to execute specialHandleFeeRate")
		}
	} else {
		// EstimateSmartFee returns the fees per kilobyte (BTC/kb) targeting given block confirmation
		feeResult, err := ob.rpc.EstimateSmartFee(ctx, 1, &btcjson.EstimateModeEconomical)
		if err != nil {
			return errors.Wrap(err, "unable to estimate smart fee")
		}
		if feeResult.Errors != nil || feeResult.FeeRate == nil {
			return fmt.Errorf("error getting gas price: %s", feeResult.Errors)
		}
		if *feeResult.FeeRate > math.MaxInt64 {
			return fmt.Errorf("gas price is too large: %f", *feeResult.FeeRate)
		}
		feeRateEstimated = common.FeeRateToSatPerByte(*feeResult.FeeRate).Uint64()
	}

	// query the current block number
	blockNumber, err := ob.rpc.GetBlockCount(ctx)
	if err != nil {
		return errors.Wrap(err, "GetBlockCount error")
	}

	// UTXO has no concept of priority fee (like eth)
	const priorityFee = 0

	// #nosec G115 always positive
	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), feeRateEstimated, priorityFee, uint64(blockNumber))
	if err != nil {
		return errors.Wrap(err, "PostVoteGasPrice error")
	}

	return nil
}

// FetchUTXOs fetches TSS-owned UTXOs from the Bitcoin node
// TODO(revamp): move to UTXO file
func (ob *Observer) FetchUTXOs(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			ob.logger.UTXOs.Error().Msgf("BTC FetchUTXOs: caught panic error: %v", err)
		}
	}()

	// noop
	if !ob.isNodeEnabled() {
		return nil
	}

	// This is useful when a zetaclient's pending nonce lagged behind for whatever reason.
	ob.refreshPendingNonce(ctx)

	// get the current block height.
	bh, err := ob.rpc.GetBlockCount(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get block height")
	}

	maxConfirmations := int(bh)

	// List all unspent UTXOs (160ms)
	tssAddr, err := ob.TSS().PubKey().AddressBTC(ob.Chain().ChainId)
	if err != nil {
		return errors.Wrap(err, "unable to get tss address")
	}

	utxos, err := ob.rpc.ListUnspentMinMaxAddresses(ctx, 0, maxConfirmations, []btcutil.Address{tssAddr})
	if err != nil {
		return errors.Wrap(err, "unable to list unspent utxo")
	}

	// rigid sort to make utxo list deterministic
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Amount == utxos[j].Amount {
			if utxos[i].TxID == utxos[j].TxID {
				return utxos[i].Vout < utxos[j].Vout
			}
			return utxos[i].TxID < utxos[j].TxID
		}
		return utxos[i].Amount < utxos[j].Amount
	})

	// filter UTXOs good to spend for next TSS transaction
	utxosFiltered := make([]btcjson.ListUnspentResult, 0)
	for _, utxo := range utxos {
		// UTXOs big enough to cover the cost of spending themselves
		if utxo.Amount < common.DefaultDepositorFee {
			continue
		}
		// we don't want to spend other people's unconfirmed UTXOs as they may not be safe to spend
		if utxo.Confirmations == 0 {
			if !ob.isTssTransaction(utxo.TxID) {
				continue
			}
		}
		utxosFiltered = append(utxosFiltered, utxo)
	}

	ob.Mu().Lock()
	ob.TelemetryServer().SetNumberOfUTXOs(len(utxosFiltered))
	ob.utxos = utxosFiltered
	ob.Mu().Unlock()
	return nil
}

// SaveBroadcastedTx saves successfully broadcasted transaction
// TODO(revamp): move to db file
func (ob *Observer) SaveBroadcastedTx(txHash string, nonce uint64) {
	outboundID := ob.OutboundID(nonce)
	ob.Mu().Lock()
	ob.broadcastedTx[outboundID] = txHash
	ob.Mu().Unlock()

	broadcastEntry := clienttypes.ToOutboundHashSQLType(txHash, outboundID)
	if err := ob.DB().Client().Save(&broadcastEntry).Error; err != nil {
		ob.logger.Outbound.Error().
			Err(err).
			Msgf("SaveBroadcastedTx: error saving broadcasted txHash %s for outbound %s", txHash, outboundID)
	}
	ob.logger.Outbound.Info().Msgf("SaveBroadcastedTx: saved broadcasted txHash %s for outbound %s", txHash, outboundID)
}

// GetBlockByNumberCached gets cached block (and header) by block number
func (ob *Observer) GetBlockByNumberCached(ctx context.Context, blockNumber int64) (*BTCBlockNHeader, error) {
	if result, ok := ob.BlockCache().Get(blockNumber); ok {
		if block, ok := result.(*BTCBlockNHeader); ok {
			return block, nil
		}
		return nil, errors.New("cached value is not of type *BTCBlockNHeader")
	}

	// Get the block hash
	hash, err := ob.rpc.GetBlockHash(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	// Get the block header
	header, err := ob.rpc.GetBlockHeader(ctx, hash)
	if err != nil {
		return nil, err
	}
	// Get the block with verbose transactions
	block, err := ob.rpc.GetBlockVerbose(ctx, hash)
	if err != nil {
		return nil, err
	}
	blockNheader := &BTCBlockNHeader{
		Header: header,
		Block:  block,
	}
	ob.BlockCache().Add(blockNumber, blockNheader)
	ob.BlockCache().Add(hash, blockNheader)
	return blockNheader, nil
}

// LoadLastBlockScanned loads the last scanned block from the database
func (ob *Observer) LoadLastBlockScanned(ctx context.Context) error {
	err := ob.Observer.LoadLastBlockScanned(ob.Logger().Chain)
	if err != nil {
		return errors.Wrapf(err, "error LoadLastBlockScanned for chain %d", ob.Chain().ChainId)
	}

	// observer will scan from the last block when 'lastBlockScanned == 0', this happens when:
	// 1. environment variable is set explicitly to "latest"
	// 2. environment variable is empty and last scanned block is not found in DB
	if ob.LastBlockScanned() == 0 {
		blockNumber, err := ob.rpc.GetBlockCount(ctx)
		if err != nil {
			return errors.Wrapf(err, "error GetBlockCount for chain %d", ob.Chain().ChainId)
		}
		// #nosec G115 always positive
		ob.WithLastBlockScanned(uint64(blockNumber))
	}

	// bitcoin regtest starts from hardcoded block 100
	if chains.IsBitcoinRegnet(ob.Chain().ChainId) {
		ob.WithLastBlockScanned(RegnetStartBlock)
	}
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from block %d", ob.Chain().ChainId, ob.LastBlockScanned())

	return nil
}

// LoadBroadcastedTxMap loads broadcasted transactions from the database
func (ob *Observer) LoadBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutboundHashSQLType
	if err := ob.DB().Client().Find(&broadcastedTransactions).Error; err != nil {
		ob.logger.Chain.Error().Err(err).Msgf("error iterating over db for chain %d", ob.Chain().ChainId)
		return err
	}
	for _, entry := range broadcastedTransactions {
		ob.broadcastedTx[entry.Key] = entry.Hash
	}
	return nil
}

// specialHandleFeeRate handles the fee rate for regnet and testnet
func (ob *Observer) specialHandleFeeRate(ctx context.Context) (uint64, error) {
	switch ob.Chain().NetworkType {
	case chains.NetworkType_privnet:
		// hardcode gas price for regnet
		return 1, nil
	case chains.NetworkType_testnet:
		feeRateEstimated, err := common.GetRecentFeeRate(ctx, ob.rpc, ob.netParams)
		if err != nil {
			return 0, errors.Wrapf(err, "error GetRecentFeeRate")
		}
		return feeRateEstimated, nil
	default:
		return 0, fmt.Errorf(" unsupported bitcoin network type %d", ob.Chain().NetworkType)
	}
}

// isTssTransaction checks if a given transaction was sent by TSS itself.
// An unconfirmed transaction is safe to spend only if it was sent by TSS and verified by ourselves.
func (ob *Observer) isTssTransaction(txid string) bool {
	_, found := ob.includedTxHashes[txid]
	return found
}
