package zetaclient

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	OutboundTxSignCount = "zetaclient_Outbound_tx_sign_count"
)

type CoreObserver struct {
	bridge    *ZetaCoreBridge
	signerMap map[common.Chain]ChainSigner
	clientMap map[common.Chain]ChainClient
	metrics   *metrics.Metrics
	tss       *TSS
	logger    zerolog.Logger
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]ChainSigner, clientMap map[common.Chain]ChainClient, metrics *metrics.Metrics, tss *TSS) *CoreObserver {
	co := CoreObserver{}
	co.logger = log.With().Str("module", "CoreObserver").Logger()
	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap

	co.clientMap = clientMap
	co.metrics = metrics
	co.logger.Info().Msg("starting core observer")
	err := metrics.RegisterCounter(OutboundTxSignCount, "number of Outbound tx signed")
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
	log.Info().Msgf("monitorCore started by signer %s", myid)
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
						co.logger.Error().Err(err).Msgf("SetTSS fail %s", chain.String())
					}
				}

				// Keysign test: sanity test
				co.logger.Info().Msgf("test keysign...")
				_ = TestKeysign(co.tss.CurrentPubkey, co.tss.Server)
				co.logger.Info().Msg("test keysign finished. exit keygen loop. ")

				for _, chain := range config.ChainsEnabled {
					err = co.clientMap[chain].PostNonceIfNotRecorded()
					if err != nil {
						co.logger.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain.String())
					}
				}

				return
			}
		}()
		return
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
			//tStart := time.Now()
			sendList, err := co.bridge.GetAllPendingCctx()
			if err != nil {
				logger.Error().Err(err).Msg("error requesting sends from zetacore")
				continue
			}
			//logger.Info().Dur("elapsed", time.Since(tStart)).Msgf("GetAllPendingCctx %d", len(sendList))
			sendMap := SplitAndSortSendListByChain(sendList)

			// schedule sends
			for chain, sendList := range sendMap {
				chainName := common.ParseStringToObserverChain(chain)
				c := GetChainFromChainName(chainName)

				found := false
				for _, enabledChain := range GetSupportedChains() {
					if enabledChain.ChainId == c.ChainId {
						found = true
						break
					}
				}
				if !found {
					log.Warn().Msgf("chain %s is not enabled; skip scheduling", c.String())
					continue
				}
				if bn%10 == 0 {
					logger.Info().Msgf("outstanding %d CCTX's on chain %s: range [%d,%d]", len(sendList), chain, sendList[0].OutboundTxParams.OutboundTxTssNonce, sendList[len(sendList)-1].OutboundTxParams.OutboundTxTssNonce)
				}
				signer := co.signerMap[*c]
				chainClient := co.clientMap[*c]
				for idx, send := range sendList {
					ob, err := co.getTargetChainOb(send)
					if err != nil {
						logger.Error().Err(err).Msgf("getTargetChainOb fail %s", chain)
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
					included, _, err := ob.IsSendOutTxProcessed(send.Index, int(send.OutboundTxParams.OutboundTxTssNonce), send.OutboundTxParams.CoinType)
					if err != nil {
						logger.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", chain)
						continue
					}
					if included {
						logger.Info().Msgf("send outTx already included; do not schedule")
						continue
					}
					chain := GetTargetChain(send)
					outTxID := fmt.Sprintf("%s", chain, send.Index) // should be the outTxID?
					nonce := send.OutboundTxParams.OutboundTxTssNonce

					// FIXME: config this schedule; this value is for localnet fast testing
					if nonce%1 == bn%1 && !outTxMan.IsOutTxActive(outTxID) {
						outTxMan.StartTryProcess(outTxID)
						fmt.Printf("chain %s: Sign outtx %s with value %d\n", chain, send.Index, send.ZetaMint)
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
		if sends[i].OutboundTxParams.OutboundTxTssNonce > sends[i-1].OutboundTxParams.OutboundTxTssNonce+1000 {
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
			return sends[i].OutboundTxParams.OutboundTxTssNonce < sends[j].OutboundTxParams.OutboundTxTssNonce
		})
		start := trimSends(sends)
		sendMap[chain] = sends[start:]
		log.Debug().Msgf("chain %s, start %d, len %d, start nonce %d", chain, start, len(sendMap[chain]), sends[start].OutboundTxParams.OutboundTxTssNonce)
	}
	return sendMap
}

func GetTargetChain(send *types.CrossChainTx) string {
	if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		return send.OutboundTxParams.ReceiverChain
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		return send.InboundTxParams.SenderChain
	}
	return ""
}

func (co *CoreObserver) getTargetChainOb(send *types.CrossChainTx) (ChainClient, error) {
	chainStr := GetTargetChain(send)
	chainName := common.ParseStringToObserverChain(chainStr)
	c := GetChainFromChainName(chainName)
	//c, err := common.ParseChain(chainStr)
	//if err != nil {
	//	return nil, err
	//}
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
