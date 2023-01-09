package zetaclient

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"gitlab.com/thorchain/tss/go-tss/keygen"

	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"

	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/zeta-chain/zetacore/x/crosschain/types"

	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
)

const (
	OutboundTxSignCount = "zetaclient_outbound_tx_sign_count"
)

type CoreObserver struct {
	bridge    *ZetaCoreBridge
	signerMap map[common.Chain]*EVMSigner
	clientMap map[common.Chain]ChainClient
	metrics   *metrics.Metrics
	tss       *TSS
	logger    zerolog.Logger
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]*EVMSigner, clientMap map[common.Chain]ChainClient, metrics *metrics.Metrics, tss *TSS) *CoreObserver {
	co := CoreObserver{}
	co.logger = log.With().Str("module", "CoreObserver").Logger()
	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap

	co.clientMap = clientMap
	co.metrics = metrics

	err := metrics.RegisterCounter(OutboundTxSignCount, "number of outbound tx signed")
	if err != nil {
		co.logger.Error().Err(err).Msg("error registering counter")
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
	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	log.Info().Msgf("MonitorCore started by signer %s", myid)
	go co.startSendScheduler()

	noKeygen := os.Getenv("DISABLE_TSS_KEYGEN")
	if noKeygen == "" {
		go co.keygenObserve()
	}
}

func (co *CoreObserver) keygenObserve() {
	log.Info().Msgf("keygen observe started")
	observeTicker := time.NewTicker(2 * time.Second)
	for range observeTicker.C {
		kg, err := co.bridge.GetKeyGen()
		if err != nil {
			continue
		}
		bn, _ := co.bridge.GetZetaBlockHeight()
		if bn != kg.BlockNumber {
			continue
		}

		go func() {
			for {
				log.Info().Msgf("Detected KeyGen, initiate keygen at blocknumm %d, # signers %d", kg.BlockNumber, len(kg.Pubkeys))
				var req keygen.Request
				req = keygen.NewRequest(kg.Pubkeys, int64(kg.BlockNumber), "0.14.0")
				res, err := co.tss.Server.Keygen(req)
				if err != nil || res.Status != tsscommon.Success {
					co.logger.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
					continue
				}
				// Keygen succeed! Report TSS address
				co.logger.Info().Msgf("Keygen success! keygen response: %v...", res)
				err = co.tss.InsertPubKey(res.PubKey)
				if err != nil {
					co.logger.Error().Msgf("InsertPubKey fail")
					continue
				}
				co.tss.CurrentPubkey = res.PubKey

				for _, chain := range config.ChainsEnabled {
					_, err = co.bridge.SetTSS(chain, co.tss.EVMAddress().Hex(), co.tss.CurrentPubkey)
					if err != nil {
						co.logger.Error().Err(err).Msgf("SetTSS fail %s", chain)
					}
				}

				// Keysign test: sanity test
				co.logger.Info().Msgf("test keysign...")
				_ = TestKeysign(co.tss.CurrentPubkey, co.tss.Server)
				co.logger.Info().Msg("test keysign finished. exit keygen loop. ")

				for _, chain := range config.ChainsEnabled {
					err = co.clientMap[chain].PostNonceIfNotRecorded()
					if err != nil {
						co.logger.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain)
					}
				}

				return
			}
		}()
		return
	}
}

type OutTxProcessorManager struct {
	outTxStartTime     map[string]time.Time
	outTxEndTime       map[string]time.Time
	outTxActive        map[string]struct{}
	mu                 sync.Mutex
	logger             zerolog.Logger
	numActiveProcessor int64
}

func NewOutTxProcessorManager() *OutTxProcessorManager {
	return &OutTxProcessorManager{
		outTxStartTime:     make(map[string]time.Time),
		outTxEndTime:       make(map[string]time.Time),
		outTxActive:        make(map[string]struct{}),
		mu:                 sync.Mutex{},
		logger:             log.With().Str("module", "OutTxProcessorManager").Logger(),
		numActiveProcessor: 0,
	}
}

func (outTxMan *OutTxProcessorManager) StartTryProcess(outTxID string) {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	outTxMan.outTxStartTime[outTxID] = time.Now()
	outTxMan.outTxActive[outTxID] = struct{}{}
	outTxMan.numActiveProcessor++
	outTxMan.logger.Info().Msgf("StartTryProcess %s, numActiveProcessor %d", outTxID, outTxMan.numActiveProcessor)
}

func (outTxMan *OutTxProcessorManager) EndTryProcess(outTxID string) {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	outTxMan.outTxEndTime[outTxID] = time.Now()
	delete(outTxMan.outTxActive, outTxID)
	outTxMan.numActiveProcessor--
	outTxMan.logger.Info().Msgf("EndTryProcess %s, numActiveProcessor %d, time elapsed %s", outTxID, outTxMan.numActiveProcessor, time.Since(outTxMan.outTxStartTime[outTxID]))
}

func (outTxMan *OutTxProcessorManager) IsOutTxActive(outTxID string) bool {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	_, found := outTxMan.outTxActive[outTxID]
	return found
}

func (outTxMan *OutTxProcessorManager) TimeInTryProcess(outTxID string) time.Duration {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	if _, found := outTxMan.outTxActive[outTxID]; found {
		return time.Since(outTxMan.outTxStartTime[outTxID])
	}
	return 0
}

