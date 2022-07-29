package zetaclient

import (
	"context"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/contracts/evm"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
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
	Nonce    int
}

// Chain configuration struct
// Filled with above constants depending on chain
type ChainObserver struct {
	chain                  common.Chain
	endpoint               string
	ticker                 *time.Ticker
	Connector              *evm.Connector
	ConnectorAddress       ethcommon.Address
	EvmClient              *ethclient.Client
	zetaClient             *ZetaCoreBridge
	Tss                    TSSSigner
	LastBlock              uint64
	confCount              uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime              uint64 // block time in seconds
	txWatchList            map[ethcommon.Hash]string
	mu                     *sync.Mutex
	db                     *leveldb.DB
	sampleLogger           *zerolog.Logger
	metrics                *metrics.Metrics
	outTXConfirmedReceipts map[int]*ethtypes.Receipt
	MinNonce               int
	MaxNonce               int
	OutTxChan              chan OutTx // send to this channel if you want something back!
	ZetaPriceQuerier       ZetaPriceQuerier
	stop                   chan struct{}
	wg                     sync.WaitGroup

	fileLogger *zerolog.Logger // for critical info
}

// Return configuration based on supplied target chain
func NewChainObserver(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metrics.Metrics) (*ChainObserver, error) {
	ob := ChainObserver{}
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLogger = &sampled
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.metrics = metrics
	ob.outTXConfirmedReceipts = make(map[int]*ethtypes.Receipt)
	ob.OutTxChan = make(chan OutTx, 100)
	addr := ethcommon.HexToAddress(config.Chains[chain.String()].ConnectorContractAddress)
	if addr == ethcommon.HexToAddress("0x0") {
		return nil, fmt.Errorf("Connector contract address %s not configured for chain %s", config.Chains[chain.String()].ConnectorContractAddress, chain.String())
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
	log.Info().Msgf("Chain %s endpoint %s", ob.chain, ob.endpoint)
	client, err := ethclient.Dial(ob.endpoint)
	if err != nil {
		log.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client

	// initialize the connector
	connector, err := evm.NewConnector(addr, ob.EvmClient)
	if err != nil {
		log.Error().Err(err).Msg("Connector")
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
		log.Info().Msg("Using dummy price of 1:1")
		ob.ZetaPriceQuerier = dummyQuerior
	}
	if dbpath != "" {
		err := ob.BuildBlockIndex(dbpath, chain.String())
		if err != nil {
			return nil, err
		}
		ob.BuildReceiptsMap()

	}
	log.Info().Msgf("%s: start scanning from block %d", chain, ob.LastBlock)

	return &ob, nil
}

func (ob *ChainObserver) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.WatchExchangeRate()    // Observers ZetaPriceQuerier for Zeta prices and posts to core
	go ob.observeOutTx()
}

func (ob *ChainObserver) Stop() {
	log.Info().Msgf("ob %s is stopping", ob.chain)
	close(ob.stop) // this notifies all goroutines to stop

	log.Info().Msg("closing ob.db")
	err := ob.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("error closing db")
	}

	log.Info().Msgf("%s observer stopped", ob.chain)
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *ChainObserver) IsSendOutTxProcessed(sendHash string, nonce int) (bool, bool, error) {
	receipt, found := ob.outTXConfirmedReceipts[nonce]
	if found && receipt.Status == 1 {
		logs := receipt.Logs
		for _, vLog := range logs {
			receivedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReceived(*vLog)
			if err == nil {
				log.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					sendhash := vLog.Topics[3].Hex()
					//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
					mMint := receivedLog.ZetaAmount.String()
					zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
						sendhash,
						vLog.TxHash.Hex(),
						vLog.BlockNumber,
						mMint,
						common.ReceiveStatus_Success,
						ob.chain.String(),
						nonce,
					)
					if err != nil {
						log.Error().Err(err).Msg("error posting confirmation to meta core")
						continue
					}
					log.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
					return true, true, nil
				} else {
					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
					return true, false, nil
				}
			}
			revertedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReverted(*vLog)
			if err == nil {
				log.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					sendhash := vLog.Topics[3].Hex()
					mMint := revertedLog.ZetaAmount.String()
					metaHash, err := ob.zetaClient.PostReceiveConfirmation(
						sendhash,
						vLog.TxHash.Hex(),
						vLog.BlockNumber,
						mMint,
						common.ReceiveStatus_Success,
						ob.chain.String(),
						nonce,
					)
					if err != nil {
						log.Err(err).Msg("error posting confirmation to meta core")
						continue
					}
					log.Info().Msgf("Zeta tx hash: %s", metaHash)
					return true, true, nil
				} else {
					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
					return true, false, nil
				}
			}
		}
	} else if found && receipt.Status == 0 {
		//FIXME: check nonce here by getTransaction RPC
		log.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, receipt.TxHash.Hex())
		zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, ob.chain.String(), nonce)
		if err != nil {
			log.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
		}
		log.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
		return true, true, nil
	}

	return false, false, fmt.Errorf("IsSendOutTxProcessed: error on chain %s", ob.chain)
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *ChainObserver) observeOutTx() {
	ticker := time.NewTicker(12 * time.Second)
	for {
		select {
		case <-ticker.C:
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain)
			if err != nil {
				return
			}
			if len(trackers) > 0 {
				err = ob.SetMinAndMaxNonce(trackers)
				if err != nil {
					return
				}
			}
			outTimeout := time.After(12 * time.Second)
		TRACKERLOOP:
			for _, tracker := range trackers {
				nonceInt, err := strconv.Atoi(tracker.Nonce)
				if err != nil {
					return
				}
			TXHASHLOOP:
				for _, txHash := range tracker.HashList {
					inTimeout := time.After(1000 * time.Millisecond)
					select {
					case <-outTimeout:
						log.Warn().Msgf("Timeout chain %s nonce %d", ob.chain, nonceInt)
						break TRACKERLOOP
					default:
						receipt, err := ob.queryTxByHash(txHash.TxHash, nonceInt)
						if err == nil && receipt != nil { // confirmed
							ob.mu.Lock()
							ob.outTXConfirmedReceipts[nonceInt] = receipt
							value, err := receipt.MarshalJSON()
							if err != nil {
								log.Error().Err(err).Msgf("receipt marshal error %s", receipt.TxHash.Hex())
							}
							ob.mu.Unlock()
							err = ob.db.Put([]byte(NonceTxKeyPrefix+fmt.Sprintf("%d", nonceInt)), value, nil)
							if err != nil {
								log.Error().Err(err).Msgf("PurgeTxHashWatchList: error putting nonce %d tx hashes %s to db", nonceInt, receipt.TxHash.Hex())
							}
							break TXHASHLOOP
						}
						<-inTimeout
					}
				}
			}
		case <-ob.stop:
			log.Info().Msg("observeOutTx: stopped")
			return
		}
	}
}

