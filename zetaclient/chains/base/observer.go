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
	rpcAlertLatency int64,
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
		rpcAlertLatency:  time.Duration(rpcAlertLatency) * time.Second,
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
	ob.mu.Lock()
	defer ob.mu.Unlock()

	return ob.chainParams
}

// SetChainParams attaches a new chain params to the observer.
func (ob *Observer) SetChainParams(params observertypes.ChainParams) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.chainParams = params
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

// TSS returns the tss signer for the observer.
func (ob *Observer) TSS() interfaces.TSSSigner {
	return ob.tss
}

// WithTSS attaches a new tss signer to the observer.
func (ob *Observer) WithTSS(tss interfaces.TSSSigner) *Observer {
	ob.tss = tss
	return ob
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

// IsBlockConfirmed checks if the given block number is confirmed.
//
// Note: block 100 is confirmed if the last block is 100 and confirmation count is 1.
func (ob *Observer) IsBlockConfirmed(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confBlock := blockNumber + ob.chainParams.ConfirmationCount - 1
	return lastBlock >= confBlock
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
		logs.FieldMethod:   "PostVoteInbound",
		logs.FieldTx:       txHash,
		logs.FieldCoinType: coinType.String(),
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

// AlertOnRPCLatency prints an alert if the RPC latency exceeds the threshold.
// Returns true if the RPC latency is too high.
func (ob *Observer) AlertOnRPCLatency(latestBlockTime time.Time, defaultAlertLatency time.Duration) bool {
	elapsedTime := time.Since(latestBlockTime)

	alertLatency := ob.rpcAlertLatency
	if alertLatency == 0 {
		alertLatency = defaultAlertLatency
	}

	lf := map[string]any{
		"rpc_latency_alert_ms": alertLatency.Milliseconds(),
		"rpc_latency_real_ms":  elapsedTime.Milliseconds(),
	}

	if elapsedTime > alertLatency {
		ob.logger.Chain.Error().Fields(lf).Msg("RPC latency is too high, please check the node or explorer")
		return true
	}

	ob.logger.Chain.Info().Fields(lf).Msg("RPC latency is OK")

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
