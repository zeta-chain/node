package zetaclient

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"

	prom "github.com/prometheus/client_golang/prometheus"
)

const (
	OutboundTxSignCount = "zetaclient_Outbound_tx_sign_count"
)

type ZetaCoreLog struct {
	ChainLogger      zerolog.Logger
	ZetaChainWatcher zerolog.Logger
}

type CoreObserver struct {
	bridge    *ZetaCoreBridge
	signerMap map[common.Chain]ChainSigner
	clientMap map[common.Chain]ChainClient
	scanner   *CctxScanner
	metrics   *metrics.Metrics
	tss       *TSS
	logger    ZetaCoreLog
	cfg       *config.Config
	ts        *TelemetryServer
	stop      chan struct{}
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]ChainSigner, clientMap map[common.Chain]ChainClient, dbpath string, metrics *metrics.Metrics, tss *TSS, logger zerolog.Logger, cfg *config.Config, ts *TelemetryServer) (*CoreObserver, error) {
	co := CoreObserver{
		ts:   ts,
		stop: make(chan struct{}),
	}
	co.cfg = cfg
	chainLogger := logger.With().
		Str("chain", "ZetaChain").
		Logger()
	co.logger = ZetaCoreLog{
		ChainLogger:      chainLogger,
		ZetaChainWatcher: chainLogger.With().Str("module", "ZetaChainWatcher").Logger(),
	}

	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap

	co.clientMap = clientMap
	co.metrics = metrics
	co.logger.ChainLogger.Info().Msg("starting core observer")
	err := metrics.RegisterCounter(OutboundTxSignCount, "number of Outbound tx signed")
	if err != nil {
		co.logger.ChainLogger.Error().Err(err).Msg("error registering counter")
	}

	scanner, err := NewCctxScanner(bridge, dbpath, false, &co.logger.ZetaChainWatcher)
	if err != nil {
		return nil, err
	}
	co.scanner = scanner

	return &co, nil
}

func (co *CoreObserver) GetPromCounter(name string) (prom.Counter, error) {
	cnt, found := metrics.Counters[name]
	if !found {
		return nil, errors.New("counter not found")
	}
	return cnt, nil
}

func (co *CoreObserver) MonitorCore() {
	myid := co.bridge.keys.GetAddress()
	co.logger.ZetaChainWatcher.Info().Msgf("Starting Send Scheduler for %s", myid)
	go co.startSendScheduler()

	go func() {
		// bridge queries UpgradePlan from zetacore and send to its pause channel if upgrade height is reached
		<-co.bridge.pause
		// now stop everything
		close(co.stop) // this stops the startSendScheduler() loop
		for _, c := range co.clientMap {
			c.Stop()
		}
	}()
}

