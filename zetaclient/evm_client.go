package zetaclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	math2 "math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	erc20custody "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

type TxHashEnvelope struct {
	TxHash string
	Done   chan struct{}
}

type OutTx struct {
	SendHash string
	TxHash   string
	Nonce    int64
}
type EVMLog struct {
	ChainLogger          zerolog.Logger // Parent logger
	ExternalChainWatcher zerolog.Logger // Observes external Chains for incoming trasnactions
	WatchGasPrice        zerolog.Logger // Observes external Chains for Gas prices and posts to core
	ObserveOutTx         zerolog.Logger // Observes external Chains for Outgoing transactions

}

const (
	DonationMessage = "I am rich!"
)

// Chain configuration struct
// Filled with above constants depending on chain
type EVMChainClient struct {
	*ChainMetrics
	chain                     common.Chain
	EvmClient                 *ethclient.Client
	KlaytnClient              *KlaytnClient
	zetaClient                *ZetaCoreBridge
	Tss                       TSSSigner
	lastBlockScanned          int64
	lastBlock                 int64
	BlockTimeExternalChain    uint64 // block time in seconds
	txWatchList               map[ethcommon.Hash]string
	mu                        *sync.Mutex
	db                        *gorm.DB
	outTXConfirmedReceipts    map[string]*ethtypes.Receipt
	outTXConfirmedTransaction map[string]*ethtypes.Transaction
	MinNonce                  int64
	MaxNonce                  int64
	OutTxChan                 chan OutTx // send to this channel if you want something back!
	stop                      chan struct{}
	fileLogger                *zerolog.Logger // for critical info
	logger                    EVMLog
	cfg                       *config.Config
	params                    observertypes.CoreParams
	ts                        *TelemetryServer

	BlockCache *lru.Cache
}

var _ ChainClient = (*EVMChainClient)(nil)

// Return configuration based on supplied target chain
func NewEVMChainClient(bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics, logger zerolog.Logger, cfg *config.Config, evmCfg config.EVMConfig, ts *TelemetryServer) (*EVMChainClient, error) {
	ob := EVMChainClient{
		ChainMetrics: NewChainMetrics(evmCfg.Chain.ChainName.String(), metrics),
		ts:           ts,
	}
	chainLogger := logger.With().Str("chain", evmCfg.Chain.ChainName.String()).Logger()
	ob.logger = EVMLog{
		ChainLogger:          chainLogger,
		ExternalChainWatcher: chainLogger.With().Str("module", "ExternalChainWatcher").Logger(),
		WatchGasPrice:        chainLogger.With().Str("module", "WatchGasPrice").Logger(),
		ObserveOutTx:         chainLogger.With().Str("module", "ObserveOutTx").Logger(),
	}
	ob.cfg = cfg
	ob.params = evmCfg.CoreParams
	ob.stop = make(chan struct{})
	ob.chain = evmCfg.Chain
	ob.mu = &sync.Mutex{}
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.outTXConfirmedReceipts = make(map[string]*ethtypes.Receipt)
	ob.outTXConfirmedTransaction = make(map[string]*ethtypes.Transaction)
	ob.OutTxChan = make(chan OutTx, 100)

	logFile, err := os.OpenFile(ob.chain.ChainName.String()+"_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error creating a logFile chain %s", ob.chain.ChainName.String())
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	ob.fileLogger = &fileLogger

	ob.logger.ChainLogger.Info().Msgf("Chain %s endpoint %s", ob.chain.ChainName.String(), evmCfg.Endpoint)
	client, err := ethclient.Dial(evmCfg.Endpoint)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client

	ob.BlockCache, err = lru.New(1000)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("failed to create block cache")
		return nil, err
	}

	if ob.chain.IsKlaytnChain() {
		kclient, err := Dial(evmCfg.Endpoint)
		if err != nil {
			ob.logger.ChainLogger.Err(err).Msg("klaytn Client Dial")
			return nil, err
		}
		ob.KlaytnClient = kclient
	}

	// create metric counters
	err = ob.RegisterPromCounter("rpc_getLogs_count", "Number of getLogs")
	if err != nil {
		return nil, err
	}
	err = ob.RegisterPromCounter("rpc_getBlockByNumber_count", "Number of getBlockByNumber")
	if err != nil {
		return nil, err
	}
	err = ob.RegisterPromGauge(metricsPkg.PendingTxs, "Number of pending transactions")
	if err != nil {
		return nil, err
	}

	err = ob.LoadDB(dbpath, ob.chain)
	if err != nil {
		return nil, err
	}

	ob.logger.ChainLogger.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.GetLastBlockHeightScanned())

	return &ob, nil
}

