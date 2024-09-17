package base

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// EnvVarLatestBlock is the environment variable that forces the observer to scan from the latest block
	EnvVarLatestBlock = "latest"

	// DefaultBlockCacheSize is the default number of blocks that the observer will keep in cache for performance (without RPC calls)
	// Cached blocks can be used to get block information and verify transactions
	DefaultBlockCacheSize = 1000

	// DefaultHeaderCacheSize is the default number of headers that the observer will keep in cache for performance (without RPC calls)
	// Cached headers can be used to get header information
	DefaultHeaderCacheSize = 1000
)

// Observer is the base structure for chain observers, grouping the common logic for each chain observer client.
// The common logic includes: chain, chainParams, contexts, zetacore client, tss, lastBlock, db, metrics, loggers etc.
type Observer struct {
	// chain contains static information about the observed chain
	chain chains.Chain

	// chainParams contains the dynamic chain parameters of the observed chain
	chainParams observertypes.ChainParams

	// zetacoreClient is the client to interact with ZetaChain
	zetacoreClient interfaces.ZetacoreClient

	// tss is the TSS signer
	tss interfaces.TSSSigner

	// lastBlock is the last block height of the observed chain
	lastBlock uint64

	// lastBlockScanned is the last block height scanned by the observer
	lastBlockScanned uint64

	// lastTxScanned is the last transaction hash scanned by the observer
	lastTxScanned string

	// rpcAlertLatency is the threshold of RPC latency to trigger an alert
	rpcAlertLatency time.Duration

	// blockCache is the cache for blocks
	blockCache *lru.Cache

	// headerCache is the cache for headers
	headerCache *lru.Cache

	// db is the database to persist data
	db *db.DB

	// ts is the telemetry server for metrics
	ts *metrics.TelemetryServer

	// logger contains the loggers used by observer
	logger ObserverLogger

	// mu protects fields from concurrent access
	// Note: base observer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu      *sync.Mutex
	started bool

	// stop is the channel to signal the observer to stop
	stop chan struct{}
}

// NewObserver creates a new base observer.
func NewObserver(
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	blockCacheSize int,
	headerCacheSize int,
	rpcAlertLatency int64,
	ts *metrics.TelemetryServer,
	database *db.DB,
	logger Logger,
) (*Observer, error) {
	ob := Observer{
		chain:            chain,
		chainParams:      chainParams,
		zetacoreClient:   zetacoreClient,
		tss:              tss,
		lastBlock:        0,
		lastBlockScanned: 0,
		lastTxScanned:    "",
		rpcAlertLatency:  time.Duration(rpcAlertLatency) * time.Second,
		ts:               ts,
		db:               database,
		mu:               &sync.Mutex{},
		stop:             make(chan struct{}),
	}

	// setup loggers
	ob.WithLogger(logger)

	// create block cache
	var err error
	ob.blockCache, err = lru.New(blockCacheSize)
	if err != nil {
		return nil, errors.Wrap(err, "error creating block cache")
	}

	// create header cache
	ob.headerCache, err = lru.New(headerCacheSize)
	if err != nil {
		return nil, errors.Wrap(err, "error creating header cache")
	}

	return &ob, nil
}

// Start starts the observer. Returns true if the observer was already started (noop).
func (ob *Observer) Start() bool {
	ob.mu.Lock()
	defer ob.Mu().Unlock()

	// noop
	if ob.started {
		return true
	}

	ob.started = true

	return false
}

// Stop notifies all goroutines to stop and closes the database.
func (ob *Observer) Stop() {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if !ob.started {
		ob.logger.Chain.Info().Msgf("Observer already stopped for chain %d", ob.Chain().ChainId)
		return
	}

	ob.logger.Chain.Info().Msgf("Stopping observer for chain %d", ob.Chain().ChainId)

	close(ob.stop)
	ob.started = false

	// close database
	if err := ob.db.Close(); err != nil {
		ob.Logger().Chain.Error().Err(err).Msgf("unable to close db for chain %d", ob.Chain().ChainId)
	}

	ob.Logger().Chain.Info().Msgf("observer stopped for chain %d", ob.Chain().ChainId)
}

