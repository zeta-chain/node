// Package observer implements the Bitcoin chain observer
package observer

import (
	"context"
	"fmt"
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
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

type BitcoinClient interface {
	Healthcheck(context.Context) (time.Time, error)

	GetBlockCount(context.Context) (int64, error)
	GetBlockHash(_ context.Context, blockHeight int64) (*hash.Hash, error)
	GetBlockHeader(context.Context, *hash.Hash) (*wire.BlockHeader, error)
	GetBlockVerbose(context.Context, *hash.Hash) (*btcjson.GetBlockVerboseTxResult, error)

	GetRawTransaction(context.Context, *hash.Hash) (*btcutil.Tx, error)
	GetRawTransactionVerbose(context.Context, *hash.Hash) (*btcjson.TxRawResult, error)
	GetRawTransactionResult(context.Context,
		*hash.Hash,
		*btcjson.GetTransactionResult,
	) (btcjson.TxRawResult, error)
	GetMempoolEntry(_ context.Context, txHash string) (*btcjson.GetMempoolEntryResult, error)

	GetEstimatedFeeRate(_ context.Context, confTarget int64) (uint64, error)
	GetTransactionFeeAndRate(_ context.Context, tx *btcjson.TxRawResult) (int64, int64, error)

	IsTxStuckInMempool(_ context.Context,
		txHash string,
		maxWaitBlocks int64,
	) (stuck bool, pendingFor time.Duration, err error)

	EstimateSmartFee(_ context.Context,
		confTarget int64,
		_ *btcjson.EstimateSmartFeeMode,
	) (*btcjson.EstimateSmartFeeResult, error)

	ListUnspentMinMaxAddresses(_ context.Context,
		minConf int,
		maxConf int,
		_ []btcutil.Address,
	) ([]btcjson.ListUnspentResult, error)

	GetBlockHeightByStr(_ context.Context, blockHash string) (int64, error)
	GetTransactionByStr(_ context.Context, hash string) (*hash.Hash, *btcjson.GetTransactionResult, error)
	GetRawTransactionByStr(_ context.Context, hash string) (*btcutil.Tx, error)

	GetTransactionInputSpender(_ context.Context, txid string, vout uint32) (string, error)
	GetTransactionInitiator(_ context.Context, txid string) (string, error)
}

const (
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
	*base.Observer

	// netParams contains the Bitcoin network parameters
	netParams *chaincfg.Params

	// bitcoinClient is the Bitcoin RPC client that interacts with the Bitcoin node
	bitcoinClient BitcoinClient

	// pendingNonce is the outbound artificial pending nonce
	pendingNonce uint64

	// feeBumpWaitBlocks is the number of blocks to await before considering a tx stuck in mempool
	feeBumpWaitBlocks int64

	// lastStuckTx contains the last stuck outbound tx information
	// Note: nil if outbound is not stuck
	lastStuckTx *LastStuckOutbound

	// utxos contains the UTXOs owned by the TSS address
	utxos []btcjson.ListUnspentResult

	// tssOutboundHashes keeps track of outbound hashes sent from TSS address
	tssOutboundHashes map[string]bool

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

// New BTC Observer constructor.
func New(baseObserver *base.Observer,
	bitcoinClient BitcoinClient,
	chain chains.Chain,
) (*Observer, error) {
	// get the bitcoin network params
	netParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get BTC net params")
	}

	isRegnet := chains.IsBitcoinRegnet(chain.ChainId)

	feeBumpWaitBlocks := pendingTxFeeBumpWaitBlocks
	if isRegnet {
		feeBumpWaitBlocks = pendingTxFeeBumpWaitBlocksRegnet
	}

	// create bitcoin observer
	ob := &Observer{
		Observer:      baseObserver,
		netParams:     netParams,
		bitcoinClient: bitcoinClient,

		pendingNonce:      0,
		feeBumpWaitBlocks: int64(feeBumpWaitBlocks),
		lastStuckTx:       nil,
		utxos:             []btcjson.ListUnspentResult{},

		tssOutboundHashes: make(map[string]bool),
		includedTxResults: make(map[string]*btcjson.GetTransactionResult),
		broadcastedTx:     make(map[string]string),

		logger: Logger{
			ObserverLogger: *baseObserver.Logger(),
			UTXOs:          baseObserver.Logger().Chain.With().Str("module", "utxos").Logger(),
		},

		nodeEnabled: atomic.Bool{},
	}

	ob.nodeEnabled.Store(true)

	// load last scanned block
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ob.LoadLastBlockScanned(ctx); err != nil {
		return nil, errors.Wrap(err, "unable to load last scanned block")
	}

	// load broadcasted transactions
	if err = ob.loadBroadcastedTxMap(); err != nil {
		return nil, errors.Wrap(err, "unable to load broadcasted tx map")
	}

	return ob, nil
}

// GetPendingNonce returns the artificial pending nonce
// Note: pending nonce is accessed concurrently
func (ob *Observer) GetPendingNonce() uint64 {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.pendingNonce
}

func (ob *Observer) setPendingNonce(nonce uint64) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.pendingNonce = nonce
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
	hash, err := ob.bitcoinClient.GetBlockHash(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	// Get the block header
	header, err := ob.bitcoinClient.GetBlockHeader(ctx, hash)
	if err != nil {
		return nil, err
	}
	// Get the block with verbose transactions
	block, err := ob.bitcoinClient.GetBlockVerbose(ctx, hash)
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

// LastStuckOutbound returns the last stuck outbound tx information
func (ob *Observer) LastStuckOutbound() (tx *LastStuckOutbound, found bool) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()

	return ob.lastStuckTx, ob.lastStuckTx != nil
}

// setLastStuckOutbound sets the information of last stuck outbound
func (ob *Observer) setLastStuckOutbound(stuckTx *LastStuckOutbound) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()

	logger := ob.logger.Outbound

	if stuckTx != nil {
		logger.Warn().
			Uint64(logs.FieldNonce, stuckTx.Nonce).
			Str(logs.FieldTx, stuckTx.Tx.MsgTx().TxID()).
			Float64("duration_in_minutes", stuckTx.StuckFor.Minutes()).
			Msg("bitcoin outbound is stuck")
	} else if ob.lastStuckTx != nil {
		logger.Info().
			Uint64(logs.FieldNonce, ob.lastStuckTx.Nonce).
			Str(logs.FieldTx, ob.lastStuckTx.Tx.MsgTx().TxID()).
			Msg("bitcoin outbound is no longer stuck")
	}
	ob.lastStuckTx = stuckTx
}

