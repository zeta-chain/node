package evm

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.non-eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
)

// Logger is the logger for evm chains
// TODO: Merge this logger with the one in bitcoin
// https://github.com/zeta-chain/node/issues/2022
type Logger struct {
	// Chain is the parent logger for the chain
	Chain zerolog.Logger

	// Inbound is the logger for incoming transactions
	Inbound zerolog.Logger

	// Outbound is the logger for outgoing transactions
	Outbound zerolog.Logger

	// GasPrice is the logger for gas prices
	GasPrice zerolog.Logger

	// Compliance is the logger for compliance checks
	Compliance zerolog.Logger
}

var _ interfaces.ChainClient = &ChainClient{}

// ChainClient represents the chain configuration for an EVM chain
// Filled with above constants depending on chain
type ChainClient struct {
	Tss interfaces.TSSSigner

	Mu *sync.Mutex

	chain                         chains.Chain
	evmClient                     interfaces.EVMRPCClient
	evmJSONRPC                    interfaces.EVMJSONRPCClient
	zetaBridge                    interfaces.ZetaCoreBridger
	lastBlockScanned              uint64
	lastBlock                     uint64
	db                            *gorm.DB
	outboundPendingTransactions   map[string]*ethtypes.Transaction
	outboundConfirmedReceipts     map[string]*ethtypes.Receipt
	outboundConfirmedTransactions map[string]*ethtypes.Transaction
	stop                          chan struct{}
	logger                        Logger
	coreContext                   *corecontext.ZetaCoreContext
	chainParams                   observertypes.ChainParams
	ts                            *metrics.TelemetryServer

	blockCache  *lru.Cache
	headerCache *lru.Cache
}

// NewEVMChainClient returns a new configuration based on supplied target chain
func NewEVMChainClient(
	appContext *appcontext.AppContext,
	bridge interfaces.ZetaCoreBridger,
	tss interfaces.TSSSigner,
	dbpath string,
	loggers clientcommon.ClientLogger,
	evmCfg config.EVMConfig,
	ts *metrics.TelemetryServer,
) (*ChainClient, error) {
	ob := ChainClient{
		ts: ts,
	}

	chainLogger := loggers.Std.With().Str("chain", evmCfg.Chain.ChainName.String()).Logger()
	ob.logger = Logger{
		Chain:      chainLogger,
		Inbound:    chainLogger.With().Str("module", "WatchInbound").Logger(),
		Outbound:   chainLogger.With().Str("module", "WatchOutbound").Logger(),
		GasPrice:   chainLogger.With().Str("module", "WatchGasPrice").Logger(),
		Compliance: loggers.Compliance,
	}

	ob.coreContext = appContext.ZetaCoreContext()
	chainParams, found := ob.coreContext.GetEVMChainParams(evmCfg.Chain.ChainId)
	if !found {
		return nil, fmt.Errorf("evm chains params not initialized for chain %d", evmCfg.Chain.ChainId)
	}

	ob.chainParams = *chainParams
	ob.stop = make(chan struct{})
	ob.chain = evmCfg.Chain
	ob.Mu = &sync.Mutex{}
	ob.zetaBridge = bridge
	ob.Tss = tss
	ob.outboundPendingTransactions = make(map[string]*ethtypes.Transaction)
	ob.outboundConfirmedReceipts = make(map[string]*ethtypes.Receipt)
	ob.outboundConfirmedTransactions = make(map[string]*ethtypes.Transaction)

	ob.logger.Chain.Info().Msgf("Chain %s endpoint %s", ob.chain.ChainName.String(), evmCfg.Endpoint)
	client, err := ethclient.Dial(evmCfg.Endpoint)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}

	ob.evmClient = client
	ob.evmJSONRPC = ethrpc.NewEthRPC(evmCfg.Endpoint)

	// create block header and block caches
	ob.blockCache, err = lru.New(1000)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("failed to create block cache")
		return nil, err
	}

	ob.headerCache, err = lru.New(1000)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("failed to create header cache")
		return nil, err
	}

	err = ob.LoadDB(dbpath, ob.chain)
	if err != nil {
		return nil, err
	}

	ob.logger.Chain.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.GetLastBlockHeightScanned())

	return &ob, nil
}

// WithChain attaches a new chain to the chain client
func (ob *ChainClient) WithChain(chain chains.Chain) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.chain = chain
}