// Chain returns the chain for the observer.
func (ob *Observer) Chain() chains.Chain {
	return ob.chain
}

// WithChain attaches a new chain to the observer.
func (ob *Observer) WithChain(chain chains.Chain) *Observer {
	ob.chain = chain
	return ob
}

// ChainParams returns the chain params for the observer.
func (ob *Observer) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}

// WithChainParams attaches a new chain params to the observer.
func (ob *Observer) WithChainParams(params observertypes.ChainParams) *Observer {
	ob.chainParams = params
	return ob
}

// ZetacoreClient returns the zetacore client for the observer.
func (ob *Observer) ZetacoreClient() interfaces.ZetacoreClient {
	return ob.zetacoreClient
}

// WithZetacoreClient attaches a new zetacore client to the observer.
func (ob *Observer) WithZetacoreClient(client interfaces.ZetacoreClient) *Observer {
	ob.zetacoreClient = client
	return ob
}

// Tss returns the tss signer for the observer.
func (ob *Observer) TSS() interfaces.TSSSigner {
	return ob.tss
}

// WithTSS attaches a new tss signer to the observer.
func (ob *Observer) WithTSS(tss interfaces.TSSSigner) *Observer {
	ob.tss = tss
	return ob
}

// LastBlock get external last block height.
func (ob *Observer) LastBlock() uint64 {
	return atomic.LoadUint64(&ob.lastBlock)
}

// WithLastBlock set external last block height.
func (ob *Observer) WithLastBlock(lastBlock uint64) *Observer {
	atomic.StoreUint64(&ob.lastBlock, lastBlock)
	return ob
}

// LastBlockScanned get last block scanned (not necessarily caught up with the chain; could be slow/paused).
func (ob *Observer) LastBlockScanned() uint64 {
	height := atomic.LoadUint64(&ob.lastBlockScanned)
	return height
}

// WithLastBlockScanned set last block scanned (not necessarily caught up with the chain; could be slow/paused).
func (ob *Observer) WithLastBlockScanned(blockNumber uint64) *Observer {
	atomic.StoreUint64(&ob.lastBlockScanned, blockNumber)
	metrics.LastScannedBlockNumber.WithLabelValues(ob.chain.Name).Set(float64(blockNumber))
	return ob
}

// LastTxScanned get last transaction scanned.
func (ob *Observer) LastTxScanned() string {
	return ob.lastTxScanned
}

// WithLastTxScanned set last transaction scanned.
func (ob *Observer) WithLastTxScanned(txHash string) *Observer {
	ob.lastTxScanned = txHash
	return ob
}

// BlockCache returns the block cache for the observer.
func (ob *Observer) BlockCache() *lru.Cache {
	return ob.blockCache
}

// WithBlockCache attaches a new block cache to the observer.
func (ob *Observer) WithBlockCache(cache *lru.Cache) *Observer {
	ob.blockCache = cache
	return ob
}

// HeaderCache returns the header cache for the observer.
func (ob *Observer) HeaderCache() *lru.Cache {
	return ob.headerCache
}

// WithHeaderCache attaches a new header cache to the observer.
func (ob *Observer) WithHeaderCache(cache *lru.Cache) *Observer {
	ob.headerCache = cache
	return ob
}

// OutboundID returns a unique identifier for the outbound transaction.
// The identifier is now used as the key for maps that store outbound related data (e.g. transaction, receipt, etc).
func (ob *Observer) OutboundID(nonce uint64) string {
	// all chains uses EVM address as part of the key except bitcoin
	tssAddress := ob.tss.EVMAddress().String()
	if ob.chain.Consensus == chains.Consensus_bitcoin {
		tssAddress = ob.tss.BTCAddress()
	}
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddress, nonce)
}

// DB returns the database for the observer.
func (ob *Observer) DB() *db.DB {
	return ob.db
}