func (ob *EVMChainClient) SetCoreParams(params observertypes.CoreParams) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.params = params
}

func (ob *EVMChainClient) GetCoreParams() observertypes.CoreParams {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.params
}

func (ob *EVMChainClient) GetConnectorContract() (*zetaconnector.ZetaConnectorNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetCoreParams().ConnectorContractAddress)
	return zetaconnector.NewZetaConnectorNonEth(addr, ob.EvmClient)
}

func (ob *EVMChainClient) GetERC20CustodyContract() (*erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.GetCoreParams().Erc20CustodyContractAddress)
	return erc20custody.NewERC20Custody(addr, ob.EvmClient)
}

func (ob *EVMChainClient) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.observeOutTx()         // Populates receipts and confirmed outbound transactions
}

func (ob *EVMChainClient) Stop() {
	ob.logger.ChainLogger.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.ChainLogger.Info().Msg("closing ob.db")
	dbInst, err := ob.db.DB()
	if err != nil {
		ob.logger.ChainLogger.Info().Msg("error getting database instance")
	}
	err = dbInst.Close()
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("error closing database")
	}

	ob.logger.ChainLogger.Info().Msgf("%s observer stopped", ob.chain.String())
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *EVMChainClient) IsSendOutTxProcessed(sendHash string, nonce uint64, cointype common.CoinType, logger zerolog.Logger) (bool, bool, error) {
	ob.mu.Lock()
	params := ob.params
	receipt, found1 := ob.outTXConfirmedReceipts[ob.GetTxID(nonce)]
	transaction, found2 := ob.outTXConfirmedTransaction[ob.GetTxID(nonce)]
	ob.mu.Unlock()
	found := found1 && found2
	if !found {
		return false, false, nil
	}

	sendID := fmt.Sprintf("%s-%d", ob.chain.String(), nonce)
	logger = logger.With().Str("sendID", sendID).Logger()
	if cointype == common.CoinType_Cmd {
		recvStatus := common.ReceiveStatus_Failed
		if receipt.Status == 1 {
			recvStatus = common.ReceiveStatus_Success
		}
		zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
			sendHash,
			receipt.TxHash.Hex(),
			receipt.BlockNumber.Uint64(),
			receipt.GasUsed,
			transaction.GasPrice(),
			transaction.Gas(),
			transaction.Value(),
			recvStatus,
			ob.chain,
			nonce,
			common.CoinType_Cmd,
		)
		if err != nil {
			logger.Error().Err(err).Msg("error posting confirmation to meta core")
		}
		logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
		return true, true, nil

	} else if cointype == common.CoinType_Gas { // the outbound is a regular Ether/BNB/Matic transfer; no need to check events
		if receipt.Status == 1 {
			zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				transaction.Value(),
				common.ReceiveStatus_Success,
				ob.chain,
				nonce,
				common.CoinType_Gas,
			)
			if err != nil {
				logger.Error().Err(err).Msg("error posting confirmation to meta core")
			}
			logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
			return true, true, nil
		} else if receipt.Status == 0 { // the same as below events flow
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				big.NewInt(0),
				common.ReceiveStatus_Failed,
				ob.chain,
				nonce,
				common.CoinType_Gas,
			)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	} else if cointype == common.CoinType_Zeta { // the outbound is a Zeta transfer; need to check events ZetaReceived
		if receipt.Status == 1 {
			logs := receipt.Logs
			for _, vLog := range logs {
				confHeight := vLog.BlockNumber + params.ConfirmationCount
				if confHeight < 0 || confHeight >= math2.MaxInt64 {
					return false, false, fmt.Errorf("confHeight is out of range")
				}
				// TODO rewrite this to return early if not confirmed
				connector, err := ob.GetConnectorContract()
				if err != nil {
					return false, false, fmt.Errorf("error getting connector contract: %w", err)
				}
				receivedLog, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if int64(confHeight) < ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						if len(vLog.Topics) != 4 {
							logger.Error().Msgf("wrong number of topics in log %d", len(vLog.Topics))
							return false, false, fmt.Errorf("wrong number of topics in log %d", len(vLog.Topics))
						}
						sendhash := vLog.Topics[3].Hex()
						//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
						mMint := receivedLog.ZetaValue
						zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							receipt.GasUsed,
							transaction.GasPrice(),
							transaction.Gas(),
							mMint,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_Zeta,
						)
						if err != nil {
							logger.Error().Err(err).Msg("error posting confirmation to meta core")
							continue
						}
						logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+params.ConfirmationCount)-int(ob.GetLastBlockHeight()), ob.chain.String(), nonce)
					return true, false, nil
				}
				revertedLog, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReverted(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if int64(confHeight) < ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						if len(vLog.Topics) != 3 {
							logger.Error().Msgf("wrong number of topics in log %d", len(vLog.Topics))
							return false, false, fmt.Errorf("wrong number of topics in log %d", len(vLog.Topics))
						}
						sendhash := vLog.Topics[2].Hex()
						mMint := revertedLog.RemainingZetaValue
						metaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							receipt.GasUsed,
							transaction.GasPrice(),
							transaction.Gas(),
							mMint,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_Zeta,
						)
						if err != nil {
							logger.Err(err).Msg("error posting confirmation to meta core")
							continue
						}
						logger.Info().Msgf("Zeta tx hash: %s", metaHash)
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+params.ConfirmationCount)-int(ob.GetLastBlockHeight()), ob.chain.String(), nonce)
					return true, false, nil
				}
			}
		} else if receipt.Status == 0 {
			//FIXME: check nonce here by getTransaction RPC
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				big.NewInt(0),
				common.ReceiveStatus_Failed,
				ob.chain,
				nonce,
				common.CoinType_Zeta,
			)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	} else if cointype == common.CoinType_ERC20 {
		if receipt.Status == 1 {
			logs := receipt.Logs
			ERC20Custody, err := ob.GetERC20CustodyContract()
			if err != nil {
				logger.Warn().Msgf("NewERC20Custody err: %s", err)
			}
			for _, vLog := range logs {
				event, err := ERC20Custody.ParseWithdrawn(*vLog)
				confHeight := vLog.BlockNumber + params.ConfirmationCount
				if confHeight < 0 || confHeight >= math2.MaxInt64 {
					return false, false, fmt.Errorf("confHeight is out of range")
				}
				if err == nil {
					logger.Info().Msgf("Found (ERC20Custody.Withdrawn Event) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if int64(confHeight) < ob.GetLastBlockHeight() {

						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendHash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							receipt.GasUsed,
							transaction.GasPrice(),
							transaction.Gas(),
							event.Amount,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_ERC20,
						)
						if err != nil {
							logger.Error().Err(err).Msg("error posting confirmation to meta core")
							continue
						}
						logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+params.ConfirmationCount)-int(ob.GetLastBlockHeight()), ob.chain.String(), nonce)
					return true, false, nil
				}
			}
		}
	}

	return false, false, nil
}

