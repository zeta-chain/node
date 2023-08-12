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

	"github.com/zeta-chain/zetacore/x/crosschain/types"
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
	metrics   *metrics.Metrics
	tss       *TSS
	logger    ZetaCoreLog
	cfg       *config.Config
	ts        *TelemetryServer
	stop      chan struct{}
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]ChainSigner, clientMap map[common.Chain]ChainClient, metrics *metrics.Metrics, tss *TSS, logger zerolog.Logger, cfg *config.Config, ts *TelemetryServer) *CoreObserver {
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

	return &co
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
						if c == nil {
							co.logger.ZetaChainWatcher.Error().Msg("chain nil")
							continue
						}
						signer := co.signerMap[*c]
						chainClient := co.clientMap[*c]
						sendList, err := co.bridge.GetAllPendingCctx(uint64(c.ChainId))
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("failed to GetAllPendingCctx for chain %s", c.ChainName.String())
							continue
						}
						res, err := co.bridge.GetAllOutTxTrackerByChain(*c)
						if err != nil {
							co.logger.ZetaChainWatcher.Warn().Err(err).Msgf("failed to GetAllOutTxTrackerByChain for chain %s", c.ChainName.String())
						}
						trackerMap := make(map[uint64]bool)
						for _, v := range res {
							trackerMap[v.Nonce] = true
						}

						for idx, send := range sendList {
							if send.GetCurrentOutTxParam().ReceiverChainId != c.ChainId {
								log.Warn().Msgf("mismatch chainid: want %d, got %d", c.ChainId, send.GetCurrentOutTxParam().ReceiverChainId)
								continue
							}
							ob, err := co.getTargetChainOb(send)
							if err != nil {
								co.logger.ZetaChainWatcher.Error().Err(err).Msgf("getTargetChainOb fail %s", c.ChainName)
								continue
							}

							// Monitor Core Logger for OutboundTxTssNonce
							included, _, err := ob.IsSendOutTxProcessed(send.Index, int(send.GetCurrentOutTxParam().OutboundTxTssNonce), send.GetCurrentOutTxParam().CoinType, co.logger.ZetaChainWatcher)
							if err != nil {
								co.logger.ZetaChainWatcher.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", c.ChainName)
								continue
							}
							if included {
								co.logger.ZetaChainWatcher.Info().Msgf("send outTx already included; do not schedule")
								continue
							}
							chain, err := GetTargetChain(send)
							if err != nil {
								co.logger.ZetaChainWatcher.Error().Err(err).Msgf("GetTargetChain fail , Chain ID : %s", c.ChainName)
								continue
							}
							nonce := send.GetCurrentOutTxParam().OutboundTxTssNonce
							outTxID := fmt.Sprintf("%s-%d-%d", send.Index, send.GetCurrentOutTxParam().ReceiverChainId, nonce) // should be the outTxID?

							// FIXME: config this schedule; this value is for localnet fast testing
							if bn >= math.MaxInt64 {
								continue
							}
							currentHeight := uint64(bn)
							var interval uint64
							var lookahead int64
							// FIXME: fix these ugly type switches and conversions
							switch v := ob.(type) {
							case *EVMChainClient:
								interval = uint64(v.GetChainConfig().CoreParams.OutboundTxScheduleInterval)
								lookahead = v.GetChainConfig().CoreParams.OutboundTxScheduleLookahead
							case *BitcoinChainClient:
								interval = uint64(v.GetChainConfig().CoreParams.OutboundTxScheduleInterval)
								lookahead = v.GetChainConfig().CoreParams.OutboundTxScheduleLookahead
							default:
								co.logger.ZetaChainWatcher.Error().Msgf("unknown ob type on chain %s: type %T", chain, ob)
								continue
							}

							// determining critical outtx; if it satisfies following criteria
							// 1. it's the first pending outtx for this chain
							// 2. the following 5 nonces have been in tracker
							isCritical := false
							criticalInterval := uint64(10) // for critical pending outTx we reduce re-try interval
							if nonce%criticalInterval == currentHeight%criticalInterval && idx < 2 {
								numNoncesInTracker := 0
								for i := nonce + 1; i <= nonce+10; i++ {
									if _, found := trackerMap[i]; found {
										numNoncesInTracker++
									}
								}
								if numNoncesInTracker >= 7 {
									isCritical = true
								}
							}

							if (isCritical || nonce%interval == currentHeight%interval) && !outTxMan.IsOutTxActive(outTxID) {
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

func GetTargetChain(send *types.CrossChainTx) (string, error) {
	chainID := send.GetCurrentOutTxParam().ReceiverChainId
	chain := common.GetChainFromChainID(chainID)
	if chain == nil {
		return "", fmt.Errorf("chain %d not found", chainID)
	}
	return chain.GetChainName().String(), nil
}

func (co *CoreObserver) getTargetChainOb(send *types.CrossChainTx) (ChainClient, error) {
	chainStr, err := GetTargetChain(send)
	if err != nil {
		return nil, fmt.Errorf("chain %d not found", send.GetCurrentOutTxParam().ReceiverChainId)
	}
	chainName := common.ParseChainName(chainStr)
	c := common.GetChainFromChainName(chainName)
	if c == nil {
		return nil, fmt.Errorf("chain %s not found", chainName)
	}
	chainOb, found := co.clientMap[*c]
	if !found {
		return nil, fmt.Errorf("chain %s not found", c)
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