// WithLogger attaches a new logger to the chain client
func (ob *ChainClient) WithLogger(logger zerolog.Logger) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.logger = Logger{
		Chain:    logger,
		Inbound:  logger.With().Str("module", "WatchInbound").Logger(),
		Outbound: logger.With().Str("module", "WatchOutbound").Logger(),
		GasPrice: logger.With().Str("module", "WatchGasPrice").Logger(),
	}
}

// WithEvmClient attaches a new evm client to the chain client
func (ob *ChainClient) WithEvmClient(client interfaces.EVMRPCClient) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.evmClient = client
}

// WithEvmJSONRPC attaches a new evm json rpc client to the chain client
func (ob *ChainClient) WithEvmJSONRPC(client interfaces.EVMJSONRPCClient) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.evmJSONRPC = client
}

// WithZetaBridge attaches a new bridge to interact with ZetaCore to the chain client
func (ob *ChainClient) WithZetaBridge(bridge interfaces.ZetaCoreBridger) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.zetaBridge = bridge
}

// WithBlockCache attaches a new block cache to the chain client
func (ob *ChainClient) WithBlockCache(cache *lru.Cache) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.blockCache = cache
}

// SetChainParams sets the chain params for the chain client
func (ob *ChainClient) SetChainParams(params observertypes.ChainParams) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.chainParams = params
}

// GetChainParams returns the chain params for the chain client
func (ob *ChainClient) GetChainParams() observertypes.ChainParams {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.chainParams
}

func (ob *ChainClient) GetConnectorContract() (ethcommon.Address, *zetaconnector.ZetaConnectorNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := FetchConnectorContract(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *ChainClient) GetConnectorContractEth() (ethcommon.Address, *zetaconnectoreth.ZetaConnectorEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := FetchConnectorContractEth(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *ChainClient) GetZetaTokenNonEthContract() (ethcommon.Address, *zeta.ZetaNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ZetaTokenContractAddress)
	contract, err := FetchZetaZetaNonEthTokenContract(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *ChainClient) GetERC20CustodyContract() (ethcommon.Address, *erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().Erc20CustodyContractAddress)
	contract, err := FetchERC20CustodyContract(addr, ob.evmClient)
	return addr, contract, err
}

func FetchConnectorContract(addr ethcommon.Address, client interfaces.EVMRPCClient) (*zetaconnector.ZetaConnectorNonEth, error) {
	return zetaconnector.NewZetaConnectorNonEth(addr, client)
}

func FetchConnectorContractEth(addr ethcommon.Address, client interfaces.EVMRPCClient) (*zetaconnectoreth.ZetaConnectorEth, error) {
	return zetaconnectoreth.NewZetaConnectorEth(addr, client)
}

func FetchZetaZetaNonEthTokenContract(addr ethcommon.Address, client interfaces.EVMRPCClient) (*zeta.ZetaNonEth, error) {
	return zeta.NewZetaNonEth(addr, client)
}

func FetchERC20CustodyContract(addr ethcommon.Address, client interfaces.EVMRPCClient) (*erc20custody.ERC20Custody, error) {
	return erc20custody.NewERC20Custody(addr, client)
}

// Start all observation routines for the evm chain
func (ob *ChainClient) Start() {
	// watch evm chain for incoming txs and post votes to zetacore
	go ob.WatchInbound()

	// watch evm chain for outgoing txs status
	go ob.WatchOutbound()

	// watch evm chain for gas prices and post to zetacore
	go ob.WatchGasPrice()

	// watch zetacore for intx trackers
	go ob.WatchInboundTracker()

	// watch the RPC status of the evm chain
	go ob.WatchRPCStatus()
}

// WatchRPCStatus watches the RPC status of the evm chain
func (ob *ChainClient) WatchRPCStatus() {
	ob.logger.Chain.Info().Msgf("Starting RPC status check for chain %s", ob.chain.String())
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}
			bn, err := ob.evmClient.BlockNumber(context.Background())
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			gasPrice, err := ob.evmClient.SuggestGasPrice(context.Background())
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			header, err := ob.evmClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(bn))
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			// #nosec G701 always in range
			blockTime := time.Unix(int64(header.Time), 0).UTC()
			elapsedSeconds := time.Since(blockTime).Seconds()
			if elapsedSeconds > 100 {
				ob.logger.Chain.Warn().Msgf("RPC Status Check warning: RPC stale or chain stuck (check explorer)? Latest block %d timestamp is %.0fs ago", bn, elapsedSeconds)
				continue
			}
			ob.logger.Chain.Info().Msgf("[OK] RPC status: latest block num %d, timestamp %s ( %.0fs ago), suggested gas price %d", header.Number, blockTime.String(), elapsedSeconds, gasPrice.Uint64())
		case <-ob.stop:
			return
		}
	}
}

