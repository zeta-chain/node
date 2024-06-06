package base

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	// EnvVarLatestBlock is the environment variable to indicate latest block height
	EnvVarLatestBlock = "latest"

	// DefaultBlockCacheSize is the default size of the block cache
	DefaultBlockCacheSize = 1000
)

// Observer is the base observer.
type Observer struct {
	// chain contains static information about the observed chain
	chain chains.Chain

	// chainParams contains the dynamic chain parameters of the observed chain
	chainParams observertypes.ChainParams

	// coreContext contains context data of ZetaChain
	zetacoreContext *context.ZetacoreContext

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

	// db is the database to persist data
	db *gorm.DB

	// logger contains the loggers used by observer
	logger ObserverLogger

	// stop is the channel to signal the observer to stop
	stop chan struct{}
}

// NewObserver creates a new base observer.
func NewObserver(
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	zetacoreContext *context.ZetacoreContext,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	blockCacheSize int,
	dbPath string,
	logger Logger,
) (*Observer, error) {
	ob := Observer{
		chain:            chain,
		chainParams:      chainParams,
		zetacoreContext:  zetacoreContext,
		zetacoreClient:   zetacoreClient,
		tss:              tss,
		lastBlock:        0,
		lastBlockScanned: 0,
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

	// open database
	err = ob.OpenDB(dbPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error opening observer db for chain: %s", chain.ChainName))
	}

	return &ob, nil
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

// ZetacoreContext returns the zetacore context for the observer.
func (ob *Observer) ZetacoreContext() *context.ZetacoreContext {
	return ob.zetacoreContext
}

// WithZetacoreContext attaches a new zetacore context to the observer.
func (ob *Observer) WithZetacoreContext(context *context.ZetacoreContext) *Observer {
	ob.zetacoreContext = context
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

// DB returns the database for the observer.
func (ob *Observer) DB() *gorm.DB {
	return ob.db
}

// Logger returns the logger for the observer.
func (ob *Observer) Logger() ObserverLogger {
	return ob.logger
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

// Stop returns the stop channel for the observer.
func (ob *Observer) Stop() chan struct{} {
	return ob.stop
}

// OpenDB open sql database in the given path.
func (ob *Observer) OpenDB(dbPath string) error {
	if dbPath != "" {
		// create db path if not exist
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			err := os.MkdirAll(dbPath, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "error creating db path")
			}
		}

		// open db by chain name
		chainName := ob.chain.ChainName.String()
		path := fmt.Sprintf("%s/%s", dbPath, chainName)
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return errors.Wrap(err, "error opening db")
		}

		// migrate db
		err = db.AutoMigrate(&clienttypes.ReceiptSQLType{},
			&clienttypes.TransactionSQLType{},
			&clienttypes.LastBlockSQLType{})
		if err != nil {
			return errors.Wrap(err, "error migrating db")
		}
		ob.db = db
	}
	return nil
}

// LoadLastBlockScanned loads last scanned block from environment variable or from database.
// The last scanned block is the height from which the observer should continue scanning.
func (ob *Observer) LoadLastBlockScanned(logger zerolog.Logger) (fromLatest bool, err error) {
	// get environment variable
	envvar := EnvVarLatestBlockByChain(ob.chain)
	scanFromBlock := os.Getenv(envvar)

	// load from environment variable if set
	if scanFromBlock != "" {
		logger.Info().
			Msgf("LoadLastBlockScanned: envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == EnvVarLatestBlock {
			return true, nil
		}
		blockNumber, err := strconv.ParseUint(scanFromBlock, 10, 64)
		if err != nil {
			return false, err
		}
		ob.WithLastBlockScanned(blockNumber)
		return false, nil
	}

	// load from DB otherwise. If not found, start from latest block
	blockNumber, err := ob.ReadLastBlockScannedFromDB()
	if err != nil {
		logger.Info().Msgf("LoadLastBlockScanned: chain %d starts scanning from latest block", ob.chain.ChainId)
		return true, nil
	}
	ob.WithLastBlockScanned(blockNumber)
	logger.Info().
		Msgf("LoadLastBlockScanned: chain %d starts scanning from block %d", ob.chain.ChainId, ob.LastBlockScanned())

	return false, nil
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
	return chain.ChainName.String() + "_SCAN_FROM"
}
