package observer

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/adapters/bridge"
	obs "github.com/zeta-chain/zetacore/zetaclient/adapters/observer"
	"github.com/zeta-chain/zetacore/zetaclient/adapters/pricer"
	"github.com/zeta-chain/zetacore/zetaclient/adapters/signer"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/model"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
)

type Observer struct {
	chain                  common.Chain
	chainObserver          obs.ChainObserver
	endpoint               string
	ticker                 *time.Ticker
	zetaClient             *bridge.ZetaCoreBridge
	Tss                    signer.TSSSigner
	lastBlock              uint64
	confCount              uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime              uint64 // block time in seconds
	mu                     *sync.Mutex
	db                     *leveldb.DB
	sampleLogger           *zerolog.Logger
	metrics                *metricsPkg.Metrics
	outTXConfirmedReceipts map[int]*ethtypes.Receipt
	outTxChan              chan model.OutTx // send to this channel if you want something back!
	ZetaPriceQuerier       pricer.ZetaPriceQuerier
	stop                   chan struct{}
	fileLogger             *zerolog.Logger // for critical info
	logger                 zerolog.Logger
}

func NewObserver(ctx context.Context, chain common.Chain, bridge *bridge.ZetaCoreBridge, tss signer.TSSSigner, dbpath string, metrics *metricsPkg.Metrics) (*Observer, error) {
	ob := Observer{}
	ob.ctx = ctx
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLogger = &sampled
	ob.logger = log.With().Str("chain", chain.String()).Logger()
	ob.zetaClient = bridge
	ob.Tss = tss
	ob.metrics = metrics
	ob.outTXConfirmedReceipts = make(map[int]*ethtypes.Receipt)
	ob.outTxChan = make(chan model.OutTx, 100)
	logFile, err := os.OpenFile(ob.chain.String()+"_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		// Can we log an error before we have our logger? :)
		log.Error().Err(err).Msgf("there was an error creating a logFile chain %s", ob.chain.String())
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	ob.fileLogger = &fileLogger

	// initialize chain observer
	switch ob.chain.Type {
	case ETH:
		ob.chainObserver = ethInfra.NewEthChainObserver(ob.ctx, ob.chain, ob.fileLogger)
	case BTC:
		ob.chainObserver = btcInfra.NewBtcChainObserver(ob.ctx, ob.chain, ob.fileLogger)
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

	if dbpath != "" {
		err := ob.BuildBlockIndex(dbpath, chain.String())
		if err != nil {
			return nil, err
		}
		ob.chainObserver.BuildReceiptsMap()
	}
	ob.logger.Info().Msgf("%s: start scanning from block %d", chain, ob.GetLastBlock())

	return &ob, nil
}

func (ob *Observer) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.WatchExchangeRate()    // Observers ZetaPriceQuerier for Zeta prices and posts to core
	go ob.ObserveOutTx()
}

func (ob *Observer) Stop() {
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
func (ob *Observer) IsSendOutTxProcessed(sendHash string, nonce int) (bool, bool, error) {
	ob.mu.Lock()
	receipt, found := ob.outTXConfirmedReceipts[nonce]
	ob.mu.Unlock()
	sendID := fmt.Sprintf("%s/%d", ob.chain.String(), nonce)
	logger := ob.logger.With().Str("sendID", sendID).Logger()
	if found && receipt.Status == 1 {
		logs := receipt.Logs
		for _, vLog := range logs {
			receivedLog, err := ob.chainObserver.GetConnectorReceivedLog(ob.ctx, *vLog)
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
					mMint := receivedLog.ZetaValue.String()
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
						logger.Error().Err(err).Msg("error posting confirmation to meta core")
						continue
					}
					logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
					return true, true, nil
				}
				logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlock()), ob.chain, nonce)
				return true, false, nil
			}
			revertedLog, err := ob.chainObserver.GetConnectorRevertedLog(ob.ctx, *vLog)
			if err == nil {
				logger.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
				if vLog.BlockNumber+ob.confCount < ob.GetLastBlock() {
					logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					if len(vLog.Topics) != 3 {
						logger.Error().Msgf("wrong number of topics in log %d", len(vLog.Topics))
						return false, false, fmt.Errorf("wrong number of topics in log %d", len(vLog.Topics))
					}
					sendhash := vLog.Topics[2].Hex()
					mMint := revertedLog.RemainingZetaValue.String()
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
		zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, ob.chain.String(), nonce)
		if err != nil {
			logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
		}
		logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
		return true, true, nil
	}

	return false, false, nil
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *Observer) ObserveOutTx() {
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
					inTimeout := time.After(1000 * time.Millisecond)
					select {
					case <-outTimeout:
						logger.Warn().Msgf("observeOutTx timeout on nonce %d", nonceInt)
						break TRACKERLOOP
					default:
						receipt, err := ob.chainObserver.QueryTxByHash(ob.ctx, txHash.TxHash, int64(nonceInt))
						if err == nil && receipt != nil { // confirmed
							ob.mu.Lock()
							ob.outTXConfirmedReceipts[nonceInt] = receipt
							value, err := receipt.MarshalJSON()
							if err != nil {
								logger.Error().Err(err).Msgf("receipt marshal error %s", receipt.TxHash.Hex())
							}
							ob.mu.Unlock()
							err = ob.db.Put([]byte(model.NonceTxKeyPrefix+fmt.Sprintf("%d", nonceInt)), value, nil)
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

func (ob *Observer) setLastBlock(block uint64) {
	atomic.StoreUint64(&ob.lastBlock, block)
}

func (ob *Observer) GetLastBlock() uint64 {
	return atomic.LoadUint64(&ob.lastBlock)
}

func (ob *Observer) BlockTimeSeconds() uint64 {
	return ob.BlockTime
}

func (ob *Observer) Chain() *common.Chain {
	return &ob.chain
}

func (ob *Observer) ConfirmationsCount() uint64 {
	return ob.confCount
}

func (ob *Observer) CriticalLog() *zerolog.Logger {
	return ob.fileLogger
}

func (ob *Observer) Log() zerolog.Logger {
	return ob.logger
}

func (ob *Observer) Endpoint() string {
	return ob.endpoint
}

func (ob *Observer) LastBlock() uint64 {
	return ob.lastBlock
}

func (ob *Observer) OutTxChan() chan model.OutTx {
	return ob.outTxChan
}

func (ob *Observer) TSSSigner() signer.TSSSigner {
	return ob.Tss
}

func (ob *Observer) Ticker() *time.Ticker {
	return ob.ticker
}
