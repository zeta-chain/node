package zetaclient

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"gitlab.com/thorchain/tss/go-tss/keygen"

	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"

	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/zeta-chain/zetacore/x/zetacore/types"

	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
)

const (
	OUTBOUND_TX_SIGN_COUNT = "zetaclient_outbound_tx_sign_count"
)

type CoreObserver struct {
	bridge    *ZetaCoreBridge
	signerMap map[common.Chain]*Signer
	clientMap map[common.Chain]*ChainObserver
	metrics   *metrics.Metrics
	tss       *TSS
	logger    zerolog.Logger
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]*Signer, clientMap map[common.Chain]*ChainObserver, metrics *metrics.Metrics, tss *TSS) *CoreObserver {
	co := CoreObserver{}
	co.logger = log.With().Str("module", "CoreObserver").Logger()
	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap

	co.clientMap = clientMap
	co.metrics = metrics

	err := metrics.RegisterCounter(OUTBOUND_TX_SIGN_COUNT, "number of outbound tx signed")
	if err != nil {
		co.logger.Error().Err(err).Msg("error registering counter")
	}

	return &co
}

func (co *CoreObserver) GetPromCounter(name string) (prom.Counter, error) {
	if cnt, found := metrics.Counters[name]; found {
		return cnt, nil
	} else {
		return nil, errors.New("counter not found")
	}
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
					_, err = co.bridge.SetTSS(chain, co.tss.Address().Hex(), co.tss.CurrentPubkey)
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

func (OutTxMan *OutTxProcessorManager) TimeInTryProcess(outTxID string) time.Duration {
	OutTxMan.mu.Lock()
	defer OutTxMan.mu.Unlock()
	if _, found := OutTxMan.outTxActive[outTxID]; found {
		return time.Since(OutTxMan.outTxStartTime[outTxID])
	} else {
		return 0
	}
}

// suicide whole zetaclient if keysign appears deadlocked.
func (OutTxMan *OutTxProcessorManager) StartMonitorHealth() {
	logger := OutTxMan.logger
	logger.Info().Msgf("StartMonitorHealth")
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		count := 0
		for outTxID, _ := range OutTxMan.outTxActive {
			if OutTxMan.TimeInTryProcess(outTxID).Minutes() > 5 {
				count++
			}
		}
		if count > 0 {
			logger.Warn().Msgf("Health: %d OutTx are more than 5min in process!", count)
		} else {
			logger.Info().Msgf("Monitor: healthy; numActiveProcessor %d", OutTxMan.numActiveProcessor)
		}
		if count > 10 {
			// suicide:
			logger.Error().Msgf("suicide zetaclient because keysign appears deadlocked; kill this process and the process supervisor should restart it")
			logger.Info().Msgf("numActiveProcessor: %d", OutTxMan.numActiveProcessor)
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
	var lastBlockNum uint64 = 0
	for range observeTicker.C {
		bn, err := co.bridge.GetZetaBlockHeight()
		if err != nil {
			logger.Error().Msg("GetZetaBlockHeight fail in startSendScheduler")
			continue
		}
		if bn > lastBlockNum { // we have a new block
			if bn%10 == 0 {
				logger.Info().Msgf("ZetaCore heart beat: %d", bn)
			}
			sendList, err := co.bridge.GetAllPendingSend()
			if err != nil {
				logger.Error().Err(err).Msg("error requesting sends from zetacore")
				continue
			}
			if len(sendList) > 0 && bn%5 == 0 {
				logger.Info().Msgf("#pending send: %d", len(sendList))
			}
			sendMap := splitAndSortSendListByChain(sendList)

			// schedule sends
			for chain, sendList := range sendMap {
				if bn%10 == 0 {
					logger.Info().Msgf("outstanding %d sends on chain %s: range [%d,%d]", len(sendList), chain, sendList[0].Nonce, sendList[len(sendList)-1].Nonce)
				}
				for idx, send := range sendList {
					ob, err := co.getTargetChainOb(send)
					if err != nil {
						logger.Error().Err(err).Msgf("getTargetChainOb fail %s", chain)
						continue
					}
					// update metrics
					if idx == 0 {
						pTxs, err := ob.GetPromGauge(metrics.PENDING_TXS)
						if err != nil {
							co.logger.Warn().Msgf("cannot get prometheus counter [%s]", metrics.PENDING_TXS)
							continue
						}
						pTxs.Set(float64(len(sendList)))
					}
					included, confirmed, err := ob.IsSendOutTxProcessed(send.Index, int(send.Nonce))
					if err != nil {
						logger.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", chain)
					}
					if included || confirmed {
						logger.Info().Msgf("send outTx already included; do not schedule")
						continue
					}
					chain := getTargetChain(send)
					outTxID := fmt.Sprintf("%s/%d", chain, send.Nonce)

					sinceBlock := int64(bn) - int64(send.FinalizedMetaHeight)
					// if there are many outstanding sends, then all first 20 has priority
					// otherwise, only the first one has priority

					if isScheduled(sinceBlock, idx < 40) && !outTxMan.IsOutTxActive(outTxID) {
						outTxMan.StartTryProcess(outTxID)
						go co.TryProcessOutTx(send, sinceBlock, outTxMan)
					}
					if idx > 50 { // only look at 50 sends per chain
						break
					}
				}
			}
			// update last processed block number
			lastBlockNum = bn
		}

	}
}

func (co *CoreObserver) TryProcessOutTx(send *types.Send, sinceBlock int64, outTxMan *OutTxProcessorManager) {
	chain := getTargetChain(send)
	outTxID := fmt.Sprintf("%s/%d", chain, send.Nonce)

	logger := co.logger.With().
		Str("sendHash", send.Index).
		Str("outTxID", outTxID).
		Int64("sinceBlock", sinceBlock).Logger()
	logger.Info().Msgf("start processing outTxID %s", outTxID)
	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()

	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	amount, ok := new(big.Int).SetString(send.ZetaMint, 0)
	if !ok {
		logger.Error().Msg("error converting MBurnt to big.Int")
		return
	}

	var to ethcommon.Address
	var err error
	var toChain common.Chain
	if send.Status == types.SendStatus_PendingRevert {
		to = ethcommon.HexToAddress(send.Sender)
		toChain, err = common.ParseChain(send.SenderChain)
		logger.Info().Msgf("Abort: reverting inbound")
	} else if send.Status == types.SendStatus_PendingOutbound {
		to = ethcommon.HexToAddress(send.Receiver)
		toChain, err = common.ParseChain(send.ReceiverChain)
	}
	if err != nil {
		logger.Error().Err(err).Msg("ParseChain fail; skip")
		return
	}

	// Early return if the send is already processed
	included, confirmed, _ := co.clientMap[toChain].IsSendOutTxProcessed(send.Index, int(send.Nonce))
	if included || confirmed {
		logger.Info().Msgf("sendHash already processed; exit signer")
		return
	}

	signer := co.signerMap[toChain]
	message, err := base64.StdEncoding.DecodeString(send.Message)
	if err != nil {
		logger.Err(err).Msgf("decode send.Message %s error", send.Message)
	}

	gasLimit := send.GasLimit
	if gasLimit < 50_000 {
		gasLimit = 50_000
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", send.GasLimit, gasLimit)
	}
	if gasLimit > 1_000_000 {
		gasLimit = 1_000_000
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", send.GasLimit, gasLimit)
	}

	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, amount, to.Hex(), send.Nonce, send.FinalizedMetaHeight)
	sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		logger.Error().Err(err).Msgf("decode sendHash %s error", send.Index)
		return
	}
	var sendhash [32]byte
	copy(sendhash[:32], sendHash[:32])
	gasprice, ok := new(big.Int).SetString(send.GasPrice, 10)
	if !ok {
		logger.Error().Err(err).Msgf("cannot convert gas price  %s ", send.GasPrice)
		return
	}
	// use 33% higher gas price for timely confirmation
	gasprice = gasprice.Mul(gasprice, big.NewInt(4))
	gasprice = gasprice.Div(gasprice, big.NewInt(3))
	var tx *ethtypes.Transaction

	srcChainID := config.Chains[send.SenderChain].ChainID
	if send.Status == types.SendStatus_PendingRevert {
		logger.Info().Msgf("SignRevertTx: %s => %s, nonce %d", send.SenderChain, toChain, send.Nonce)
		toChainID := config.Chains[send.ReceiverChain].ChainID
		tx, err = signer.SignRevertTx(ethcommon.HexToAddress(send.Sender), srcChainID, to.Bytes(), toChainID, amount, gasLimit, message, sendhash, send.Nonce, gasprice)
	} else if send.Status == types.SendStatus_PendingOutbound {
		logger.Info().Msgf("SignOutboundTx: %s => %s, nonce %d", send.SenderChain, toChain, send.Nonce)
		tx, err = signer.SignOutboundTx(ethcommon.HexToAddress(send.Sender), srcChainID, to, amount, gasLimit, message, sendhash, send.Nonce, gasprice)
	}

	if err != nil {
		logger.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %s", send.Nonce, send.ReceiverChain)
		return
	} else {
		logger.Info().Msgf("Keysign success: %s => %s, nonce %d", send.SenderChain, toChain, send.Nonce)
	}
	cnt, err := co.GetPromCounter(OUTBOUND_TX_SIGN_COUNT)
	if err != nil {
		log.Error().Err(err).Msgf("GetPromCounter error")
	} else {
		cnt.Inc()
	}
	if tx != nil {
		outTxHash := tx.Hash().Hex()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, send.Nonce, outTxHash, myid)
		if myid == send.Signers[send.Broadcaster] || myid == send.Signers[int(send.Broadcaster+1)%len(send.Signers)] {
			backOff := 1000 * time.Millisecond
			// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
			for i := 0; i < 5; i++ {
				logger.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, send.Nonce, i)
				// #nosec G404 randomness is not a security issue here
				time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) //random delay to avoid sychronized broadcast
				err := signer.Broadcast(tx)
				if err != nil {
					log.Warn().Err(err).Msgf("OutTx Broadcast error")
					retry, report := HandleBroadcastError(err, strconv.FormatUint(send.Nonce, 10), toChain.String(), outTxHash)
					if report {
						zetaHash, err := co.bridge.AddTxHashToWatchlist(toChain.String(), tx.Nonce(), outTxHash)
						if err != nil {
							logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
						}
						logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
					}
					if !retry {
						break
					}
					backOff *= 2
					continue
				}
				logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
				zetaHash, err := co.bridge.AddTxHashToWatchlist(toChain.String(), tx.Nonce(), outTxHash)
				if err != nil {
					logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
				}
				logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
				break // successful broadcast; no need to retry
			}

		}
	}

}

