// Package observer implements the Bitcoin chain observer
package observer

import (
	"math/big"
	"sync/atomic"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

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
	btcClient interfaces.BTCRPCClient

	// pendingNonce is the outbound artificial pending nonce
	pendingNonce uint64

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

// NewObserver returns a new Bitcoin chain observer
func NewObserver(
	chain chains.Chain,
	btcClient interfaces.BTCRPCClient,
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
		return nil, errors.Wrapf(err, "unable to create base observer for chain %d", chain.ChainId)
	}

	// get the bitcoin network params
	netParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get BTC net params for chain %d", chain.ChainId)
	}

	// create bitcoin observer
	ob := &Observer{
		Observer:          *baseObserver,
		netParams:         netParams,
		btcClient:         btcClient,
		utxos:             []btcjson.ListUnspentResult{},
		tssOutboundHashes: make(map[string]bool),
		includedTxResults: make(map[string]*btcjson.GetTransactionResult),
		broadcastedTx:     make(map[string]string),
		logger: Logger{
			ObserverLogger: *baseObserver.Logger(),
			UTXOs:          baseObserver.Logger().Chain.With().Str("module", "utxos").Logger(),
		},
	}

	ob.nodeEnabled.Store(true)

	// load last scanned block
	if err = ob.LoadLastBlockScanned(); err != nil {
		return nil, errors.Wrap(err, "unable to load last scanned block")
	}

	// load broadcasted transactions
	if err = ob.LoadBroadcastedTxMap(); err != nil {
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

// SetPendingNonce sets the artificial pending nonce
func (ob *Observer) SetPendingNonce(nonce uint64) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.pendingNonce = nonce
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

// GetBlockByNumberCached gets cached block (and header) by block number
func (ob *Observer) GetBlockByNumberCached(blockNumber int64) (*BTCBlockNHeader, error) {
	if result, ok := ob.BlockCache().Get(blockNumber); ok {
		if block, ok := result.(*BTCBlockNHeader); ok {
			return block, nil
		}
		return nil, errors.New("cached value is not of type *BTCBlockNHeader")
	}

	// Get the block hash
	hash, err := ob.btcClient.GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	// Get the block header
	header, err := ob.btcClient.GetBlockHeader(hash)
	if err != nil {
		return nil, err
	}
	// Get the block with verbose transactions
	block, err := ob.btcClient.GetBlockVerboseTx(hash)
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

// GetLastStuckOutbound returns the last stuck outbound tx information
func (ob *Observer) GetLastStuckOutbound() *LastStuckOutbound {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.lastStuckTx
}

// SetLastStuckOutbound sets the information of last stuck outbound
func (ob *Observer) SetLastStuckOutbound(stuckTx *LastStuckOutbound) {
	lf := map[string]any{
		logs.FieldMethod: "SetLastStuckOutbound",
	}

	ob.Mu().Lock()
	defer ob.Mu().Unlock()

	if stuckTx != nil {
		lf[logs.FieldNonce] = stuckTx.Nonce
		lf[logs.FieldTx] = stuckTx.Tx.MsgTx().TxID()
		ob.logger.Outbound.Warn().
			Fields(lf).
			Msgf("Bitcoin outbound is stuck for %f minutes", stuckTx.StuckFor.Minutes())
	} else if ob.lastStuckTx != nil {
		lf[logs.FieldNonce] = ob.lastStuckTx.Nonce
		lf[logs.FieldTx] = ob.lastStuckTx.Tx.MsgTx().TxID()
		ob.logger.Outbound.Info().Fields(lf).Msgf("Bitcoin outbound is no longer stuck")
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

func (ob *Observer) isNodeEnabled() bool {
	return ob.nodeEnabled.Load()
}