// The lowest nonce we observe outTx for each chain
var lowestOutTxNonceToObserve = map[int64]uint64{
	5:     70000,  // Goerli
	97:    95000,  // BSC testnet
	80001: 120000, // Mumbai
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *EVMChainClient) observeOutTx() {
	// read env variables if set
	timeoutNonce, err := strconv.Atoi(os.Getenv("OS_TIMEOUT_NONCE"))
	if err != nil || timeoutNonce <= 0 {
		timeoutNonce = 100 * 3 // process up to 100 hashes
	}
	rpcRestTime, err := strconv.Atoi(os.Getenv("OS_RPC_REST_TIME"))
	if err != nil || rpcRestTime <= 0 {
		rpcRestTime = 20 // 20ms
	}
	ob.logger.ObserveOutTx.Info().Msgf("observeOutTx using timeoutNonce %d seconds, rpcRestTime %d ms", timeoutNonce, rpcRestTime)

	ticker := NewDynamicTicker(fmt.Sprintf("EVM_observeOutTx_%d", ob.chain.ChainId), ob.GetCoreParams().OutTxTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain, Ascending)
			if err != nil {
				continue
			}
			sort.Slice(trackers, func(i, j int) bool {
				return trackers[i].Nonce < trackers[j].Nonce
			})
			outTimeout := time.After(time.Duration(timeoutNonce) * time.Second)
		TRACKERLOOP:
			// Skip old gabbage trackers as we spent too much time on querying them
			for _, tracker := range trackers {
				nonceInt := tracker.Nonce
				if nonceInt < lowestOutTxNonceToObserve[ob.chain.ChainId] {
					continue
				}
			TXHASHLOOP:
				for _, txHash := range tracker.HashList {
					//inTimeout := time.After(3000 * time.Millisecond)
					select {
					case <-outTimeout:
						ob.logger.ObserveOutTx.Warn().Msgf("observeOutTx timeout on chain %d nonce %d", ob.chain.ChainId, nonceInt)
						break TRACKERLOOP
					default:
						ob.mu.Lock()
						_, found := ob.outTXConfirmedReceipts[ob.GetTxID(nonceInt)]
						ob.mu.Unlock()
						if found {
							continue
						}

						receipt, transaction, err := ob.queryTxByHash(txHash.TxHash, nonceInt)
						time.Sleep(time.Duration(rpcRestTime) * time.Millisecond)
						if err == nil && receipt != nil { // confirmed
							ob.mu.Lock()
							ob.outTXConfirmedReceipts[ob.GetTxID(nonceInt)] = receipt
							ob.outTXConfirmedTransaction[ob.GetTxID(nonceInt)] = transaction
							ob.mu.Unlock()

							break TXHASHLOOP
						}
						if err != nil {
							ob.logger.ObserveOutTx.Debug().Err(err).Msgf("error queryTxByHash: chain %s hash %s", ob.chain.String(), txHash.TxHash)
						}
						//<-inTimeout
					}
				}
			}
			ticker.UpdateInterval(ob.GetCoreParams().OutTxTicker, ob.logger.ObserveOutTx)
		case <-ob.stop:
			ob.logger.ObserveOutTx.Info().Msg("observeOutTx: stopped")
			return
		}
	}
}

