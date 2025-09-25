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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/node/zetaclient/context"
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

	// MonitoringErrHandlerRoutineTimeout is the timeout for the handleMonitoring routine that waits for an error from the monitorVote channel
	MonitoringErrHandlerRoutineTimeout = 5 * time.Minute
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

	// tssSigner is the TSS signer
	tssSigner interfaces.TSSSigner

	// lastBlock is the last block height of the observed chain
	lastBlock uint64

	// lastBlockScanned is the last block height scanned by the observer
	lastBlockScanned uint64

	// lastTxScanned is the last transaction hash scanned by the observer
	lastTxScanned string

	// auxStringMap is a key-value map to store any auxiliary string values used by the observer
	// it is now only used by Sui observer to store old/new Sui gateway inbound cursors
	auxStringMap map[string]string

	blockCache *lru.Cache

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

	forceResetLastScanned bool
}

// NewObserver creates a new base observer.
func NewObserver(
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tssSigner interfaces.TSSSigner,
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
		chain:                 chain,
		chainParams:           chainParams,
		zetacoreClient:        zetacoreClient,
		tssSigner:             tssSigner,
		lastBlock:             0,
		lastBlockScanned:      0,
		lastTxScanned:         "",
		auxStringMap:          make(map[string]string),
		ts:                    ts,
		db:                    database,
		blockCache:            blockCache,
		mu:                    &sync.Mutex{},
		logger:                newObserverLogger(chain, logger),
		stop:                  make(chan struct{}),
		forceResetLastScanned: false,
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

// ZetacoreClient returns the zetacore client for the observer.
func (ob *Observer) ZetacoreClient() interfaces.ZetacoreClient {
	return ob.zetacoreClient
}

// TSS returns the tss signer for the observer.
func (ob *Observer) TSS() interfaces.TSSSigner {
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
	ob.mu.Lock()
	defer ob.mu.Unlock()
	height := atomic.LoadUint64(&ob.lastBlockScanned)
	return height
}

// WithLastBlockScanned set last block scanned (not necessarily caught up with the chain; could be slow/paused).
// it also set the value of forceResetLastScanned and returns the previous value.
// If forceResetLastScanned was true before, it means the monitoring thread would have updated it and so it skips updating the last scanned block.
func (ob *Observer) WithLastBlockScanned(blockNumber uint64, forceResetLastScanned bool) (*Observer, bool) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	wasForceReset := ob.forceResetLastScanned
	ob.forceResetLastScanned = forceResetLastScanned

	// forceResetLastScanned was set to true before; it means the monitoring thread would have updated it
	// In this case we should not update the last scanned block and just return
	if wasForceReset && !forceResetLastScanned {
		return ob, wasForceReset
	}

	atomic.StoreUint64(&ob.lastBlockScanned, blockNumber)
	metrics.LastScannedBlockNumber.WithLabelValues(ob.chain.Name).Set(float64(blockNumber))
	return ob, wasForceReset
}

// LastTxScanned get last transaction scanned.
func (ob *Observer) LastTxScanned() string {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.lastTxScanned
}

// WithLastTxScanned set last transaction scanned.
func (ob *Observer) WithLastTxScanned(txHash string) *Observer {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if ob.lastTxScanned == "" {
		ob.logger.Chain.Info().Str("tx", txHash).Msg("initializing last tx scanned")
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
		ob.WithLastBlockScanned(blockNumber, false)
		return nil
	}

	// load from DB otherwise. If not found, start from latest block
	blockNumber, err := ob.ReadLastBlockScannedFromDB()
	if err != nil {
		logger.Info().Msg("last scanned block not found in the database")
		return nil
	}
	ob.WithLastBlockScanned(blockNumber, false)

	return nil
}

// SaveLastBlockScanned saves the last scanned block to memory and database.
func (ob *Observer) SaveLastBlockScanned(blockNumber uint64) error {
	_, forceResetLastScannedBeforeUpdate := ob.WithLastBlockScanned(blockNumber, false)
	if forceResetLastScannedBeforeUpdate {
		return nil
	}
	return ob.WriteLastBlockScannedToDB(blockNumber)
}

