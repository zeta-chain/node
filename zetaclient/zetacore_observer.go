package zetaclient

import (
	"fmt"
	"math"
	"sort"
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
						signer := co.signerMap[*c]
						chainClient := co.clientMap[*c]
						sendList, err := co.bridge.GetAllPendingCctx(uint64(c.ChainId))
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("failed to GetAllPendingCctx for chain %s", c.ChainName.String())
							continue
						}
						ob, err := co.getTargetChainOb(c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("getTargetChainOb fail, Chain ID: %s", c.ChainName)
							continue
						}
						chain, err := common.GetChainNameFromChainID(c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("GetTargetChain fail, Chain ID: %s", c.ChainName)
							continue
						}

						// Any necessary preparation work (e.g. update pending sends)
						ob.PreSendSchedule(sendList)

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
							interval := uint64(ob.GetCoreParameters().OutboundTxScheduleInterval)
							lookahead := uint64(ob.GetCoreParameters().OutboundTxScheduleLookahead)
							if nonce%interval == currentHeight%interval && !outTxMan.IsOutTxActive(outTxID) {
								outTxMan.StartTryProcess(outTxID)
								co.logger.ZetaChainWatcher.Debug().Msgf("chain %s: Sign outtx %s with value %d sats", chain, send.Index, params.Amount)
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

// trim "bogus" pending sends that are not actually pending
// input sends must be sorted by nonce ascending
func trimSends(sends []*types.CrossChainTx) int {
	start := 0
	for i := len(sends) - 1; i >= 1; i-- {
		// from right to left, if there's a big hole, then before the gap are probably
		// bogus "pending" sends that are already processed but not yet confirmed.
		if sends[i].GetCurrentOutTxParam().OutboundTxTssNonce > sends[i-1].GetCurrentOutTxParam().OutboundTxTssNonce+1000 {
			start = i
			break
		}
	}
	return start
}

func SplitAndSortSendListByChain(sendList []*types.CrossChainTx) map[string][]*types.CrossChainTx {
	sendMap := make(map[string][]*types.CrossChainTx)
	for _, send := range sendList {
		targetChain, err := common.GetChainNameFromChainID(send.GetCurrentOutTxParam().ReceiverChainId)
		if targetChain == "" || err != nil {
			continue
		}
		if _, found := sendMap[targetChain]; !found {
			sendMap[targetChain] = make([]*types.CrossChainTx, 0)
		}
		sendMap[targetChain] = append(sendMap[targetChain], send)
	}
	for chain, sends := range sendMap {
		sort.Slice(sends, func(i, j int) bool {
			return sends[i].GetCurrentOutTxParam().OutboundTxTssNonce < sends[j].GetCurrentOutTxParam().OutboundTxTssNonce
		})
		start := trimSends(sends)
		sendMap[chain] = sends[start:]
		log.Debug().Msgf("chain %s, start %d, len %d, start nonce %d", chain, start, len(sendMap[chain]), sends[start].GetCurrentOutTxParam().OutboundTxTssNonce)
	}
	return sendMap
}

func (co *CoreObserver) getTargetChainOb(chainID int64) (ChainClient, error) {
	chainStr, err := common.GetChainNameFromChainID(chainID)
	if err != nil {
		return nil, fmt.Errorf("chain %d not found", chainID)
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