// return the status of txHash
// receipt nil, err non-nil: txHash not found
// receipt nil, err nil: txHash receipt recorded, but may not be confirmed
// receipt non-nil, err nil: txHash confirmed
func (ob *EVMChainClient) queryTxByHash(txHash string, nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction, error) {
	logger := ob.logger.ObserveOutTx.With().Str("txHash", txHash).Uint64("nonce", nonce).Logger()
	if ob.outTXConfirmedReceipts[ob.GetTxID(nonce)] != nil && ob.outTXConfirmedTransaction[ob.GetTxID(nonce)] != nil {
		return nil, nil, fmt.Errorf("queryTxByHash: txHash %s receipts already recorded", txHash)
	}
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	receipt, err := ob.EvmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			logger.Warn().Err(err).Msg("TransactionReceipt/TransactionByHash error")
		}
		return nil, nil, err
	}
	transaction, _, err := ob.EvmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		return nil, nil, err
	}
	if transaction.Nonce() != nonce {
		return nil, nil, fmt.Errorf("queryTxByHash: txHash %s nonce mismatch: wanted %d, got tx nonce %d", txHash, nonce, transaction.Nonce())
	}
	confHeight := receipt.BlockNumber.Uint64() + ob.GetCoreParams().ConfirmationCount
	if confHeight < 0 || confHeight >= math2.MaxInt64 {
		return nil, nil, fmt.Errorf("confHeight is out of range")
	}

	if int64(confHeight) > ob.GetLastBlockHeight() {
		log.Warn().Msgf("included but not confirmed: receipt block %d, current block %d", receipt.BlockNumber, ob.GetLastBlockHeight())
		return nil, nil, fmt.Errorf("included but not confirmed")
	}
	return receipt, transaction, nil
}

// SetLastBlockHeightScanned set last block height scanned (not necessarily caught up with external block; could be slow/paused)
func (ob *EVMChainClient) SetLastBlockHeightScanned(block int64) {
	if block < 0 {
		panic("lastBlockScanned is negative")
	}
	if block >= math2.MaxInt64 {
		panic("lastBlockScanned is too large")
	}
	atomic.StoreInt64(&ob.lastBlockScanned, block)
	ob.ts.SetLastScannedBlockNumber(ob.chain.ChainId, block)
}

// GetLastBlockHeightScanned get last block height scanned (not necessarily caught up with external block; could be slow/paused)
func (ob *EVMChainClient) GetLastBlockHeightScanned() int64 {
	height := atomic.LoadInt64(&ob.lastBlockScanned)
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	if height >= math2.MaxInt64 {
		panic("lastBlockScanned is too large")
	}
	return height
}

// SetLastBlockHeight set external last block height (confirmed with confirmation count)
func (ob *EVMChainClient) SetLastBlockHeight(block int64) {
	if block < 0 {
		panic("lastBlock is negative")
	}
	if block >= math2.MaxInt64 {
		panic("lastBlock is too large")
	}
	atomic.StoreInt64(&ob.lastBlock, block)
}