// IsTSSTransaction checks if a given transaction was sent by TSS itself.
// An unconfirmed transaction is safe to spend only if it was sent by TSS self.
func (ob *Observer) IsTSSTransaction(txid string) bool {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()

	_, found := ob.tssOutboundHashes[txid]
	return found
}

// GetBroadcastedTx gets successfully broadcasted transaction by nonce
func (ob *Observer) GetBroadcastedTx(nonce uint64) (string, bool) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()

	outboundID := ob.OutboundID(nonce)
	txHash, found := ob.broadcastedTx[outboundID]
	return txHash, found
}

// CheckRPCStatus checks the RPC status of the Bitcoin chain
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	if !ob.isNodeEnabled() {
		return nil
	}

	blockTime, err := ob.bitcoinClient.Healthcheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}

func (ob *Observer) isNodeEnabled() bool {
	return ob.nodeEnabled.Load()
}

// updateLastBlock is a helper function to update the last block number.
// Note: keep last block up-to-date helps to avoid inaccurate confirmation.
func (ob *Observer) updateLastBlock(ctx context.Context) error {
	blockNumber, err := ob.bitcoinClient.GetBlockCount(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting block number")
	}
	if blockNumber < 0 {
		return fmt.Errorf("block number is negative: %d", blockNumber)
	}

	// 0 will be returned if the node is not synced
	if blockNumber == 0 {
		ob.nodeEnabled.Store(false)
		ob.Logger().Chain.Debug().Err(err).Msg("bitcoin node is not enabled")
		return nil
	}
	ob.nodeEnabled.Store(true)

	// #nosec G115 checked positive
	if uint64(blockNumber) < ob.LastBlock() {
		return fmt.Errorf("block number should not decrease: current %d last %d", blockNumber, ob.LastBlock())
	}
	ob.WithLastBlock(uint64(blockNumber))

	return nil
}
