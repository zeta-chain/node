package base

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
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

	// blockCache is the cache for blocks
	blockCache *lru.Cache

	// headerCache is the cache for headers
	headerCache *lru.Cache

	// db is the database to persist data
	db *gorm.DB

	// ts is the telemetry server for metrics
	ts *metrics.TelemetryServer

	// logger contains the loggers used by observer
	logger ObserverLogger

	// mu protects fields from concurrent access
	// Note: base observer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu *sync.Mutex

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
	ts *metrics.TelemetryServer,
	logger Logger,
) (*Observer, error) {
	ob := Observer{
		chain:            chain,
		chainParams:      chainParams,
		zetacoreClient:   zetacoreClient,
		tss:              tss,
		lastBlock:        0,
		lastBlockScanned: 0,
		ts:               ts,
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

// Stop notifies all goroutines to stop and closes the database.
func (ob *Observer) Stop() {
	ob.logger.Chain.Info().Msgf("observer is stopping for chain %d", ob.Chain().ChainId)
	close(ob.stop)

	// close database
	if ob.db != nil {
		err := ob.CloseDB()
		if err != nil {
			ob.Logger().Chain.Error().Err(err).Msgf("CloseDB failed for chain %d", ob.Chain().ChainId)
		}
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
	metrics.LastScannedBlockNumber.WithLabelValues(ob.chain.ChainName.String()).Set(float64(blockNumber))
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

// DB returns the database for the observer.
func (ob *Observer) DB() *gorm.DB {
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
	chainLogger := logger.Std.With().Int64("chain", ob.chain.ChainId).Logger()
	ob.logger = ObserverLogger{
		Chain:      chainLogger,
		Inbound:    chainLogger.With().Str("module", "inbound").Logger(),
		Outbound:   chainLogger.With().Str("module", "outbound").Logger(),
		GasPrice:   chainLogger.With().Str("module", "gasprice").Logger(),
		Headers:    chainLogger.With().Str("module", "headers").Logger(),
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

// OpenDB open sql database in the given path.
func (ob *Observer) OpenDB(dbPath string, dbName string) error {
	// create db path if not exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err := os.MkdirAll(dbPath, 0o750)
		if err != nil {
			return errors.Wrapf(err, "error creating db path: %s", dbPath)
		}
	}

	// use custom dbName or chain name if not provided
	if dbName == "" {
		dbName = ob.chain.ChainName.String()
	}
	path := fmt.Sprintf("%s/%s", dbPath, dbName)

	// use memory db if specified
	if strings.Contains(dbPath, ":memory:") {
		path = dbPath
	}

	// open db
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return errors.Wrap(err, "error opening db")
	}

	// migrate db
	err = db.AutoMigrate(&clienttypes.LastBlockSQLType{})
	if err != nil {
		return errors.Wrap(err, "error migrating db")
	}
	ob.db = db

	return nil
}

// CloseDB close the database.
func (ob *Observer) CloseDB() error {
	dbInst, err := ob.db.DB()
	if err != nil {
		return fmt.Errorf("error getting database instance: %w", err)
	}
	err = dbInst.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}
	return nil
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
			return err
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
	logger.Info().
		Msgf("LoadLastBlockScanned: chain %d starts scanning from block %d", ob.chain.ChainId, ob.LastBlockScanned())

	return nil
}

// SaveLastBlockScanned saves the last scanned block to memory and database.
func (ob *Observer) SaveLastBlockScanned(blockNumber uint64) error {
	ob.WithLastBlockScanned(blockNumber)
	return ob.WriteLastBlockScannedToDB(blockNumber)
}

// WriteLastBlockScannedToDB saves the last scanned block to the database.
func (ob *Observer) WriteLastBlockScannedToDB(lastScannedBlock uint64) error {
	return ob.db.Save(clienttypes.ToLastBlockSQLType(lastScannedBlock)).Error
}

// ReadLastBlockScannedFromDB reads the last scanned block from the database.
func (ob *Observer) ReadLastBlockScannedFromDB() (uint64, error) {
	var lastBlock clienttypes.LastBlockSQLType
	if err := ob.db.First(&lastBlock, clienttypes.LastBlockNumID).Error; err != nil {
		// record not found
		return 0, err
	}
	return lastBlock.Num, nil
}

// EnvVarLatestBlock returns the environment variable for the latest block by chain.
func EnvVarLatestBlockByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM", chain.ChainId)
}