// WithTelemetryServer attaches a new telemetry server to the observer.
func (ob *Observer) WithTelemetryServer(ts *metrics.TelemetryServer) *Observer {
	ob.ts = ts
	return ob
}

// TelemetryServer returns the telemetry server for the observer.
func (ob *Observer) TelemetryServer() *metrics.TelemetryServer {
	return ob.ts
}

// Logger returns the logger for the observer.
func (ob *Observer) Logger() *ObserverLogger {
	return &ob.logger
}

// WithLogger attaches a new logger to the observer.
func (ob *Observer) WithLogger(logger Logger) *Observer {
	chainLogger := logger.Std.With().Int64(logs.FieldChain, ob.chain.ChainId).Logger()
	ob.logger = ObserverLogger{
		Chain:      chainLogger,
		Inbound:    chainLogger.With().Str(logs.FieldModule, logs.ModNameInbound).Logger(),
		Outbound:   chainLogger.With().Str(logs.FieldModule, logs.ModNameOutbound).Logger(),
		GasPrice:   chainLogger.With().Str(logs.FieldModule, logs.ModNameGasPrice).Logger(),
		Headers:    chainLogger.With().Str(logs.FieldModule, logs.ModNameHeaders).Logger(),
		Compliance: logger.Compliance,
	}
	return ob
}

// Mu returns the mutex for the observer.
func (ob *Observer) Mu() *sync.Mutex {
	return ob.mu
}

// StopChannel returns the stop channel for the observer.
func (ob *Observer) StopChannel() chan struct{} {
	return ob.stop
}

// LoadLastBlockScanned loads last scanned block from environment variable or from database.
// The last scanned block is the height from which the observer should continue scanning.
func (ob *Observer) LoadLastBlockScanned(logger zerolog.Logger) error {
	// get environment variable
	envvar := EnvVarLatestBlockByChain(ob.chain)
	scanFromBlock := os.Getenv(envvar)

	// load from environment variable if set
	if scanFromBlock != "" {
		logger.Info().
			Msgf("LoadLastBlockScanned: envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == EnvVarLatestBlock {
			return nil
		}
		blockNumber, err := strconv.ParseUint(scanFromBlock, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to parse block number from ENV %s=%s", envvar, scanFromBlock)
		}
		ob.WithLastBlockScanned(blockNumber)
		return nil
	}

	// load from DB otherwise. If not found, start from latest block
	blockNumber, err := ob.ReadLastBlockScannedFromDB()
	if err != nil {
		logger.Info().Msgf("LoadLastBlockScanned: last scanned block not found in db for chain %d", ob.chain.ChainId)
		return nil
	}
	ob.WithLastBlockScanned(blockNumber)

	return nil
}

// SaveLastBlockScanned saves the last scanned block to memory and database.
func (ob *Observer) SaveLastBlockScanned(blockNumber uint64) error {
	ob.WithLastBlockScanned(blockNumber)
	return ob.WriteLastBlockScannedToDB(blockNumber)
}

// WriteLastBlockScannedToDB saves the last scanned block to the database.
func (ob *Observer) WriteLastBlockScannedToDB(lastScannedBlock uint64) error {
	return ob.db.Client().Save(clienttypes.ToLastBlockSQLType(lastScannedBlock)).Error
}

// ReadLastBlockScannedFromDB reads the last scanned block from the database.
func (ob *Observer) ReadLastBlockScannedFromDB() (uint64, error) {
	var lastBlock clienttypes.LastBlockSQLType
	if err := ob.db.Client().First(&lastBlock, clienttypes.LastBlockNumID).Error; err != nil {
		// record not found
		return 0, err
	}
	return lastBlock.Num, nil
}