// GetLastBlockHeight get external last block height (confirmed with confirmation count)
func (ob *EVMChainClient) GetLastBlockHeight() int64 {
	height := atomic.LoadInt64(&ob.lastBlock)
	if height < 0 {
		panic("lastBlock is negative")
	}
	if height >= math2.MaxInt64 {
		panic("lastBlock is too large")
	}
	return height
}

func (ob *EVMChainClient) ExternalChainWatcher() {
	// At each tick, query the Connector contract
	ticker := NewDynamicTicker(fmt.Sprintf("EVM_ExternalChainWatcher_%d", ob.chain.ChainId), ob.GetCoreParams().InTxTicker)
	defer ticker.Stop()
	ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher started")
	for {
		select {
		case <-ticker.C():
			err := ob.observeInTX()
			if err != nil {
				ob.logger.ExternalChainWatcher.Err(err).Msg("observeInTX error")
			}
			ticker.UpdateInterval(ob.GetCoreParams().InTxTicker, ob.logger.ExternalChainWatcher)
		case <-ob.stop:
			ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher stopped")
			return
		}
	}
}

func (ob *EVMChainClient) observeInTX() error {
	header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	// "confirmed" current block number
	confirmedBlockNum := header.Number.Uint64() - ob.GetCoreParams().ConfirmationCount
	ob.SetLastBlockHeight(int64(confirmedBlockNum))

	crosschainFlags, err := ob.zetaClient.GetCrosschainFlags()
	if err != nil {
		return err
	}
	if !crosschainFlags.IsInboundEnabled {
		return errors.New("inbound TXS / Send has been disabled by the protocol")
	}
	counter, err := ob.GetPromCounter("rpc_getBlockByNumber_count")
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("GetPromCounter:")
	}
	counter.Inc()

	// skip if no new block is produced.
	sampledLogger := ob.logger.ExternalChainWatcher.Sample(&zerolog.BasicSampler{N: 10})
	if confirmedBlockNum < 0 || confirmedBlockNum > math2.MaxUint64 {
		sampledLogger.Error().Msg("Skipping observer , confirmedBlockNum is negative or too large ")
		return nil
	}
	if confirmedBlockNum <= uint64(ob.GetLastBlockHeightScanned()) {
		sampledLogger.Debug().Msg("Skipping observer , No new block is produced ")
		return nil
	}
	lastBlock := ob.GetLastBlockHeightScanned()
	startBlock := lastBlock + 1
	toBlock := lastBlock + config.MaxBlocksPerPeriod // read at most 10 blocks in one go
	if uint64(toBlock) >= confirmedBlockNum {
		toBlock = int64(confirmedBlockNum)
	}
	if startBlock < 0 || startBlock >= math2.MaxInt64 {
		return fmt.Errorf("startBlock is negative or too large")
	}
	if toBlock < 0 || toBlock >= math2.MaxInt64 {
		return fmt.Errorf("toBlock is negative or too large")
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("Checking for all inTX : startBlock %d, toBlock %d", startBlock, toBlock)
	//task 1:  Query evm chain for zeta sent logs
	func() {
		tb := uint64(toBlock)
		connector, err := ob.GetConnectorContract()
		if err != nil {
			ob.logger.ChainLogger.Warn().Err(err).Msgf("observeInTx: GetConnectorContract error:")
			return
		}
		cnt, err := ob.GetPromCounter("rpc_getLogs_count")
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msg("GetPromCounter:")
		} else {
			cnt.Inc()
		}
		logs, err := connector.FilterZetaSent(&bind.FilterOpts{
			Start:   uint64(startBlock),
			End:     &tb,
			Context: context.TODO(),
		}, []ethcommon.Address{}, []*big.Int{})
		if err != nil {
			ob.logger.ChainLogger.Warn().Err(err).Msgf("observeInTx: FilterZetaSent error:")
			return
		}
		// Pull out arguments from logs
		for logs.Next() {
			event := logs.Event
			ob.logger.ExternalChainWatcher.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
			destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
			if destChain == nil {
				ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
				continue
			}
			destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
			if destChain.IsExternalChain() {
				cfgDest, found := ob.cfg.GetEVMConfig(destChain.ChainId)
				if !found {
					ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not present in EVMChainConfigs  %d", event.DestinationChainId.Int64())
					continue
				}
				if strings.EqualFold(destAddr, cfgDest.ZetaTokenContractAddress) {
					ob.logger.ExternalChainWatcher.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
					continue
				}
			}
			zetaHash, err := ob.zetaClient.PostSend(
				event.ZetaTxSenderAddress.Hex(),
				ob.chain.ChainId,
				event.SourceTxOriginAddress.Hex(),
				clienttypes.BytesToEthHex(event.DestinationAddress),
				destChain.ChainId,
				math.NewUintFromBigInt(event.ZetaValueAndGas),
				base64.StdEncoding.EncodeToString(event.Message),
				event.Raw.TxHash.Hex(),
				event.Raw.BlockNumber,
				event.DestinationGasLimit.Uint64(),
				common.CoinType_Zeta,
				PostSendNonEVMGasLimit,
				"",
			)
			if err != nil {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
				continue
			}
			ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
		}
	}()

	// task 2: Query evm chain for deposited logs
	func() {
		toB := uint64(toBlock)
		erc20custody, err := ob.GetERC20CustodyContract()
		if err != nil {
			ob.logger.ExternalChainWatcher.Warn().Err(err).Msgf("observeInTx: GetERC20CustodyContract error:")
			return
		}
		depositedLogs, err := erc20custody.FilterDeposited(&bind.FilterOpts{
			Start:   uint64(startBlock),
			End:     &toB,
			Context: context.TODO(),
		}, []ethcommon.Address{})

		if err != nil {
			ob.logger.ExternalChainWatcher.Warn().Err(err).Msgf("observeInTx: FilterDeposited error:")
			return
		}
		cnt, err := ob.GetPromCounter("rpc_getLogs_count")
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msg("GetPromCounter:")
		} else {
			cnt.Inc()
		}

		// Pull out arguments from logs
		for depositedLogs.Next() {
			event := depositedLogs.Event
			ob.logger.ExternalChainWatcher.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
			// TODO :add logger to POSTSEND
			if bytes.Compare(event.Message, []byte(DonationMessage)) == 0 {
				ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation!: %s", event.Raw.TxHash.Hex())
				continue
			}

			// get the sender of the event's transaction
			tx, _, err := ob.EvmClient.TransactionByHash(context.Background(), event.Raw.TxHash)
			if err != nil {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg(fmt.Sprintf("failed to get transaction by hash: %s", event.Raw.TxHash.Hex()))
				continue
			}
			signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
			sender, err := signer.Sender(tx)
			if err != nil {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg(fmt.Sprintf("can't recover the sender from the tx hash: %s", event.Raw.TxHash.Hex()))
				continue
			}

			zetaHash, err := ob.zetaClient.PostSend(
				sender.Hex(),
				ob.chain.ChainId,
				"",
				clienttypes.BytesToEthHex(event.Recipient),
				common.ZetaChain().ChainId,
				math.NewUintFromBigInt(event.Amount),
				hex.EncodeToString(event.Message),
				event.Raw.TxHash.Hex(),
				event.Raw.BlockNumber,
				1_500_000,
				common.CoinType_ERC20,
				PostSendEVMGasLimit,
				event.Asset.String(),
			)
			if err != nil {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
				continue
			}
			ob.logger.ExternalChainWatcher.Info().Msgf("ZRC20Custody Deposited event detected and reported: PostSend zeta tx: %s", zetaHash)
		}
	}()

	// task 3: query the incoming tx to TSS address ==============
	func() {
		tssAddress := ob.Tss.EVMAddress() // after keygen, ob.Tss.pubkey will be updated
		if tssAddress == (ethcommon.Address{}) {
			ob.logger.ExternalChainWatcher.Warn().Msgf("observeInTx: TSS address not set")
			return
		}

		// query incoming gas asset
		if !ob.chain.IsKlaytnChain() {
			for bn := startBlock; bn <= toBlock; bn++ {
				//block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
				block, err := ob.GetBlockByNumberCached(bn)
				if err != nil {
					ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("error getting block: %d", bn)
					continue
				}
				_ = ob.BlockCache.Add(block.Hash(), block)
				headerRLP, err := rlp.EncodeToBytes(block.Header())
				if err != nil {
					ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("error encoding block header: %d", bn)
					continue
				}

				_, balloutIdentifier, err := ob.zetaClient.PostAddBlockHeader(
					ob.chain.ChainId,
					block.Hash().Bytes(),
					block.Number().Int64(),
					common.NewEthereumHeader(headerRLP),
				)
				if err != nil {
					ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("error posting block header: %d", bn)
					continue
				}
				ob.logger.ExternalChainWatcher.Info().Msgf("successfully posted block-header: ChainID %d Hash %s , Height %s , Header %v , Ballot %s", ob.chain.ChainId, block.Hash().Hex(), block.Number().String(), common.NewEthereumHeader(headerRLP), balloutIdentifier)
				for _, tx := range block.Transactions() {
					if tx.To() == nil {
						continue
					}
					if bytes.Compare(tx.Data(), []byte(DonationMessage)) == 0 {
						ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation!: %s", tx.Hash().Hex())
						continue
					}

					if *tx.To() == tssAddress {
						receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash())
						if err != nil {
							ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionReceipt error")
							continue
						}
						if receipt.Status != 1 { // 1: successful, 0: failed
							ob.logger.ExternalChainWatcher.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
							continue
						}

						from, err := ob.EvmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
						if err != nil {
							ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionSender error; trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address")
							signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
							from, err = signer.Sender(tx)
							if err != nil {
								ob.logger.ExternalChainWatcher.Err(err).Msg("local recovery of sender address failed")
								continue
							}
						}

						zetaHash, err := ob.ReportTokenSentToTSS(tx.Hash(), tx.Value(), receipt, from, tx.Data())
						if err != nil {
							ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
							continue
						}
						ob.logger.ExternalChainWatcher.Info().Msgf("Gas Deposit detected and reported: PostSend zeta tx: %s", zetaHash)
					}
				}
			}
		} else { // for Klaytn
			for bn := startBlock; bn <= toBlock; bn++ {
				//block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
				block, err := ob.KlaytnClient.BlockByNumber(context.Background(), big.NewInt(bn))
				if err != nil {
					ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("error getting block: %d", bn)
					continue
				}
				for _, tx := range block.Transactions {
					if tx.To == nil {
						continue
					}
					if *tx.To == tssAddress {
						receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash)
						if err != nil {
							ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionReceipt error")
							continue
						}
						if receipt.Status != 1 { // 1: successful, 0: failed
							ob.logger.ExternalChainWatcher.Info().Msgf("tx %s failed; don't act", tx.Hash.Hex())
							continue
						}

						from := *tx.From
						value := tx.Value.ToInt()

						zetaHash, err := ob.ReportTokenSentToTSS(tx.Hash, value, receipt, from, tx.Input)
						if err != nil {
							ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
							continue
						}
						ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
					}
				}
			}
		}
	}()
	// ============= end of query the incoming tx to TSS address ==============
	ob.SetLastBlockHeightScanned(toBlock)
	if err := ob.db.Save(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeightScanned())).Error; err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
}

