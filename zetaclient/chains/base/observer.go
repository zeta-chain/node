package base

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	// anyStringMap is a key-value map to store any string values used by the observer
	// it is now only used by Sui observer to store old/new Sui gateway inbound cursors
	anyStringMap map[string]string

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
}

// NewObserver creates a new base observer.
func NewObserver(
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
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
		chain:            chain,
		chainParams:      chainParams,
		zetacoreClient:   zetacoreClient,
		tss:              tss,
		lastBlock:        0,
		lastBlockScanned: 0,
		lastTxScanned:    "",
		anyStringMap:     make(map[string]string),
		ts:               ts,
		db:               database,
		blockCache:       blockCache,
		mu:               &sync.Mutex{},
		logger:           newObserverLogger(chain, logger),
		stop:             make(chan struct{}),
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
		ob.logger.Chain.Info().Msg("Observer already stopped")
		return
	}

	ob.logger.Chain.Info().Msg("Stopping observer")

	close(ob.stop)
	ob.started = false

	// close database
	if err := ob.db.Close(); err != nil {
		ob.Logger().Chain.Error().Err(err).Msg("Unable to close db")
	}

	ob.Logger().Chain.Info().Msgf("observer stopped")
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

	ob.logger.Chain.Info().Any("observer.chain_params", params).Msg("updated chain params")
}

// ZetacoreClient returns the zetacore client for the observer.
func (ob *Observer) ZetacoreClient() interfaces.ZetacoreClient {
	return ob.zetacoreClient
}

// TSS returns the tss signer for the observer.
func (ob *Observer) TSS() interfaces.TSSSigner {
	return ob.tss
}

// TSSAddressString returns the TSS address for the chain.
//
// Note: all chains uses TSS EVM address except Bitcoin chain.
func (ob *Observer) TSSAddressString() string {
	switch ob.chain.Consensus {
	case chains.Consensus_bitcoin:
		address, err := ob.tss.PubKey().AddressBTC(ob.Chain().ChainId)
		if err != nil {
			return ""
		}
		return address.EncodeAddress()
	default:
		return ob.tss.PubKey().AddressEVM().String()
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

// GetAnyString get any string data by key
func (ob *Observer) GetAnyString(key string) string {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.anyStringMap[key]
}

// WithAnyString set any string data by key
func (ob *Observer) WithAnyString(key, value string) *Observer {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if ob.anyStringMap[key] == "" {
		ob.logger.Chain.Info().Str("key", key).Str("value", value).Msg("initializing any string value")
	}
	ob.anyStringMap[key] = value

	return ob
}

// WriteAnyStringToDB writes the any string data to the database.
func (ob *Observer) WriteAnyStringToDB(key, value string) error {
	// handle both insert and update cases
	anyString := clienttypes.ToAnyStringSQLType(key, value)
	return ob.db.Client().Where("key_name = ?", key).Assign(anyString).FirstOrCreate(anyString).Error
}

// LoadAnyString loads any string data from environment variable or from database.
func (ob *Observer) LoadAnyString(key string) {
	// get environment variable
	envvar := EnvVarLatestAnyStringByChain(ob.chain, key)
	value := os.Getenv(envvar)

	// load from environment variable if set
	if value != "" {
		ob.logger.Chain.Info().Str("envvar", envvar).Str("value", value).Msg("environment variable is set")
		ob.WithAnyString(key, value)
		return
	}

	// load from DB otherwise
	value, err := ob.ReadAnyStringFromDB(key)
	if err != nil {
		// if not found, let the concrete chain observer decide where to start
		chainID := ob.chain.ChainId
		ob.logger.Chain.Info().Int64(logs.FieldChain, chainID).Str("key", key).Msg("string value not found in db")
		return
	}
	ob.WithAnyString(key, value)
}

// ReadAnyStringFromDB reads the any string data from the database.
func (ob *Observer) ReadAnyStringFromDB(key string) (string, error) {
	var anyString clienttypes.AnyStringSQLType
	if err := ob.db.Client().Where("key_name = ?", key).First(&anyString).Error; err != nil {
		return "", err
	}
	return anyString.Value, nil
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

	// prepare logger fields
	lf := map[string]any{
		logs.FieldMethod:           "PostVoteInbound",
		logs.FieldTx:               txHash,
		logs.FieldCoinType:         coinType.String(),
		logs.FieldConfirmationMode: msg.ConfirmationMode.String(),
	}

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
				ob.logger.Inbound.Info().Fields(lf).Msg("inbound detected: cctx exists but the ballot does not")
				return cctxIndex, nil
			}
		}
	}

	// make sure the message is valid to avoid unnecessary retries
	if err := msg.ValidateBasic(); err != nil {
		ob.logger.Inbound.Warn().Err(err).Fields(lf).Msg("invalid inbound vote message")
		return "", nil
	}

	// post vote to zetacore
	zetaHash, ballot, err := ob.ZetacoreClient().PostVoteInbound(ctx, gasLimit, retryGasLimit, msg)
	lf[logs.FieldZetaTx] = zetaHash
	lf[logs.FieldBallot] = ballot

	switch {
	case err != nil:
		ob.logger.Inbound.Error().Err(err).Fields(lf).Msg("inbound detected: error posting vote")
		return "", err
	case zetaHash == "":
		ob.logger.Inbound.Info().Fields(lf).Msg("inbound detected: already voted on ballot")
	default:
		ob.logger.Inbound.Info().Fields(lf).Msgf("inbound detected: vote posted")
	}

	return ballot, nil
}

// EnvVarLatestBlockByChain returns the environment variable for the last block by chain.
func EnvVarLatestBlockByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_BLOCK", chain.ChainId)
}

// EnvVarLatestTxByChain returns the environment variable for the last tx by chain.
func EnvVarLatestTxByChain(chain chains.Chain) string {
	return fmt.Sprintf("CHAIN_%d_SCAN_FROM_TX", chain.ChainId)
}

// EnvVarLatestAnyStringByChain returns the environment variable for any string data by chain for the given key.
func EnvVarLatestAnyStringByChain(chain chains.Chain, key string) string {
	return fmt.Sprintf("CHAIN_%d_ANY_STRING_%s", chain.ChainId, key)
}

func newObserverLogger(chain chains.Chain, logger Logger) ObserverLogger {
	withLogFields := func(l zerolog.Logger) zerolog.Logger {
		return l.With().
			Int64(logs.FieldChain, chain.ChainId).
			Str(logs.FieldChainNetwork, chain.Network.String()).
			Logger()
	}

	log := withLogFields(logger.Std)
	complianceLog := withLogFields(logger.Compliance)

	return ObserverLogger{
		Chain:      log,
		Inbound:    log.With().Str(logs.FieldModule, logs.ModNameInbound).Logger(),
		Outbound:   log.With().Str(logs.FieldModule, logs.ModNameOutbound).Logger(),
		GasPrice:   log.With().Str(logs.FieldModule, logs.ModNameGasPrice).Logger(),
		Headers:    log.With().Str(logs.FieldModule, logs.ModNameHeaders).Logger(),
		Compliance: complianceLog,
	}
}
