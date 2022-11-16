package zetaclient

import (
	"context"
	"fmt"
	"github.com/zeta-chain/zetacore/common"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/contracts/evm"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	PosKey                 = "PosKey"
	NonceTxHashesKeyPrefix = "NonceTxHashes-"
	NonceTxKeyPrefix       = "NonceTx-"
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

// Chain configuration struct
// Filled with above constants depending on chain
type ChainObserver struct {
	chain                     zetaObserverTypes.Chain
	endpoint                  string
	ticker                    *time.Ticker
	Connector                 *evm.Connector
	ConnectorAddress          ethcommon.Address
	EvmClient                 *ethclient.Client
	zetaClient                *ZetaCoreBridge
	Tss                       TSSSigner
	lastBlock                 uint64
	confCount                 uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime                 uint64 // block time in seconds
	txWatchList               map[ethcommon.Hash]string
	mu                        *sync.Mutex
	db                        *leveldb.DB
	sampleLogger              *zerolog.Logger
	metrics                   *metricsPkg.Metrics
	outTXConfirmedReceipts    map[int]*ethtypes.Receipt
	outTXConfirmedTransaction map[int]*ethtypes.Transaction
	MinNonce                  int64
	MaxNonce                  int64
	OutTxChan                 chan OutTx // send to this channel if you want something back!
	ZetaPriceQuerier          ZetaPriceQuerier
	stop                      chan struct{}
	fileLogger                *zerolog.Logger // for critical info
	logger                    zerolog.Logger
}

// Return configuration based on supplied target chain
func NewChainObserver(chain zetaObserverTypes.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics) (*ChainObserver, error) {
	ob := ChainObserver{}
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLogger = &sampled
	ob.logger = log.With().Str("chain", chain.String()).Logger()
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.metrics = metrics
	ob.outTXConfirmedReceipts = make(map[int]*ethtypes.Receipt)
	ob.outTXConfirmedTransaction = make(map[int]*ethtypes.Transaction)
	ob.OutTxChan = make(chan OutTx, 100)
	addr := ethcommon.HexToAddress(config.Chains[chain.String()].ConnectorContractAddress)
	if addr == ethcommon.HexToAddress("0x0") {
		return nil, fmt.Errorf("connector contract address %s not configured for chain %s", config.Chains[chain.String()].ConnectorContractAddress, chain.String())
	}
	ob.ConnectorAddress = addr
	ob.endpoint = config.Chains[chain.String()].Endpoint
	logFile, err := os.OpenFile(ob.chain.String()+"_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		// Can we log an error before we have our logger? :)
		log.Error().Err(err).Msgf("there was an error creating a logFile chain %s", ob.chain.String())
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	ob.fileLogger = &fileLogger

	// initialize the Client
	ob.logger.Info().Msgf("Chain %s endpoint %s", ob.chain, ob.endpoint)
	client, err := ethclient.Dial(ob.endpoint)
	if err != nil {
		ob.logger.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client

	// initialize the connector
	connector, err := evm.NewConnector(addr, ob.EvmClient)
	if err != nil {
		ob.logger.Error().Err(err).Msg("Connector")
		return nil, err
	}
	ob.Connector = connector

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

	uniswapV3ABI, err := abi.JSON(strings.NewReader(config.UNISWAPV3POOL))
	if err != nil {
		return nil, err
	}
	uniswapV2ABI, err := abi.JSON(strings.NewReader(config.PANCAKEPOOL))
	if err != nil {
		return nil, err
	}

	uniswapv3Querier, uniswapv2Querier, dummyQuerior := ob.GetPriceQueriers(chain.String(), uniswapV3ABI, uniswapV2ABI)
	ob.SetChainDetails(chain, uniswapv3Querier, uniswapv2Querier)
	if os.Getenv("DUMMY_PRICE") != "" {
		ob.logger.Info().Msg("Using dummy price of 1:1")
		ob.ZetaPriceQuerier = dummyQuerior
	}

	if dbpath != "" {
		err := ob.BuildBlockIndex(dbpath, chain.String())
		if err != nil {
			return nil, err
		}
		ob.BuildReceiptsMap()

	}
	ob.logger.Info().Msgf("%s: start scanning from block %d", chain, ob.GetLastBlock())

	return &ob, nil
}

func (ob *ChainObserver) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.WatchExchangeRate()    // Observers ZetaPriceQuerier for Zeta prices and posts to core
	go ob.observeOutTx()
}

func (ob *ChainObserver) Stop() {
	ob.logger.Info().Msgf("ob %s is stopping", ob.chain)
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.Info().Msg("closing ob.db")
	err := ob.db.Close()
	if err != nil {
		ob.logger.Error().Err(err).Msg("error closing db")
	}

	ob.logger.Info().Msgf("%s observer stopped", ob.chain)
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *ChainObserver) IsSendOutTxProcessed(sendHash string, nonce int, fromOrToZeta bool) (bool, bool, error) {
	ob.mu.Lock()
	receipt, found1 := ob.outTXConfirmedReceipts[nonce]
	transaction, found2 := ob.outTXConfirmedTransaction[nonce]
	ob.mu.Unlock()
	found := found1 && found2
	sendID := fmt.Sprintf("%s/%d", ob.chain.String(), nonce)
	logger := ob.logger.With().Str("sendID", sendID).Logger()
	if fromOrToZeta {
		if found && receipt.Status == 1 {
			zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
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
		} else if found && receipt.Status == 0 { // the same as below events flow
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), big.NewInt(0), common.ReceiveStatus_Failed, ob.chain, nonce, common.CoinType_Gas)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	} else {
		if found && receipt.Status == 1 {
			logs := receipt.Logs
			for _, vLog := range logs {
				receivedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReceived(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
					if vLog.BlockNumber+ob.confCount < ob.GetLastBlock() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						if len(vLog.Topics) != 4 {
							logger.Error().Msgf("wrong number of topics in log %d", len(vLog.Topics))
							return false, false, fmt.Errorf("wrong number of topics in log %d", len(vLog.Topics))
						}
						sendhash := vLog.Topics[3].Hex()
						//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
						mMint := receivedLog.ZetaValue
						logger.Warn().Msgf("######## chain %s ########", ob.chain.String())
						zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
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
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlock()), ob.chain, nonce)
					return true, false, nil
				}
				revertedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReverted(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
					if vLog.BlockNumber+ob.confCount < ob.GetLastBlock() {
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
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlock()), ob.chain, nonce)
					return true, false, nil
				}
			}
		} else if found && receipt.Status == 0 {
			//FIXME: check nonce here by getTransaction RPC
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), big.NewInt(0), common.ReceiveStatus_Failed, ob.chain, nonce, common.CoinType_Zeta)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	}

	return false, false, nil
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *ChainObserver) observeOutTx() {
	logger := ob.logger
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain)
			if err != nil {
				return
			}
			outTimeout := time.After(90 * time.Second)
		TRACKERLOOP:
			for _, tracker := range trackers {
				nonceInt, err := strconv.Atoi(tracker.Nonce)
				if err != nil {
					return
				}
			TXHASHLOOP:
				for _, txHash := range tracker.HashList {
					inTimeout := time.After(3000 * time.Millisecond)
					select {
					case <-outTimeout:
						logger.Warn().Msgf("observeOutTx timeout on nonce %d", nonceInt)
						break TRACKERLOOP
					default:
						receipt, transaction, err := ob.queryTxByHash(txHash.TxHash, int64(nonceInt))
						if err == nil && receipt != nil { // confirmed
							ob.mu.Lock()
							ob.outTXConfirmedReceipts[nonceInt] = receipt
							ob.outTXConfirmedTransaction[nonceInt] = transaction
							value, err := receipt.MarshalJSON()
							if err != nil {
								logger.Error().Err(err).Msgf("receipt marshal error %s", receipt.TxHash.Hex())
							}
							ob.mu.Unlock()
							err = ob.db.Put([]byte(NonceTxKeyPrefix+fmt.Sprintf("%d", nonceInt)), value, nil)
							if err != nil {
								logger.Error().Err(err).Msgf("PurgeTxHashWatchList: error putting nonce %d tx hashes %s to db", nonceInt, receipt.TxHash.Hex())
							}
							break TXHASHLOOP
						}
						<-inTimeout
					}
				}
			}
		case <-ob.stop:
			logger.Info().Msg("observeOutTx: stopped")
			return
		}
	}
}

