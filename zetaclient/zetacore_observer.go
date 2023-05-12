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
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]ChainSigner, clientMap map[common.Chain]ChainClient, metrics *metrics.Metrics, tss *TSS, logger zerolog.Logger, cfg *config.Config) *CoreObserver {
	co := CoreObserver{}
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
}

// ZetaCore block is heart beat; each block we schedule some send according to
// retry schedule. ses
func (co *CoreObserver) startSendScheduler() {
	outTxMan := NewOutTxProcessorManager(co.logger.ChainLogger)
	go outTxMan.StartMonitorHealth()
	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum int64
	for range observeTicker.C {
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
			//tStart := time.Now()
			sendList, err := co.bridge.GetAllPendingCctx()
			if err != nil {
				co.logger.ZetaChainWatcher.Error().Err(err).Msg("error requesting sends from zetacore")
				continue
			}
			//logger.Info().Dur("elapsed", time.Since(tStart)).Msgf("GetAllPendingCctx %d", len(sendList))
			sendMap := SplitAndSortSendListByChain(sendList)

			// schedule sends
			for chain, sendList := range sendMap {
				chainName := common.ParseChainName(chain)
				c := common.GetChainFromChainName(chainName)

				found := false
				for _, enabledChain := range GetSupportedChains() {
					if enabledChain.ChainId == c.ChainId {
						found = true
						break
					}
				}
				if !found {
					co.logger.ZetaChainWatcher.Warn().Msgf("chain %s is not enabled; skip scheduling", c.String())
					continue
				}
				if bn%10 == 0 {
					co.logger.ZetaChainWatcher.Info().Msgf("outstanding %d CCTX's on chain %s: range [%d,%d]", len(sendList), chain, sendList[0].GetCurrentOutTxParam().OutboundTxTssNonce, sendList[len(sendList)-1].GetCurrentOutTxParam().OutboundTxTssNonce)
				}
				signer := co.signerMap[*c]
				chainClient := co.clientMap[*c]
				cnt := 0
				for _, send := range sendList {
					ob, err := co.getTargetChainOb(send)
					if err != nil {
						co.logger.ZetaChainWatcher.Error().Err(err).Msgf("getTargetChainOb fail %s", chain)
						continue
					}
					// update metrics
					//if idx == 0 {
					//	pTxs, err := ob.GetPromGauge(metrics.PendingTxs)
					//	if err != nil {
					//		co.logger.Warn().Msgf("cannot get prometheus counter [%s]", metrics.PendingTxs)
					//	} else {
					//		pTxs.Set(float64(len(sendList)))
					//	}
					//}
					// Monitor Core Logger for OutboundTxTssNonce
					included, _, err := ob.IsSendOutTxProcessed(send.Index, int(send.GetCurrentOutTxParam().OutboundTxTssNonce), send.GetCurrentOutTxParam().CoinType, co.logger.ZetaChainWatcher)
					if err != nil {
						co.logger.ZetaChainWatcher.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", chain)
						continue
					}
					if included {
						co.logger.ZetaChainWatcher.Info().Msgf("send outTx already included; do not schedule")
						continue
					}
					chain := GetTargetChain(send)
					nonce := send.GetCurrentOutTxParam().OutboundTxTssNonce
					outTxID := fmt.Sprintf("%s-%d-%d", send.Index, send.GetCurrentOutTxParam().ReceiverChainId, nonce) // should be the outTxID?

					// FIXME: config this schedule; this value is for localnet fast testing
					if bn >= math.MaxInt64 {
						continue
					}
					currentHeight := uint64(bn)
					if nonce%10 == currentHeight%10 {
						cnt++
						outTxMan.StartTryProcess(outTxID)
						co.logger.ZetaChainWatcher.Debug().Msgf("chain %s: Sign outtx %s with value %d\n", chain, send.Index, send.GetCurrentOutTxParam().Amount)
						go signer.TryProcessOutTx(send, outTxMan, outTxID, chainClient, co.bridge)
					}
					if cnt == 3 {
						break
					}
				}
			}
			// update last processed block number
			lastBlockNum = bn
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
		targetChain := GetTargetChain(send)
		if targetChain == "" {
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

func GetTargetChain(send *types.CrossChainTx) string {
	chainID := send.GetCurrentOutTxParam().ReceiverChainId
	return common.GetChainFromChainID(chainID).GetChainName().String()
}

func (co *CoreObserver) getTargetChainOb(send *types.CrossChainTx) (ChainClient, error) {
	chainStr := GetTargetChain(send)
	chainName := common.ParseChainName(chainStr)
	c := common.GetChainFromChainName(chainName)

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