func (ob *ChainClient) Stop() {
	ob.logger.Chain.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.Chain.Info().Msg("closing ob.db")
	dbInst, err := ob.db.DB()
	if err != nil {
		ob.logger.Chain.Info().Msg("error getting database instance")
	}
	err = dbInst.Close()
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("error closing database")
	}

	ob.logger.Chain.Info().Msgf("%s observer stopped", ob.chain.String())
}

// WatchOutbound watches evm chain for outgoing txs status
func (ob *ChainClient) WatchOutbound() {
	ticker, err := clienttypes.NewDynamicTicker(fmt.Sprintf("EVM_WatchOutbound_%d", ob.chain.ChainId), ob.GetChainParams().OutboundTicker)
	if err != nil {
		ob.logger.Outbound.Error().Err(err).Msg("error creating ticker")
		return
	}

	ob.logger.Outbound.Info().Msgf("WatchOutbound started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.Outbound.Sample(&zerolog.BasicSampler{N: 10})
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !corecontext.IsOutboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchOutbound: outbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			trackers, err := ob.zetaBridge.GetAllOutboundTrackerByChainbound(ob.chain.ChainId, interfaces.Ascending)
			if err != nil {
				continue
			}
			for _, tracker := range trackers {
				nonceInt := tracker.Nonce
				if ob.IsTxConfirmed(nonceInt) { // Go to next tracker if this one already has a confirmed tx
					continue
				}
				txCount := 0
				var outtxReceipt *ethtypes.Receipt
				var outtx *ethtypes.Transaction
				for _, txHash := range tracker.HashList {
					if receipt, tx, ok := ob.checkConfirmedTx(txHash.TxHash, nonceInt); ok {
						txCount++
						outtxReceipt = receipt
						outtx = tx
						ob.logger.Outbound.Info().Msgf("WatchOutbound: confirmed outTx %s for chain %d nonce %d", txHash.TxHash, ob.chain.ChainId, nonceInt)
						if txCount > 1 {
							ob.logger.Outbound.Error().Msgf(
								"WatchOutbound: checkConfirmedTx passed, txCount %d chain %d nonce %d receipt %v transaction %v", txCount, ob.chain.ChainId, nonceInt, outtxReceipt, outtx)
						}
					}
				}
				if txCount == 1 { // should be only one txHash confirmed for each nonce.
					ob.SetTxNReceipt(nonceInt, outtxReceipt, outtx)
				} else if txCount > 1 { // should not happen. We can't tell which txHash is true. It might happen (e.g. glitchy/hacked endpoint)
					ob.logger.Outbound.Error().Msgf("WatchOutbound: confirmed multiple (%d) outTx for chain %d nonce %d", txCount, ob.chain.ChainId, nonceInt)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutboundTicker, ob.logger.Outbound)
		case <-ob.stop:
			ob.logger.Outbound.Info().Msg("WatchOutbound: stopped")
			return
		}
	}
}

// SetPendingTx sets the pending transaction in memory
func (ob *ChainClient) SetPendingTx(nonce uint64, transaction *ethtypes.Transaction) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.outboundPendingTransactions[ob.GetTxID(nonce)] = transaction
}

// GetPendingTx gets the pending transaction from memory
func (ob *ChainClient) GetPendingTx(nonce uint64) *ethtypes.Transaction {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.outboundPendingTransactions[ob.GetTxID(nonce)]
}

// SetTxNReceipt sets the receipt and transaction in memory
func (ob *ChainClient) SetTxNReceipt(nonce uint64, receipt *ethtypes.Receipt, transaction *ethtypes.Transaction) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	delete(ob.outboundPendingTransactions, ob.GetTxID(nonce)) // remove pending transaction, if any
	ob.outboundConfirmedReceipts[ob.GetTxID(nonce)] = receipt
	ob.outboundConfirmedTransactions[ob.GetTxID(nonce)] = transaction
}

// GetTxNReceipt gets the receipt and transaction from memory
func (ob *ChainClient) GetTxNReceipt(nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	receipt := ob.outboundConfirmedReceipts[ob.GetTxID(nonce)]
	transaction := ob.outboundConfirmedTransactions[ob.GetTxID(nonce)]
	return receipt, transaction
}