// return the status of txHash
// receipt nil, err non-nil: txHash not found
// receipt nil, err nil: txHash receipt recorded, but may not be confirmed
// receipt non-nil, err nil: txHash confirmed
func (ob *ChainObserver) queryTxByHash(txHash string, nonce int) (*ethtypes.Receipt, error) {
	if ob.outTXConfirmedReceipts[nonce] != nil {
		return nil, fmt.Errorf("queryTxByHash: txHash %s recepits already recorded", txHash)
	}
	//timeStart := time.Now()
	//defer func() { log.Info().Msgf("queryTxByHash elapsed: %s", time.Since(timeStart)) }()
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	receipt, err := ob.EvmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			log.Warn().Err(err).Msgf("%s %s TransactionReceipt err", ob.chain, txHash)
		}
		return nil, err
	} else if receipt.BlockNumber.Uint64()+ob.confCount > ob.LastBlock {
		log.Info().Msgf("%s TransactionReceipt %s included in block %d but not confirmed; current block num %d", ob.chain, txHash, receipt.BlockNumber.Uint64(), ob.LastBlock)
		return nil, fmt.Errorf("included but not confirmed")
	} else {
		if receipt.Status == 0 { // failed (reverted tx)
			log.Info().Msgf("%s TransactionReceipt %s nonce %d mined and confirmed, but it's reverted!", ob.chain, txHash, nonce)
		} else if receipt.Status == 1 { // success
			log.Info().Msgf("%s TransactionReceipt %s nonce %d mined and confirmed, and it's successful", ob.chain, txHash, nonce)
		}
		return receipt, nil
	}
}