func (ob *EVMChainClient) ReportTokenSentToTSS(txhash ethcommon.Hash, value *big.Int, receipt *ethtypes.Receipt, from ethcommon.Address, data []byte) (string, error) {
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx detected: %s, blocknum %d", txhash.Hex(), receipt.BlockNumber)
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx value: %s", value.String())
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx from: %s", from.Hex())
	message := ""
	if len(data) != 0 {
		message = hex.EncodeToString(data)
	}
	zetaHash, err := ob.zetaClient.PostSend(
		from.Hex(),
		ob.chain.ChainId,
		from.Hex(),
		from.Hex(),
		common.ZetaChain().ChainId,
		math.NewUintFromBigInt(value),
		message,
		txhash.Hex(),
		receipt.BlockNumber.Uint64(),
		90_000,
		common.CoinType_Gas,
		PostSendEVMGasLimit,
		"",
	)
	return zetaHash, err
}

func (ob *EVMChainClient) WatchGasPrice() {

	err := ob.PostGasPrice()
	if err != nil {
		height, _ := ob.zetaClient.GetBlockHeight()
		ob.logger.WatchGasPrice.Error().Err(err).Msgf("PostGasPrice error at zeta block : %d  ", height)
	}
	ticker := NewDynamicTicker(fmt.Sprintf("EVM_WatchGasPrice_%d", ob.chain.ChainId), ob.GetCoreParams().GasPriceTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.PostGasPrice()
			if err != nil {
				height, _ := ob.zetaClient.GetBlockHeight()
				ob.logger.WatchGasPrice.Error().Err(err).Msgf("PostGasPrice error at zeta block : %d  ", height)
			}
			ticker.UpdateInterval(ob.GetCoreParams().GasPriceTicker, ob.logger.WatchGasPrice)
		case <-ob.stop:
			ob.logger.WatchGasPrice.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *EVMChainClient) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := ob.EvmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("Err SuggestGasPrice:")
		return err
	}
	blockNum, err := ob.EvmClient.BlockNumber(context.TODO())
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("Err Fetching Most recent Block : ")
		return err
	}

	// SUPPLY
	var supply string // lockedAmount on ETH, totalSupply on other chains
	supply = "100"

	zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("PostGasPrice to zetacore failed")
		return err
	}
	_ = zetaHash
	//ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)

	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (ob *EVMChainClient) getLastHeight() (int64, error) {
	lastheight, err := ob.zetaClient.GetLastBlockHeightByChain(ob.chain)
	if err != nil {
		return 0, errors.Wrap(err, "getLastHeight")
	}
	return int64(lastheight.LastSendHeight), nil
}

