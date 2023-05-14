package zetaclient

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"gitlab.com/thorchain/tss/go-tss/p2p"
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
	Tss       *TSS
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

	co.Tss = tss
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

// returns map(protocolID -> count); map(connID -> count)
func countActiveStreams(n network.Network) (map[string]int, map[string]int, int) {
	count := 0
	conns := n.Conns()
	protocolCount := make(map[string]int)
	connCount := make(map[string]int)
	for _, conn := range conns {
		count += len(conn.GetStreams())
		for _, stream := range conn.GetStreams() {
			protocolCount[string(stream.Protocol())]++
		}
		connCount[string(conn.ID())] += len(conn.GetStreams())
	}
	return protocolCount, connCount, count
}

var joinPartyProtocolWithLeader protocol.ID = "/p2p/join-party-leader"
var TSSProtocolID protocol.ID = "/p2p/tss"

func releaseAllStreams(n network.Network, streamMgr *p2p.StreamMgr) int {
	streams := streamMgr.UnusedStreams
	numKeys := len(streams)
	lenMap := make(map[string]int)
	for k, v := range streams {
		lenMap[k] = len(v)
	}
	log.Warn().Msgf("analyzing StreamMgr: %d keys; ", numKeys)
	log.Warn().Msgf("StreamMgr statistics by msgID: %v", lenMap)
	conns := n.Conns()
	cnt := 0
	for _, conn := range conns {
		for _, stream := range conn.GetStreams() {
			if stream.Protocol() == joinPartyProtocolWithLeader || stream.Protocol() == TSSProtocolID {
				stream.Reset()
				cnt++
			}
		}
	}
	return cnt
}

// ZetaCore block is heart beat; each block we schedule some send according to
// retry schedule. ses
func (co *CoreObserver) startSendScheduler() {
	outTxMan := NewOutTxProcessorManager(co.logger.ChainLogger)
	go outTxMan.StartMonitorHealth()
	observeTicker := time.NewTicker(1 * time.Second)
	var lastBlockNum int64
	zblockToProcessedNonce := make(map[int64]int64)
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
				if len(sendList) > 0 {
					co.logger.ZetaChainWatcher.Info().Msgf("outstanding %d CCTX's on chain %s: range [%d,%d]", len(sendList), chain, sendList[0].GetCurrentOutTxParam().OutboundTxTssNonce, sendList[len(sendList)-1].GetCurrentOutTxParam().OutboundTxTssNonce)
				} else {
					continue
				}
				signer := co.signerMap[*c]
				chainClient := co.clientMap[*c]
				cnt := 0
				maxCnt := 4
				safeMode := true // by default, be cautious and only send 1 tx per block
				if len(sendList) > 0 {
					lastProcessedNonce := int64(sendList[0].GetCurrentOutTxParam().OutboundTxTssNonce) - 1
					zblockToProcessedNonce[bn] = lastProcessedNonce
					// if for 10 blocks there is no progress, then wind down the maxCnt (lookahead)
					if nonce1, found := zblockToProcessedNonce[bn-10]; found {
						if nonce1 < lastProcessedNonce && outTxMan.numActiveProcessor < 10 {
							safeMode = false
						}
					}
					co.logger.ZetaChainWatcher.Info().Msgf("20 blocks outbound tx processing rate: %.2f", float64(lastProcessedNonce-zblockToProcessedNonce[bn-20])/20.0)
					co.logger.ZetaChainWatcher.Info().Msgf("100 blocks outbound tx processing rate: %.2f", float64(lastProcessedNonce-zblockToProcessedNonce[bn-100])/100.0)
					co.logger.ZetaChainWatcher.Info().Msgf("since block 0 outbound tx processing rate: %.2f", float64(lastProcessedNonce)/(1.0*float64(bn)))
				}
				streamMgr := co.Tss.Server.P2pCommunication.StreamMgr

				host := co.Tss.Server.P2pCommunication.GetHost()
				pCount, cCount, numStreams := countActiveStreams(host.Network())
				co.logger.ZetaChainWatcher.Info().Msgf("numStreams: %d; protocol: %+v; conn: %+v", numStreams, pCount, cCount)
				if outTxMan.numActiveProcessor == 0 {
					co.logger.ZetaChainWatcher.Warn().Msgf("no active outbound tx processor; safeMode: %v", safeMode)
					numStreamsReleased := releaseAllStreams(host.Network(), streamMgr)
					co.logger.ZetaChainWatcher.Warn().Msgf("released %d streams", numStreamsReleased)
				}

				for _, send := range sendList {
					ob, err := co.getTargetChainOb(send)
					if err != nil {
						co.logger.ZetaChainWatcher.Error().Err(err).Msgf("getTargetChainOb fail %s", chain)
						continue
					}
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
					if nonce%10 == currentHeight%10 && !outTxMan.IsOutTxActive(outTxID) {
						if safeMode && nonce != sendList[0].GetCurrentOutTxParam().OutboundTxTssNonce {
							break
						}
						cnt++
						outTxMan.StartTryProcess(outTxID)
						co.logger.ZetaChainWatcher.Debug().Msgf("chain %s: Sign outtx %s with value %d\n", chain, send.Index, send.GetCurrentOutTxParam().Amount)
						go signer.TryProcessOutTx(send, outTxMan, outTxID, chainClient, co.bridge)
					}
					if cnt == maxCnt {
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
