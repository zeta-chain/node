package base

import (
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
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

const (
	// EnvVarLatestBlock is the environment variable that forces the observer to scan from the latest block
	EnvVarLatestBlock = "latest"

	// DefaultBlockCacheSize is the default number of blocks that the observer will keep in cache for performance (without RPC calls)
	// Cached blocks can be used to get block information and verify transactions
	DefaultBlockCacheSize = 1000

	// MonitoringErrHandlerRoutineTimeout is the timeout for the handleMonitoring routine that waits for an error from the monitorVote channel
	MonitoringErrHandlerRoutineTimeout = 5 * time.Minute

	// defaultSampledLogInterval is the default interval for sampled logs
	defaultSampledLogInterval = 10
)

// Observer is the base structure for chain observers, grouping the common logic for each chain observer client.
// The common logic includes: chain, chainParams, contexts, zetacore client, tss, lastBlock, db, metrics, loggers etc.
type Observer struct {
	// chain contains static information about the observed chain
	chain chains.Chain

	// chainParams contains the dynamic chain parameters of the observed chain
	chainParams observertypes.ChainParams

	// zetaRepo is the repository that interacts with zetacore
	zetaRepo *zrepo.ZetaRepo

	// tssSigner is the TSS signer
	tssSigner tssrepo.TSSClient

	// lastBlock is the last block height of the observed chain
	lastBlock uint64

	// lastBlockScanned is the last block height scanned by the observer
	lastBlockScanned uint64

	// lastTxScanned is the last transaction hash scanned by the observer
	lastTxScanned string

	blockCache *lru.Cache

	// internalInboundTrackers stores trackers for inbounds that failed to vote on due to broadcasting error (e.g. tx dropped)
	// the contents of the map may vary from observer to observer, depending on individual situation
	internalInboundTrackers map[string]crosschaintypes.InboundTracker

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
	zetaRepo *zrepo.ZetaRepo,
	tssSigner tssrepo.TSSClient,
	blockCacheSize int,
	ts *metrics.TelemetryServer,
	database *db.DB,
	logger Logger,
) (*Observer, error) {
	blockCache, err := lru.New(blockCacheSize)
	if err != nil {
		return nil, errors.Wrap(err, "error creating block cache")
	}

	return &Observer{
		chain:                   chain,
		chainParams:             chainParams,
		zetaRepo:                zetaRepo,
		tssSigner:               tssSigner,
		lastBlock:               0,
		lastBlockScanned:        0,
		lastTxScanned:           "",
		ts:                      ts,
		db:                      database,
		blockCache:              blockCache,
		internalInboundTrackers: make(map[string]crosschaintypes.InboundTracker),
		mu:                      &sync.Mutex{},
		logger:                  newObserverLogger(chain, logger),
		stop:                    make(chan struct{}),
	}, nil
}

// Start starts the observer. Returns false if it's already started (noop).
func (ob *Observer) Start() bool {
	ob.mu.Lock()
	defer ob.Mu().Unlock()

	// noop
	if ob.started {
		return false
	}

	ob.started = true

	return true
}

// Stop notifies all goroutines to stop and closes the database.
func (ob *Observer) Stop() {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if !ob.started {
		ob.logger.Chain.Info().Msg("observer already stopped")
		return
	}

	ob.logger.Chain.Info().Msg("stopping the observer")

	close(ob.stop)
	ob.started = false

	// close database
	if err := ob.db.Close(); err != nil {
		ob.Logger().Chain.Error().Err(err).Msg("unable to close database")
	}

	ob.Logger().Chain.Info().Msg("stopped the observer")
}

// Chain returns the chain for the observer.
func (ob *Observer) Chain() chains.Chain {
	return ob.chain
}

// ChainParams returns the chain params for the observer.
func (ob *Observer) ChainParams() observertypes.ChainParams {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	return ob.chainParams
}

// SetChainParams attaches a new chain params to the observer.
func (ob *Observer) SetChainParams(params observertypes.ChainParams) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if observertypes.ChainParamsEqual(ob.chainParams, params) {
		return
	}

	ob.chainParams = params

	ob.logger.Chain.Info().Any("chain_params", params).Msg("updated chain parameters")
}