// LoadLastTxScanned loads last scanned tx from environment variable or from database.
// The last scanned tx is the tx hash from which the observer should continue scanning.
func (ob *Observer) LoadLastTxScanned() {
	// get environment variable
	envvar := EnvVarLatestTxByChain(ob.chain)
	scanFromTx := os.Getenv(envvar)

	// load from environment variable if set
	if scanFromTx != "" {
		ob.logger.Chain.Info().Msgf("LoadLastTxScanned: envvar %s is set; scan from  tx %s", envvar, scanFromTx)
		ob.WithLastTxScanned(scanFromTx)
		return
	}

	// load from DB otherwise.
	txHash, err := ob.ReadLastTxScannedFromDB()
	if err != nil {
		// If not found, let the concrete chain observer decide where to start
		ob.logger.Chain.Info().Msgf("LoadLastTxScanned: last scanned tx not found in db for chain %d", ob.chain.ChainId)
		return
	}
	ob.WithLastTxScanned(txHash)
}

// SaveLastTxScanned saves the last scanned tx hash to memory and database.
func (ob *Observer) SaveLastTxScanned(txHash string, slot uint64) error {
	// save last scanned tx to memory
	ob.WithLastTxScanned(txHash)

	// update last_scanned_block_number metrics
	ob.WithLastBlockScanned(slot)

	return ob.WriteLastTxScannedToDB(txHash)
}

// WriteLastTxScannedToDB saves the last scanned tx hash to the database.
func (ob *Observer) WriteLastTxScannedToDB(txHash string) error {
	return ob.db.Client().Save(clienttypes.ToLastTxHashSQLType(txHash)).Error
}

// ReadLastTxScannedFromDB reads the last scanned tx hash from the database.
func (ob *Observer) ReadLastTxScannedFromDB() (string, error) {
	var lastTx clienttypes.LastTransactionSQLType
	if err := ob.db.Client().First(&lastTx, clienttypes.LastTxHashID).Error; err != nil {
		// record not found
		return "", err
	}
	return lastTx.Hash, nil
}

// PostVoteInbound posts a vote for the given vote message
func (ob *Observer) PostVoteInbound(
	ctx context.Context,
	msg *crosschaintypes.MsgVoteInbound,
	retryGasLimit uint64,
) (string, error) {
	txHash := msg.InboundHash
	coinType := msg.CoinType
	chainID := ob.Chain().ChainId
	zetaHash, ballot, err := ob.ZetacoreClient().
		PostVoteInbound(ctx, zetacore.PostVoteInboundGasLimit, retryGasLimit, msg)
	if err != nil {
		ob.logger.Inbound.Err(err).
			Msgf("inbound detected: error posting vote for chain %d token %s inbound %s", chainID, coinType, txHash)
		return "", err
	} else if zetaHash != "" {
		ob.logger.Inbound.Info().Msgf("inbound detected: chain %d token %s inbound %s vote %s ballot %s", chainID, coinType, txHash, zetaHash, ballot)
	} else {
		ob.logger.Inbound.Info().Msgf("inbound detected: chain %d token %s inbound %s already voted on ballot %s", chainID, coinType, txHash, ballot)
	}

	return ballot, err
}

// AlertOnRPCLatency prints an alert if the RPC latency exceeds the threshold.
// Returns true if the RPC latency is too high.
func (ob *Observer) AlertOnRPCLatency(latestBlockTime time.Time, defaultAlertLatency time.Duration) bool {
	// use configured alert latency if set
	alertLatency := defaultAlertLatency
	if ob.rpcAlertLatency > 0 {
		alertLatency = ob.rpcAlertLatency
	}

	// latest block should not be too old
	elapsedTime := time.Since(latestBlockTime)
	if elapsedTime > alertLatency {
		ob.logger.Chain.Error().
			Msgf("RPC is stale: latest block is %.0f seconds old, RPC down or chain stuck (check explorer)?", elapsedTime.Seconds())
		return true
	}

	ob.logger.Chain.Info().Msgf("RPC is OK: latest block is %.0f seconds old", elapsedTime.Seconds())
	return false
}

// EnvVarLatestBlockByChain returns the environment variable for the last block by chain.
func EnvVarLatestBlockByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_BLOCK", chain.ChainId)
}

// EnvVarLatestTxByChain returns the environment variable for the last tx by chain.
func EnvVarLatestTxByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_TX", chain.ChainId)
}