// return the status of txHash
// receipt nil, err non-nil: txHash not found
// receipt nil, err nil: txHash receipt recorded, but may not be confirmed
// receipt non-nil, err nil: txHash confirmed
func (ob *ChainObserver) queryTxByHash(txHash string, nonce int64) (*ethtypes.Receipt, *ethtypes.Transaction, error) {
	logger := ob.logger.With().Str("txHash", txHash).Int64("nonce", nonce).Logger()
	if ob.outTXConfirmedReceipts[int(nonce)] != nil && ob.outTXConfirmedTransaction[int(nonce)] != nil {
		return nil, nil, fmt.Errorf("queryTxByHash: txHash %s receipts already recorded", txHash)
	}
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	receipt, err1 := ob.EvmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	transaction, _, err2 := ob.EvmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))

	if err1 != nil || err2 != nil {
		if err1 != ethereum.NotFound {
			logger.Warn().Err(err1).Msg("TransactionReceipt/TransactionByHash error")
		}
		return nil, nil, err1
	} else if receipt.BlockNumber.Uint64()+ob.confCount > ob.GetLastBlock() {
		log.Warn().Msgf("included but not confirmed: receipt block %d, current block %d", receipt.BlockNumber, ob.GetLastBlock())
		return nil, nil, fmt.Errorf("included but not confirmed")
	} else { // confirmed outbound tx
		return receipt, transaction, nil
	}
}

func (ob *ChainObserver) setLastBlock(block uint64) {
	atomic.StoreUint64(&ob.lastBlock, block)
}

func (ob *ChainObserver) GetLastBlock() uint64 {
	return atomic.LoadUint64(&ob.lastBlock)
}