// ForceSaveLastBlockScanned saves the last scanned block to memory if the new blocknumber is less than the current last scanned block.
// It also forces the update of the last scanned block in the database, to makes sure any other the block gets rescanned.
func (ob *Observer) ForceSaveLastBlockScanned(blockNumber uint64) error {
	currentLastScanned := ob.LastBlockScanned()
	if blockNumber > currentLastScanned {
		return nil
	}
	ob.WithLastBlockScanned(blockNumber, true)
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
	ob.WithLastBlockScanned(slot, false)

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

// GetAuxString get any auxiliary string data by key
func (ob *Observer) GetAuxString(key string) string {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.auxStringMap[key]
}

// WithAuxString set any auxiliary string data by key
func (ob *Observer) WithAuxString(key, value string) *Observer {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if ob.auxStringMap[key] == "" {
		ob.logger.Chain.Info().Str("key", key).Str("value", value).Msg("initializing auxiliary string value")
	}
	ob.auxStringMap[key] = value

	return ob
}

// WriteAuxStringToDB writes the auxiliary string data to the database.
func (ob *Observer) WriteAuxStringToDB(key, value string) error {
	// create new record if not found
	var existingRecord clienttypes.AuxStringSQLType
	if err := ob.db.Client().Where("key_name = ?", key).First(&existingRecord).Error; err != nil {
		return ob.db.Client().Create(clienttypes.ToAuxStringSQLType(key, value)).Error
	}

	// record exists, update it
	return ob.db.Client().Model(&existingRecord).Update("value", value).Error
}

// LoadAuxString loads auxiliary string data from environment variable or from database.
func (ob *Observer) LoadAuxString(key string) {
	// get environment variable
	envvar := EnvVarLatestAuxStringByChain(ob.chain, key)
	value := os.Getenv(envvar)

	// load from environment variable if set
	if value != "" {
		ob.logger.Chain.Info().Str("envvar", envvar).Str("value", value).Msg("environment variable is set")
		ob.WithAuxString(key, value)
		return
	}

	// load from DB otherwise
	value, err := ob.ReadAuxStringFromDB(key)
	if err != nil {
		// if not found, let the concrete chain observer decide where to start
		chainID := ob.chain.ChainId
		ob.logger.Chain.Info().Int64(logs.FieldChain, chainID).Str("key", key).Msg("string value not found in db")
		return
	}
	ob.WithAuxString(key, value)
}

// ReadAuxStringFromDB reads the auxiliary string data from the database.
func (ob *Observer) ReadAuxStringFromDB(key string) (string, error) {
	var record clienttypes.AuxStringSQLType
	if err := ob.db.Client().Where("key_name = ?", key).First(&record).Error; err != nil {
		return "", err
	}
	return record.Value, nil
}

// PostVoteInbound posts a vote for the given vote message and returns the ballot.
func (ob *Observer) PostVoteInbound(
	ctx context.Context,
	msg *crosschaintypes.MsgVoteInbound,
	retryGasLimit uint64,
) (string, error) {
	const gasLimit = zetacore.PostVoteInboundGasLimit

	var (
		txHash   = msg.InboundHash
		coinType = msg.CoinType
	)

	logger := ob.logger.Inbound.With().
		Str(logs.FieldTx, txHash).
		Stringer(logs.FieldCoinType, coinType).
		Stringer("confirmation_mode", msg.ConfirmationMode).
		Logger()

	cctxIndex := msg.Digest()
	// The cctx is created after the inbound ballot is finalized
	// 1. if the cctx already exists, we could try voting if the ballot is present
	// 2. if the cctx exists but the ballot does not exist, we do not need to vote
	_, err := ob.ZetacoreClient().GetCctxByHash(ctx, cctxIndex)
	if err == nil {
		// The cctx exists we should still vote if the ballot is present
		_, ballotErr := ob.ZetacoreClient().GetBallotByID(ctx, cctxIndex)
		if ballotErr != nil {
			// Verify ballot is not found
			if st, ok := status.FromError(ballotErr); ok && st.Code() == codes.NotFound {
				// Query for ballot failed, the ballot does not exist we can return
				logger.Info().Msg("inbound detected: CCTX exists but the ballot does not")
				return cctxIndex, nil
			}
		}
	}

	// make sure the message is valid to avoid unnecessary retries
	if err := msg.ValidateBasic(); err != nil {
		logger.Warn().Err(err).Msg("invalid inbound vote message")
		return "", nil
	}

	monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)

	// ctxWithTimeout is a context with timeout used for monitoring the vote transaction
	ctxWithTimeout, _ := zctx.CopyWithTimeout(ctx, context.Background(), MonitoringErrHandlerRoutineTimeout)

	// post vote to zetacore
	zetaHash, ballot, err := ob.ZetacoreClient().
		PostVoteInbound(ctxWithTimeout, gasLimit, retryGasLimit, msg, monitorErrCh)

	logger = logger.With().
		Str(logs.FieldZetaTx, zetaHash).
		Str(logs.FieldBallotIndex, ballot).
		Logger()

	switch {
	case err != nil:
		logger.Error().Err(err).Msg("inbound detected: error posting vote")
		return "", err
	case zetaHash == "":
		logger.Info().Msg("inbound detected: already voted on ballot")
	default:
		logger.Info().Msg("inbound detected: vote posted")
	}

	go func() {
		ob.handleMonitoringError(ctxWithTimeout, monitorErrCh, zetaHash)
	}()

	return ballot, nil
}