// suicide whole zetaclient if keysign appears deadlocked.
func (outTxMan *OutTxProcessorManager) StartMonitorHealth() {
	logger := outTxMan.logger
	logger.Info().Msgf("StartMonitorHealth")
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		count := 0
		for outTxID := range outTxMan.outTxActive {
			if outTxMan.TimeInTryProcess(outTxID).Minutes() > 2 {
				count++
			}
		}
		if count > 0 {
			logger.Warn().Msgf("Health: %d OutTx are more than 2min in process!", count)
		} else {
			logger.Info().Msgf("Monitor: healthy; numActiveProcessor %d", outTxMan.numActiveProcessor)
		}
		if count > 10 {
			// suicide:
			logger.Error().Msgf("suicide zetaclient because keysign appears deadlocked; kill this process and the process supervisor should restart it")
			logger.Info().Msgf("numActiveProcessor: %d", outTxMan.numActiveProcessor)
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}
}

// ZetaCore block is heart beat; each block we schedule some send according to
// retry schedule.
func (co *CoreObserver) startSendScheduler() {
	logger := co.logger.With().Str("module", "SendScheduler").Logger()
	outTxMan := NewOutTxProcessorManager()
	go outTxMan.StartMonitorHealth()

	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum uint64
	for range observeTicker.C {
		bn, err := co.bridge.GetZetaBlockHeight()
		if err != nil {
			logger.Error().Msg("GetZetaBlockHeight fail in startSendScheduler")
			continue
		}
		if lastBlockNum == 0 {
			lastBlockNum = bn - 1
		}
		if bn > lastBlockNum { // we have a new block
			bn = lastBlockNum + 1
			if bn%10 == 0 {
				logger.Info().Msgf("ZetaCore heart beat: %d", bn)
			}
			tStart := time.Now()
			sendList, err := co.bridge.GetAllPendingCctx()
			if err != nil {
				logger.Error().Err(err).Msg("error requesting sends from zetacore")
				continue
			}
			logger.Info().Dur("elapsed", time.Since(tStart)).Msgf("GetAllPendingCctx %d", len(sendList))
			sendMap := SplitAndSortSendListByChain(sendList)

			// schedule sends
			for chain, sendList := range sendMap {
				c, _ := common.ParseChain(chain)
				found := false
				for _, enabledChain := range config.ChainsEnabled {
					if enabledChain == c {
						found = true
						break
					}
				}
				if !found {
					log.Warn().Msgf("chain %s is not enabled; skip scheduling", chain)
					continue
				}
				if bn%10 == 0 {
					logger.Info().Msgf("outstanding %d CCTX's on chain %s: range [%d,%d]", len(sendList), chain, sendList[0].OutBoundTxParams.OutBoundTxTSSNonce, sendList[len(sendList)-1].OutBoundTxParams.OutBoundTxTSSNonce)
				}
				signer := co.signerMap[c]
				chainClient := co.clientMap[c]
				for idx, send := range sendList {
					ob, err := co.getTargetChainOb(send)
					if err != nil {
						logger.Error().Err(err).Msgf("getTargetChainOb fail %s", chain)
						continue
					}
					// update metrics
					if idx == 0 {
						pTxs, err := ob.GetPromGauge(metrics.PendingTxs)
						if err != nil {
							co.logger.Warn().Msgf("cannot get prometheus counter [%s]", metrics.PendingTxs)
							continue
						}
						pTxs.Set(float64(len(sendList)))
					}
					included, _, err := ob.IsSendOutTxProcessed(send)
					if err != nil {
						logger.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", chain)
						continue
					}
					if included {
						logger.Info().Msgf("send outTx already included; do not schedule")
						continue
					}
					chain := GetTargetChain(send)
					outTxID := fmt.Sprintf("%s-%d", chain, send.OutBoundTxParams.OutBoundTxTSSNonce)
					nonce := send.OutBoundTxParams.OutBoundTxTSSNonce

					if nonce%20 == bn%20 && !outTxMan.IsOutTxActive(outTxID) {
						outTxMan.StartTryProcess(outTxID)
						go signer.TryProcessOutTx(send, outTxMan, chainClient, co.bridge)
					}
					if idx > 60 { // only look at 50 sends per chain
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
		if sends[i].OutBoundTxParams.OutBoundTxTSSNonce > sends[i-1].OutBoundTxParams.OutBoundTxTSSNonce+1000 {
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
			return sends[i].OutBoundTxParams.OutBoundTxTSSNonce < sends[j].OutBoundTxParams.OutBoundTxTSSNonce
		})
		start := trimSends(sends)
		sendMap[chain] = sends[start:]
		log.Info().Msgf("chain %s, start %d, len %d, start nonce %d", chain, start, len(sendMap[chain]), sends[start].OutBoundTxParams.OutBoundTxTSSNonce)
	}
	return sendMap
}

func GetTargetChain(send *types.CrossChainTx) string {
	if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		return send.OutBoundTxParams.ReceiverChain
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		return send.InBoundTxParams.SenderChain
	}
	return ""
}

func (co *CoreObserver) getTargetChainOb(send *types.CrossChainTx) (ChainClient, error) {
	chainStr := GetTargetChain(send)
	c, err := common.ParseChain(chainStr)
	if err != nil {
		return nil, err
	}
	chainOb, found := co.clientMap[c]
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