// ZetaCore block is heart beat; each block we schedule some send according to
// retry schedule. ses
func (co *CoreObserver) startSendScheduler() {
	outTxMan := NewOutTxProcessorManager(co.logger.ChainLogger)
	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum int64
	for {
		select {
		case <-co.stop:
			co.logger.ZetaChainWatcher.Warn().Msg("stop sendScheduler")
			return
		case <-observeTicker.C:
			{
				bn, err := co.bridge.GetZetaBlockHeight()
				if err != nil {
					co.logger.ZetaChainWatcher.Error().Msg("GetZetaBlockHeight fail in startSendScheduler")
					continue
				}
				if lastBlockNum == 0 {
					lastBlockNum = bn - 1
				}
				if bn > lastBlockNum { // we have a new block
					bn = lastBlockNum + 1
					if bn%10 == 0 {
						co.logger.ZetaChainWatcher.Debug().Msgf("ZetaCore heart beat: %d", bn)
					}
					//logger.Info().Dur("elapsed", time.Since(tStart)).Msgf("GetAllPendingCctx %d", len(sendList))

					supportedChains := GetSupportedChains()
					for _, c := range supportedChains {
						if c == nil || c.ChainId == common.ZetaChain().ChainId {
							continue
						}
						signer := co.signerMap[*c]
						chainClient := co.clientMap[*c]
						sendList, pendingNonces, err := co.bridge.GetAllPendingCctx(uint64(c.ChainId))
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("failed to GetAllPendingCctx for chain %s", c.ChainName.String())
							continue
						}

						// Scan missed pending cctx in history
						missedList := co.scanner.ScanMissedPendingCctx(c.ChainId, pendingNonces)
						if missedList != nil {
							sendList = append(missedList, sendList...)
						}

						ob, err := co.getUpdatedChainOb(c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("getTargetChainOb fail, Chain ID: %s", c.ChainName)
							continue
						}
						chain, err := common.GetChainNameFromChainID(c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("GetTargetChain fail, Chain ID: %s", c.ChainName)
							continue
						}
						res, err := co.bridge.GetAllOutTxTrackerByChain(*c, Ascending)
						if err != nil {
							co.logger.ZetaChainWatcher.Warn().Err(err).Msgf("failed to GetAllOutTxTrackerByChain for chain %s", c.ChainName.String())
							continue
						}
						trackerMap := make(map[uint64]bool)
						for _, v := range res {
							trackerMap[v.Nonce] = true
						}

						gauge, err := ob.GetPromGauge(metrics.PendingTxs)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("failed to get prometheus gauge: %s", metrics.PendingTxs)
							continue
						}
						gauge.Set(float64(len(sendList)))

						for idx, send := range sendList {
							params := send.GetCurrentOutTxParam()
							if params.ReceiverChainId != c.ChainId {
								log.Warn().Msgf("mismatch chainid: want %d, got %d", c.ChainId, params.ReceiverChainId)
								continue
							}

							// Monitor Core Logger for OutboundTxTssNonce
							included, _, err := ob.IsSendOutTxProcessed(send.Index, params.OutboundTxTssNonce, params.CoinType, co.logger.ZetaChainWatcher)
							if err != nil {
								co.logger.ZetaChainWatcher.Error().Err(err).Msgf("IsSendOutTxProcessed fail, Chain ID: %s", c.ChainName)
								continue
							}
							// skip keysign for missed pending cctx as outTx was already finalized
							if co.scanner.IsMissedPendingCctx(c.ChainId, params.OutboundTxTssNonce) {
								co.scanner.UpdateMissedPendingCctx(c.ChainId, params.OutboundTxTssNonce, uint64(pendingNonces.NonceLow))
								continue
							}
							if included {
								co.logger.ZetaChainWatcher.Info().Msgf("send outTx already included; do not schedule")
								continue
							}
							nonce := params.OutboundTxTssNonce
							outTxID := fmt.Sprintf("%s-%d-%d", send.Index, params.ReceiverChainId, nonce) // should be the outTxID?

							// FIXME: config this schedule; this value is for localnet fast testing
							if bn >= math.MaxInt64 {
								continue
							}
							currentHeight := uint64(bn)
							interval := uint64(ob.GetCoreParams().OutboundTxScheduleInterval)
							lookahead := uint64(ob.GetCoreParams().OutboundTxScheduleLookahead)

							// determining critical outtx; if it satisfies following criteria
							// 1. it's the first pending outtx for this chain
							// 2. the following 5 nonces have been in tracker
							criticalInterval := uint64(10)      // for critical pending outTx we reduce re-try interval
							nonCriticalInterval := interval * 2 // for non-critical pending outTx we increase re-try interval
							if nonce%criticalInterval == currentHeight%criticalInterval {
								count := 0
								for i := nonce + 1; i <= nonce+10; i++ {
									if _, found := trackerMap[i]; found {
										count++
									}
								}
								if count >= 5 {
									interval = criticalInterval
								}
							}
							// if it's already in tracker, we increase re-try interval
							if _, ok := trackerMap[nonce]; ok {
								interval = nonCriticalInterval
							}

							// otherwise, the normal interval is used
							if nonce%interval == currentHeight%interval && !outTxMan.IsOutTxActive(outTxID) {
								outTxMan.StartTryProcess(outTxID)
								co.logger.ZetaChainWatcher.Debug().Msgf("chain %s: Sign outtx %s with value %d\n", chain, send.Index, send.GetCurrentOutTxParam().Amount)
								go signer.TryProcessOutTx(send, outTxMan, outTxID, chainClient, co.bridge, currentHeight)
							}
							if idx > int(lookahead) { // only look at 50 sends per chain
								break
							}
						}
					}
					// update last processed block number
					lastBlockNum = bn
					co.ts.SetCoreBlockNumber(lastBlockNum)
				}
			}

		}
	}
}

func (co *CoreObserver) getUpdatedChainOb(chainID int64) (ChainClient, error) {
	chainOb, err := co.getTargetChainOb(chainID)
	if err != nil {
		return nil, err
	}
	// update chain client core parameters
	curParams := chainOb.GetCoreParams()
	if common.IsEVMChain(chainID) {
		evmCfg, found := co.cfg.GetEVMConfig(chainID)
		if found && curParams != evmCfg.CoreParams {
			chainOb.SetCoreParams(evmCfg.CoreParams)
			co.logger.ZetaChainWatcher.Info().Msgf("updated core params for chainID %d, new params: %v", chainID, evmCfg.CoreParams)
		}
	} else if common.IsBitcoinChain(chainID) {
		_, btcCfg, found := co.cfg.GetBTCConfig()
		if found && curParams != btcCfg.CoreParams {
			chainOb.SetCoreParams(btcCfg.CoreParams)
			co.logger.ZetaChainWatcher.Info().Msgf("updated core params for Bitcoin, new params: %v", btcCfg.CoreParams)
		}
	}
	return chainOb, nil
}

func (co *CoreObserver) getTargetChainOb(chainID int64) (ChainClient, error) {
	c := common.GetChainFromChainID(chainID)
	if c == nil {
		return nil, fmt.Errorf("chain not found for chainID %d", chainID)
	}
	chainOb, found := co.clientMap[*c]
	if !found {
		return nil, fmt.Errorf("chain client not found for chainID %d", chainID)
	}
	return chainOb, nil
}

// returns whether to retry in a few seconds, and whether to report via AddTxHashToOutTxTracker
func HandleBroadcastError(err error, nonce, toChain, outTxHash string) (bool, bool) {
	if strings.Contains(err.Error(), "nonce too low") {
		log.Warn().Err(err).Msgf("nonce too low! this might be a unnecessary key-sign. increase re-try interval and awaits outTx confirmation")
		return false, false
	}
	if strings.Contains(err.Error(), "replacement transaction underpriced") {
		log.Warn().Err(err).Msgf("Broadcast replacement: nonce %s chain %s outTxHash %s", nonce, toChain, outTxHash)
		return false, false
	} else if strings.Contains(err.Error(), "already known") { // this is error code from QuickNode
		log.Warn().Err(err).Msgf("Broadcast duplicates: nonce %s chain %s outTxHash %s", nonce, toChain, outTxHash)
		return false, true // report to tracker, because there's possibilities a successful broadcast gets this error code
	}

	log.Error().Err(err).Msgf("Broadcast error: nonce %s chain %s outTxHash %s; retrying...", nonce, toChain, outTxHash)
	return true, false
}