func isScheduled(diff int64, priority bool) bool {
	d := diff - 1
	if d < 0 {
		return false
	}
	if priority {
		return d%15 == 0
	}
	if d < 100 && d%20 == 0 {
		return true
	} else if d >= 100 && d%100 == 0 { // after 100 blocks, schedule once per 100 blocks
		return true
	}
	return false
}

func splitAndSortSendListByChain(sendList []*types.Send) map[string][]*types.Send {
	sendMap := make(map[string][]*types.Send)
	for _, send := range sendList {
		targetChain := getTargetChain(send)
		if targetChain == "" {
			continue
		}
		if _, found := sendMap[targetChain]; !found {
			sendMap[targetChain] = make([]*types.Send, 0)
		}
		sendMap[targetChain] = append(sendMap[targetChain], send)
	}
	for chain, sends := range sendMap {
		sort.Slice(sends, func(i, j int) bool {
			return sends[i].Nonce < sends[j].Nonce
		})
		sendMap[chain] = sends
	}
	return sendMap
}

func getTargetChain(send *types.Send) string {
	if send.Status == types.SendStatus_PendingOutbound {
		return send.ReceiverChain
	} else if send.Status == types.SendStatus_PendingRevert {
		return send.SenderChain
	}
	return ""
}

func (co *CoreObserver) getTargetChainOb(send *types.Send) (*ChainObserver, error) {
	chainStr := getTargetChain(send)
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

// returns whether to retry in a few seconds, and whether to report via AddTxHashToWatchlist
func HandleBroadcastError(err error, nonce, toChain, outTxHash string) (bool, bool) {
	if strings.Contains(err.Error(), "nonce too low") {
		log.Warn().Err(err).Msgf("nonce too low! this might be a unnecessary keysign. increase re-try interval and awaits outTx confirmation")
		return false, false
	}
	if strings.Contains(err.Error(), "replacement transaction underpriced") {
		log.Warn().Err(err).Msgf("Broadcast replacement: nonce %s chain %s outTxHash %s", nonce, toChain, outTxHash)
		return false, false
	} else if strings.Contains(err.Error(), "already known") { // this is error code from QuickNode
		log.Warn().Err(err).Msgf("Broadcast duplicates: nonce %s chain %s outTxHash %s", nonce, toChain, outTxHash)
		return false, true // report to tracker, because there's possibilities a successful broadcast gets this error code
	}

	log.Error().Err(err).Msgf("Broadcast error: nonce %s chain %s outTxHash %s; retring...", nonce, toChain, outTxHash)
	return true, false
}