func (ob *EVMChainClient) BuildBlockIndex() error {
	logger := ob.logger.ChainLogger.With().Str("module", "BuildBlockIndex").Logger()
	envvar := ob.chain.ChainName.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		logger.Info().Msgf("envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == clienttypes.EnvVarLatest {
			header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(header.Number.Int64())
		} else {
			scanFromBlockInt, err := strconv.ParseInt(scanFromBlock, 10, 64)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(scanFromBlockInt)
		}
	} else { // last observed block
		var lastBlockNum clienttypes.LastBlockSQLType
		if err := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumID).Error; err != nil {
			logger.Info().Msg("db PosKey does not exist; read from ZetaCore")
			lastheight, err := ob.getLastHeight()
			if err != nil {
				logger.Warn().Err(err).Msg("getLastHeight error")
			}
			ob.SetLastBlockHeightScanned(lastheight)
			// if ZetaCore does not have last heard block height, then use current
			if ob.GetLastBlockHeightScanned() == 0 {
				header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					return err
				}
				ob.SetLastBlockHeightScanned(header.Number.Int64())
			}
			if dbc := ob.db.Save(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeightScanned())); dbc.Error != nil {
				logger.Error().Err(dbc.Error).Msg("error writing ob.LastBlock to db: ")
			}
		} else {
			ob.SetLastBlockHeightScanned(lastBlockNum.Num)
		}
	}
	return nil
}