// IsTxConfirmed returns true if there is a confirmed tx for 'nonce'
func (ob *ChainClient) IsTxConfirmed(nonce uint64) bool {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.outboundConfirmedReceipts[ob.GetTxID(nonce)] != nil && ob.outboundConfirmedTransactions[ob.GetTxID(nonce)] != nil
}

// CheckTxInclusion returns nil only if tx is included at the position indicated by the receipt ([block, index])
func (ob *ChainClient) CheckTxInclusion(tx *ethtypes.Transaction, receipt *ethtypes.Receipt) error {
	block, err := ob.GetBlockByNumberCached(receipt.BlockNumber.Uint64())
	if err != nil {
		return errors.Wrapf(err, "GetBlockByNumberCached error for block %d txHash %s nonce %d",
			receipt.BlockNumber.Uint64(), tx.Hash(), tx.Nonce())
	}

	// #nosec G701 non negative value
	if receipt.TransactionIndex >= uint(len(block.Transactions)) {
		return fmt.Errorf("transaction index %d out of range [0, %d), txHash %s nonce %d block %d",
			receipt.TransactionIndex, len(block.Transactions), tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	txAtIndex := block.Transactions[receipt.TransactionIndex]
	if !strings.EqualFold(txAtIndex.Hash, tx.Hash().Hex()) {
		ob.RemoveCachedBlock(receipt.BlockNumber.Uint64()) // clean stale block from cache
		return fmt.Errorf("transaction at index %d has different hash %s, txHash %s nonce %d block %d",
			receipt.TransactionIndex, txAtIndex.Hash, tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	return nil
}

// SetLastBlockHeightScanned set last block height scanned (not necessarily caught up with external block; could be slow/paused)
func (ob *ChainClient) SetLastBlockHeightScanned(height uint64) {
	atomic.StoreUint64(&ob.lastBlockScanned, height)
	metrics.LastScannedBlockNumber.WithLabelValues(ob.chain.ChainName.String()).Set(float64(height))
}

// GetLastBlockHeightScanned get last block height scanned (not necessarily caught up with external block; could be slow/paused)
func (ob *ChainClient) GetLastBlockHeightScanned() uint64 {
	height := atomic.LoadUint64(&ob.lastBlockScanned)
	return height
}

// SetLastBlockHeight set external last block height
func (ob *ChainClient) SetLastBlockHeight(height uint64) {
	if height >= math.MaxInt64 {
		panic("lastBlock is too large")
	}
	atomic.StoreUint64(&ob.lastBlock, height)
}

// GetLastBlockHeight get external last block height
func (ob *ChainClient) GetLastBlockHeight() uint64 {
	height := atomic.LoadUint64(&ob.lastBlock)
	if height >= math.MaxInt64 {
		panic("lastBlock is too large")
	}
	return height
}

// WatchInbound watches evm chain for incoming txs and post votes to zetacore
func (ob *ChainClient) WatchInbound() {
	ticker, err := clienttypes.NewDynamicTicker(fmt.Sprintf("EVM_WatchInbound_%d", ob.chain.ChainId), ob.GetChainParams().InboundTicker)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	ob.logger.Inbound.Info().Msgf("WatchInbound started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.Inbound.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !corecontext.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchInbound: inbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			err := ob.observeInbound(sampledLogger)
			if err != nil {
				ob.logger.Inbound.Err(err).Msg("WatchInbound: observeInbound error")
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.stop:
			ob.logger.Inbound.Info().Msgf("WatchInbound stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

// ObserveZetaSent queries the ZetaSent event from the connector contract and posts to zetabridge
// returns the last block successfully scanned
func (ob *ChainClient) ObserveZetaSent(startBlock, toBlock uint64) uint64 {
	// filter ZetaSent logs
	addrConnector, connector, err := ob.GetConnectorContract()
	if err != nil {
		ob.logger.Chain.Warn().Err(err).Msgf("ObserveZetaSent: GetConnectorContract error:")
		return startBlock - 1 // lastScanned
	}
	iter, err := connector.FilterZetaSent(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{}, []*big.Int{})
	if err != nil {
		ob.logger.Chain.Warn().Err(err).Msgf(
			"ObserveZetaSent: FilterZetaSent error from block %d to %d for chain %d", startBlock, toBlock, ob.chain.ChainId)
		return startBlock - 1 // lastScanned
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*zetaconnector.ZetaConnectorNonEthZetaSent, 0)
	for iter.Next() {
		// sanity check tx event
		err := ValidateEvmTxLog(&iter.Event.Raw, addrConnector, "", TopicsZetaSent)
		if err == nil {
			events = append(events, iter.Event)
			continue
		}
		ob.logger.Inbound.Warn().Err(err).Msgf("ObserveZetaSent: invalid ZetaSent event in tx %s on chain %d at height %d",
			iter.Event.Raw.TxHash.Hex(), ob.chain.ChainId, iter.Event.Raw.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.chain.ChainName.String()).Inc()

	// post to zetabridge
	beingScanned := uint64(0)
	guard := make(map[string]bool)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.logger.Inbound.Warn().Msgf("ObserveZetaSent: multiple remote call events detected in tx %s", event.Raw.TxHash)
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true

		msg := ob.BuildInboundVoteMsgForZetaSentEvent(event)
		if msg != nil {
			_, err = ob.PostVoteInbound(msg, coin.CoinType_Zeta, zetabridge.PostVoteInboundMessagePassingExecutionGasLimit)
			if err != nil {
				return beingScanned - 1 // we have to re-scan from this block next time
			}
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock
}

// ObserveERC20Deposited queries the ERC20CustodyDeposited event from the ERC20Custody contract and posts to zetabridge
// returns the last block successfully scanned
func (ob *ChainClient) ObserveERC20Deposited(startBlock, toBlock uint64) uint64 {
	// filter ERC20CustodyDeposited logs
	addrCustody, erc20custodyContract, err := ob.GetERC20CustodyContract()
	if err != nil {
		ob.logger.Inbound.Warn().Err(err).Msgf("ObserveERC20Deposited: GetERC20CustodyContract error:")
		return startBlock - 1 // lastScanned
	}

	iter, err := erc20custodyContract.FilterDeposited(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{})
	if err != nil {
		ob.logger.Inbound.Warn().Err(err).Msgf(
			"ObserveERC20Deposited: FilterDeposited error from block %d to %d for chain %d", startBlock, toBlock, ob.chain.ChainId)
		return startBlock - 1 // lastScanned
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*erc20custody.ERC20CustodyDeposited, 0)
	for iter.Next() {
		// sanity check tx event
		err := ValidateEvmTxLog(&iter.Event.Raw, addrCustody, "", TopicsDeposited)
		if err == nil {
			events = append(events, iter.Event)
			continue
		}
		ob.logger.Inbound.Warn().Err(err).Msgf("ObserveERC20Deposited: invalid Deposited event in tx %s on chain %d at height %d",
			iter.Event.Raw.TxHash.Hex(), ob.chain.ChainId, iter.Event.Raw.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.chain.ChainName.String()).Inc()

	// post to zetabridge
	guard := make(map[string]bool)
	beingScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		tx, _, err := ob.TransactionByHash(event.Raw.TxHash.Hex())
		if err != nil {
			ob.logger.Inbound.Error().Err(err).Msgf(
				"ObserveERC20Deposited: error getting transaction for intx %s chain %d", event.Raw.TxHash, ob.chain.ChainId)
			return beingScanned - 1 // we have to re-scan from this block next time
		}
		sender := ethcommon.HexToAddress(tx.From)

		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.logger.Inbound.Warn().Msgf("ObserveERC20Deposited: multiple remote call events detected in tx %s", event.Raw.TxHash)
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true

		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		if msg != nil {
			_, err = ob.PostVoteInbound(msg, coin.CoinType_ERC20, zetabridge.PostVoteInboundExecutionGasLimit)
			if err != nil {
				return beingScanned - 1 // we have to re-scan from this block next time
			}
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock
}

// ObserverTSSReceive queries the incoming gas asset to TSS address and posts to zetabridge
// returns the last block successfully scanned
func (ob *ChainClient) ObserverTSSReceive(startBlock, toBlock uint64) uint64 {
	// query incoming gas asset
	for bn := startBlock; bn <= toBlock; bn++ {
		// post new block header (if any) to zetabridge and ignore error
		// TODO: consider having a independent ticker(from TSS scaning) for posting block headers
		// https://github.com/zeta-chain/node/issues/1847
		blockHeaderVerification, found := ob.coreContext.GetBlockHeaderEnabledChains(ob.chain.ChainId)
		if found && blockHeaderVerification.Enabled {
			// post block header for supported chains
			err := ob.postBlockHeader(toBlock)
			if err != nil {
				ob.logger.Inbound.Error().Err(err).Msg("error posting block header")
			}
		}

		// observe TSS received gas token in block 'bn'
		err := ob.ObserveTSSReceiveInBlock(bn)
		if err != nil {
			ob.logger.Inbound.Error().Err(err).Msgf("ObserverTSSReceive: error observing TSS received token in block %d for chain %d", bn, ob.chain.ChainId)
			return bn - 1 // we have to re-scan from this block next time
		}
	}
	// successful processed all gas asset deposits in [startBlock, toBlock]
	return toBlock
}

// WatchGasPrice watches evm chain for gas prices and post to zetacore
func (ob *ChainClient) WatchGasPrice() {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice()
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.chain.ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker(fmt.Sprintf("EVM_WatchGasPrice_%d", ob.chain.ChainId), ob.GetChainParams().GasPriceTicker)
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msg("NewDynamicTicker error")
		return
	}
	ob.logger.GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.chain.ChainId, ob.GetChainParams().GasPriceTicker)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err = ob.PostGasPrice()
			if err != nil {
				ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.chain.ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.logger.GasPrice)
		case <-ob.stop:
			ob.logger.GasPrice.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *ChainClient) PostGasPrice() error {

	// GAS PRICE
	gasPrice, err := ob.evmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		ob.logger.GasPrice.Err(err).Msg("Err SuggestGasPrice:")
		return err
	}
	blockNum, err := ob.evmClient.BlockNumber(context.TODO())
	if err != nil {
		ob.logger.GasPrice.Err(err).Msg("Err Fetching Most recent Block : ")
		return err
	}

	// SUPPLY
	supply := "100" // lockedAmount on ETH, totalSupply on other chains

	zetaHash, err := ob.zetaBridge.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		ob.logger.GasPrice.Err(err).Msg("PostGasPrice to zetabridge failed")
		return err
	}
	_ = zetaHash

	return nil
}

func (ob *ChainClient) BuildLastBlock() error {
	logger := ob.logger.Chain.With().Str("module", "BuildBlockIndex").Logger()
	envvar := ob.chain.ChainName.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		logger.Info().Msgf("BuildLastBlock: envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == clienttypes.EnvVarLatest {
			header, err := ob.evmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(header.Number.Uint64())
		} else {
			scanFromBlockInt, err := strconv.ParseUint(scanFromBlock, 10, 64)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(scanFromBlockInt)
		}
	} else { // last observed block
		var lastBlockNum clienttypes.LastBlockSQLType
		if err := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumID).Error; err != nil {
			logger.Info().Msgf("BuildLastBlock: db PosKey does not exist; read from external chain %s", ob.chain.String())
			header, err := ob.evmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(header.Number.Uint64())
			if dbc := ob.db.Save(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeightScanned())); dbc.Error != nil {
				logger.Error().Err(dbc.Error).Msgf("BuildLastBlock: error writing lastBlockScanned %d to db", ob.GetLastBlockHeightScanned())
			}
		} else {
			ob.SetLastBlockHeightScanned(lastBlockNum.Num)
		}
	}
	return nil
}

// LoadDB open sql database and load data into EVMChainClient
func (ob *ChainClient) LoadDB(dbPath string, chain chains.Chain) error {
	if dbPath != "" {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			err := os.MkdirAll(dbPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		path := fmt.Sprintf("%s/%s", dbPath, chain.ChainName.String()) //Use "file::memory:?cache=shared" for temp db
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		err = db.AutoMigrate(&clienttypes.ReceiptSQLType{},
			&clienttypes.TransactionSQLType{},
			&clienttypes.LastBlockSQLType{})
		if err != nil {
			ob.logger.Chain.Error().Err(err).Msg("error migrating db")
			return err
		}

		ob.db = db
		err = ob.BuildLastBlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ob *ChainClient) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.EVMAddress().String()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

// BlockByNumber query block by number via JSON-RPC
func (ob *ChainClient) BlockByNumber(blockNumber int) (*ethrpc.Block, error) {
	block, err := ob.evmJSONRPC.EthGetBlockByNumber(blockNumber, true)
	if err != nil {
		return nil, err
	}
	for i := range block.Transactions {
		err := ValidateEvmTransaction(&block.Transactions[i])
		if err != nil {
			return nil, err
		}
	}
	return block, nil
}

// TransactionByHash query transaction by hash via JSON-RPC
func (ob *ChainClient) TransactionByHash(txHash string) (*ethrpc.Transaction, bool, error) {
	tx, err := ob.evmJSONRPC.EthGetTransactionByHash(txHash)
	if err != nil {
		return nil, false, err
	}
	err = ValidateEvmTransaction(tx)
	if err != nil {
		return nil, false, err
	}
	return tx, tx.BlockNumber == nil, nil
}

func (ob *ChainClient) GetBlockHeaderCached(blockNumber uint64) (*ethtypes.Header, error) {
	if header, ok := ob.headerCache.Get(blockNumber); ok {
		return header.(*ethtypes.Header), nil
	}
	header, err := ob.evmClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.headerCache.Add(blockNumber, header)
	return header, nil
}

// GetBlockByNumberCached get block by number from cache
// returns block, ethrpc.Block, isFallback, isSkip, error
func (ob *ChainClient) GetBlockByNumberCached(blockNumber uint64) (*ethrpc.Block, error) {
	if block, ok := ob.blockCache.Get(blockNumber); ok {
		return block.(*ethrpc.Block), nil
	}
	if blockNumber > math.MaxInt32 {
		return nil, fmt.Errorf("block number %d is too large", blockNumber)
	}
	// #nosec G701 always in range, checked above
	block, err := ob.BlockByNumber(int(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.blockCache.Add(blockNumber, block)
	return block, nil
}

// RemoveCachedBlock remove block from cache
func (ob *ChainClient) RemoveCachedBlock(blockNumber uint64) {
	ob.blockCache.Remove(blockNumber)
}

// checkConfirmedTx checks if a txHash is confirmed
// returns (receipt, transaction, true) if confirmed or (nil, nil, false) otherwise
func (ob *ChainClient) checkConfirmedTx(txHash string, nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction, bool) {
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query transaction
	transaction, isPending, err := ob.evmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: error getting transaction for outtx %s chain %d", txHash, ob.chain.ChainId)
		return nil, nil, false
	}
	if transaction == nil { // should not happen
		log.Error().Msgf("confirmTxByHash: transaction is nil for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	// check tx sender and nonce
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
	from, err := signer.Sender(transaction)
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: local recovery of sender address failed for outtx %s chain %d", transaction.Hash().Hex(), ob.chain.ChainId)
		return nil, nil, false
	}
	if from != ob.Tss.EVMAddress() { // must be TSS address
		log.Error().Msgf("confirmTxByHash: sender %s for outtx %s chain %d is not TSS address %s",
			from.Hex(), transaction.Hash().Hex(), ob.chain.ChainId, ob.Tss.EVMAddress().Hex())
		return nil, nil, false
	}
	if transaction.Nonce() != nonce { // must match cctx nonce
		log.Error().Msgf("confirmTxByHash: outtx %s nonce mismatch: wanted %d, got tx nonce %d", txHash, nonce, transaction.Nonce())
		return nil, nil, false
	}

	// save pending transaction
	if isPending {
		ob.SetPendingTx(nonce, transaction)
		return nil, nil, false
	}

	// query receipt
	receipt, err := ob.evmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			log.Warn().Err(err).Msgf("confirmTxByHash: TransactionReceipt error, txHash %s nonce %d", txHash, nonce)
		}
		return nil, nil, false
	}
	if receipt == nil { // should not happen
		log.Error().Msgf("confirmTxByHash: receipt is nil for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	// check confirmations
	if !ob.HasEnoughConfirmations(receipt, ob.GetLastBlockHeight()) {
		log.Debug().Msgf("confirmTxByHash: txHash %s nonce %d included but not confirmed: receipt block %d, current block %d",
			txHash, nonce, receipt.BlockNumber, ob.GetLastBlockHeight())
		return nil, nil, false
	}

	// cross-check tx inclusion against the block
	// Note: a guard for false BlockNumber in receipt. The blob-carrying tx won't come here
	err = ob.CheckTxInclusion(transaction, receipt)
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: checkTxInclusion error for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	return receipt, transaction, true
}

// calcBlockRangeToScan calculates the next range of blocks to scan
func (ob *ChainClient) calcBlockRangeToScan(latestConfirmed, lastScanned, batchSize uint64) (uint64, uint64) {
	startBlock := lastScanned + 1
	toBlock := lastScanned + batchSize
	if toBlock > latestConfirmed {
		toBlock = latestConfirmed
	}
	return startBlock, toBlock
}

func (ob *ChainClient) postBlockHeader(tip uint64) error {
	bn := tip

	res, err := ob.zetaBridge.GetBlockHeaderChainState(ob.chain.ChainId)
	if err == nil && res.ChainState != nil && res.ChainState.EarliestHeight > 0 {
		// #nosec G701 always positive
		bn = uint64(res.ChainState.LatestHeight) + 1 // the next header to post
	}

	if bn > tip {
		return fmt.Errorf("postBlockHeader: must post block confirmed block header: %d > %d", bn, tip)
	}

	header, err := ob.GetBlockHeaderCached(bn)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("postBlockHeader: error getting block: %d", bn)
		return err
	}
	headerRLP, err := rlp.EncodeToBytes(header)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("postBlockHeader: error encoding block header: %d", bn)
		return err
	}

	_, err = ob.zetaBridge.PostVoteBlockHeader(
		ob.chain.ChainId,
		header.Hash().Bytes(),
		header.Number.Int64(),
		proofs.NewEthereumHeader(headerRLP),
	)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("postBlockHeader: error posting block header: %d", bn)
		return err
	}
	return nil
}

func (ob *ChainClient) observeInbound(sampledLogger zerolog.Logger) error {
	// get and update latest block height
	blockNumber, err := ob.evmClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	if blockNumber < ob.GetLastBlockHeight() {
		return fmt.Errorf("observeInbound: block number should not decrease: current %d last %d", blockNumber, ob.GetLastBlockHeight())
	}
	ob.SetLastBlockHeight(blockNumber)

	// increment prom counter
	metrics.GetBlockByNumberPerChain.WithLabelValues(ob.chain.ChainName.String()).Inc()

	// skip if current height is too low
	if blockNumber < ob.GetChainParams().ConfirmationCount {
		return fmt.Errorf("observeInbound: skipping observer, current block number %d is too low", blockNumber)
	}
	confirmedBlockNum := blockNumber - ob.GetChainParams().ConfirmationCount

	// skip if no new block is confirmed
	lastScanned := ob.GetLastBlockHeightScanned()
	if lastScanned >= confirmedBlockNum {
		sampledLogger.Debug().Msgf("observeInbound: skipping observer, no new block is produced for chain %d", ob.chain.ChainId)
		return nil
	}

	// get last scanned block height (we simply use same height for all 3 events ZetaSent, Deposited, TssRecvd)
	// Note: using different heights for each event incurs more complexity (metrics, db, etc) and not worth it
	startBlock, toBlock := ob.calcBlockRangeToScan(confirmedBlockNum, lastScanned, config.MaxBlocksPerPeriod)

	// task 1:  query evm chain for zeta sent logs (read at most 100 blocks in one go)
	lastScannedZetaSent := ob.ObserveZetaSent(startBlock, toBlock)

	// task 2: query evm chain for deposited logs (read at most 100 blocks in one go)
	lastScannedDeposited := ob.ObserveERC20Deposited(startBlock, toBlock)

	// task 3: query the incoming tx to TSS address (read at most 100 blocks in one go)
	lastScannedTssRecvd := ob.ObserverTSSReceive(startBlock, toBlock)

	// note: using lowest height for all 3 events is not perfect, but it's simple and good enough
	lastScannedLowest := lastScannedZetaSent
	if lastScannedDeposited < lastScannedLowest {
		lastScannedLowest = lastScannedDeposited
	}
	if lastScannedTssRecvd < lastScannedLowest {
		lastScannedLowest = lastScannedTssRecvd
	}

	// update last scanned block height for all 3 events (ZetaSent, Deposited, TssRecvd), ignore db error
	if lastScannedLowest > lastScanned {
		sampledLogger.Info().Msgf("observeInbound: lasstScanned heights for chain %d ZetaSent %d ERC20Deposited %d TssRecvd %d",
			ob.chain.ChainId, lastScannedZetaSent, lastScannedDeposited, lastScannedTssRecvd)
		ob.SetLastBlockHeightScanned(lastScannedLowest)
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(lastScannedLowest)).Error; err != nil {
			ob.logger.Inbound.Error().Err(err).Msgf("observeInbound: error writing lastScannedLowest %d to db", lastScannedLowest)
		}
	}
	return nil
}