func (ob *Observer) handleMonitoringError(
	ctx context.Context,
	monitorErrCh <-chan zetaerrors.ErrTxMonitor,
	zetaHash string,
) {
	logger := ob.logger.Inbound
	defer func() {
		if r := recover(); r != nil {
			logger.Error().Any("panic", r).Msg("recovered from panic in monitoring error handler")
		}
	}()

	select {
	case monitorErr := <-monitorErrCh:
		if monitorErr.Err != nil {
			logger.Error().
				Err(monitorErr).
				Str(logs.FieldZetaTx, monitorErr.ZetaTxHash).
				Str(logs.FieldBallotIndex, monitorErr.BallotIndex).
				Uint64(logs.FieldBlock, monitorErr.InboundBlockHeight).
				Msg("error monitoring vote transaction")

			if monitorErr.InboundBlockHeight > 0 {
				logger.Info().Uint64(logs.FieldBlock, monitorErr.InboundBlockHeight-1).
					Str(logs.FieldBallotIndex, monitorErr.BallotIndex).
					Uint64(logs.FieldBlock, monitorErr.InboundBlockHeight).
					Msg("reset last scanned block")
				// save last scanned block as the block before the inbound block height
				err := ob.ForceSaveLastBlockScanned(monitorErr.InboundBlockHeight - 1)
				if err != nil {
					logger.Error().Err(err).
						Str(logs.FieldZetaTx, monitorErr.ZetaTxHash).
						Msg("unable to save last scanned block after monitoring error")
				}
			}
		}
	case <-ctx.Done():
		logger.Debug().
			Str(logs.FieldZetaTx, zetaHash).
			Msg("no error received for the monitoring, the transaction likely succeeded")
	}
}

// EnvVarLatestBlockByChain returns the environment variable for the last block by chain.
func EnvVarLatestBlockByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_BLOCK", chain.ChainId)
}

// EnvVarLatestTxByChain returns the environment variable for the last tx by chain.
func EnvVarLatestTxByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_TX", chain.ChainId)
}

// EnvVarLatestAuxStringByChain returns the environment variable for auxiliary string data by chain for the given key.
func EnvVarLatestAuxStringByChain(chain chains.Chain, key string) string {
	return fmt.Sprintf("CHAIN_%d_AUX_STRING_%s", chain.ChainId, key)
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
	}
}