func (ob *EVMChainClient) BuildReceiptsMap() error {
	logger := ob.logger
	var receipts []clienttypes.ReceiptSQLType
	if err := ob.db.Find(&receipts).Error; err != nil {
		logger.ChainLogger.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, receipt := range receipts {
		r, err := clienttypes.FromReceiptDBType(receipt.Receipt)
		if err != nil {
			return err
		}
		ob.outTXConfirmedReceipts[receipt.Identifier] = r
	}

	return nil
}

func (ob *EVMChainClient) BuildTransactionsMap() error {
	logger := ob.logger
	var transactions []clienttypes.TransactionSQLType
	if err := ob.db.Find(&transactions).Error; err != nil {
		logger.ChainLogger.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, transaction := range transactions {
		trans, err := clienttypes.FromTransactionDBType(transaction.Transaction)
		if err != nil {
			return err
		}
		ob.outTXConfirmedTransaction[transaction.Identifier] = trans
	}
	return nil
}

// LoadDB open sql database and load data into EVMChainClient
func (ob *EVMChainClient) LoadDB(dbPath string, chain common.Chain) error {
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
			return err
		}

		ob.db = db
		err = ob.BuildBlockIndex()
		if err != nil {
			return err
		}

		//DISABLING RECEIPT AND TRANSACTION PERSISTENCE
		//err = ob.BuildReceiptsMap()
		//if err != nil {
		//	return err
		//}
		//
		//err = ob.BuildTransactionsMap()
		//if err != nil {
		//	return err
		//}

	}
	return nil
}

func (ob *EVMChainClient) SetMinAndMaxNonce(trackers []types.OutTxTracker) error {
	minNonce, maxNonce := int64(-1), int64(0)
	for _, tracker := range trackers {
		conv := tracker.Nonce
		intNonce := int64(conv)
		if minNonce == -1 {
			minNonce = intNonce
		}
		if intNonce < minNonce {
			minNonce = intNonce
		}
		if intNonce > maxNonce {
			maxNonce = intNonce
		}
	}
	if minNonce != -1 {
		atomic.StoreInt64(&ob.MinNonce, minNonce)
	}
	if maxNonce > 0 {
		atomic.StoreInt64(&ob.MaxNonce, maxNonce)
	}
	return nil
}

func (ob *EVMChainClient) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.EVMAddress().String()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

func (ob *EVMChainClient) GetBlockByNumberCached(blockNumber int64) (*ethtypes.Block, error) {
	if block, ok := ob.BlockCache.Get(blockNumber); ok {
		return block.(*ethtypes.Block), nil
	}
	block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.BlockCache.Add(blockNumber, block)
	return block, nil
}