// ZetaRepo returns the zrepo.ZetaRepo repository for the observer.
func (ob *Observer) ZetaRepo() *zrepo.ZetaRepo {
	return ob.zetaRepo
}

// TSS returns the tss signer for the observer.
func (ob *Observer) TSS() tssrepo.TSSClient {
	return ob.tssSigner
}

// TSSAddressString returns the TSS address for the chain.
//
// Note: all chains uses TSS EVM address except Bitcoin chain.
func (ob *Observer) TSSAddressString() string {
	switch ob.chain.Consensus {
	case chains.Consensus_bitcoin:
		address, err := ob.tssSigner.PubKey().AddressBTC(ob.Chain().ChainId)
		if err != nil {
			return ""
		}
		return address.EncodeAddress()
	default:
		return ob.tssSigner.PubKey().AddressEVM().String()
	}
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
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if ob.lastTxScanned == "" {
		ob.logger.Chain.Info().
			Str(logs.FieldTx, txHash).
			Msg("initializing last scanned transaction")
	}

	ob.lastTxScanned = txHash
	return ob
}

// BlockCache returns the block cache for the observer.
func (ob *Observer) BlockCache() *lru.Cache {
	return ob.blockCache
}

// OutboundID returns a unique identifier for the outbound transaction.
// The identifier is now used as the key for maps that store outbound related data (e.g. transaction, receipt, etc).
func (ob *Observer) OutboundID(nonce uint64) string {
	tssAddress := ob.TSSAddressString()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddress, nonce)
}

// DB returns the database for the observer.
func (ob *Observer) DB() *db.DB {
	return ob.db
}

// TelemetryServer returns the telemetry server for the observer.
func (ob *Observer) TelemetryServer() *metrics.TelemetryServer {
	return ob.ts
}

// Logger returns the logger for the observer.
func (ob *Observer) Logger() *ObserverLogger {
	return &ob.logger
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
func (ob *Observer) LoadLastBlockScanned() error {
	logger := ob.logger.Chain

	// get environment variable
	envvar := EnvVarLatestBlockByChain(ob.chain)
	scanFromBlock := os.Getenv(envvar)

	// load from environment variable if set
	if scanFromBlock != "" {
		logger.Info().
			Str("envvar", envvar).
			Str(logs.FieldBlock, scanFromBlock).
			Msg("envvar is set; scan from block")
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
		logger.Info().Msg("last scanned block not found in the database")
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
	logger := ob.logger.Chain

	// get environment variable
	envvar := EnvVarLatestTxByChain(ob.chain)
	scanFromTx := os.Getenv(envvar)

	// load from environment variable if set
	if scanFromTx != "" {
		logger.Info().
			Str("envvar", envvar).
			Str(logs.FieldTx, scanFromTx).
			Msg("envvar is set; scan from tx")
		ob.WithLastTxScanned(scanFromTx)
		return
	}

	// load from DB otherwise.
	txHash, err := ob.ReadLastTxScannedFromDB()
	if err != nil {
		// If not found, let the concrete chain observer decide where to start
		logger.Info().Err(err).Msg("last scanned tx not found in the database")
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

// EnvVarLatestBlockByChain returns the environment variable for the last block by chain.
func EnvVarLatestBlockByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_BLOCK", chain.ChainId)
}

// EnvVarLatestTxByChain returns the environment variable for the last tx by chain.
func EnvVarLatestTxByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_TX", chain.ChainId)
}

func newObserverLogger(chain chains.Chain, logger Logger) ObserverLogger {
	withLogFields := func(l zerolog.Logger) zerolog.Logger {
		return l.With().
			Int64(logs.FieldChain, chain.ChainId).
			Stringer(logs.FieldNetwork, chain.Network).
			Logger()
	}

	log := withLogFields(logger.Std)
	complianceLog := withLogFields(logger.Compliance)

	return ObserverLogger{
		Chain:      log,
		Inbound:    log.With().Str(logs.FieldModule, logs.ModNameInbound).Logger(),
		Outbound:   log.With().Str(logs.FieldModule, logs.ModNameOutbound).Logger(),
		Compliance: complianceLog,
		Sampled:    log.Sample(&zerolog.BasicSampler{N: defaultSampledLogInterval}),
	}
}
