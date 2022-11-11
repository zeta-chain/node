package zetaclient

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/big"
	"math/rand"
	"net/http"
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
	OutboundTxSignCount = "zetaclient_outbound_tx_sign_count"
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
	go co.CleanUpCommand()

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

func (outTxMan *OutTxProcessorManager) StartTryProcess(outTxID string) int64 {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	outTxMan.outTxStartTime[outTxID] = time.Now()
	outTxMan.outTxActive[outTxID] = struct{}{}
	outTxMan.numActiveProcessor++
	outTxMan.logger.Info().Msgf("StartTryProcess %s, numActiveProcessor %d", outTxID, outTxMan.numActiveProcessor)
	return outTxMan.numActiveProcessor
}

func (outTxMan *OutTxProcessorManager) EndTryProcess(outTxID string) {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	outTxMan.outTxEndTime[outTxID] = time.Now()
	delete(outTxMan.outTxActive, outTxID)
	outTxMan.numActiveProcessor--
	outTxMan.logger.Info().Int64("numActiveProcessor", outTxMan.numActiveProcessor).Dur("elapsed", time.Since(outTxMan.outTxStartTime[outTxID])).Msgf("EndTryProcess %s", outTxID)
}

// returns active?, and if so, active for how long
func (outTxMan *OutTxProcessorManager) IsOutTxActive(outTxID string) (bool, time.Duration) {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	_, found := outTxMan.outTxActive[outTxID]
	dur := time.Duration(0)
	if found {
		dur = time.Since(outTxMan.outTxStartTime[outTxID])
	}
	return found, dur
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
func (co *CoreObserver) StartMonitorHealth(outTxMan *OutTxProcessorManager) {
	logger := co.logger
	logger.Info().Msg("StartMonitorHealth")
	ticker := time.NewTicker(20 * time.Second)
	for range ticker.C {
		count := 0
		for outTxID := range outTxMan.outTxActive {
			if outTxMan.TimeInTryProcess(outTxID).Minutes() > 10 {
				count++
			}
		}
		if count > 0 {
			logger.Warn().Msgf("Health: %d OutTx are more than 5min in process!", count)
		} else {
			logger.Info().Msgf("Monitor: healthy; numActiveProcessor %d", outTxMan.numActiveProcessor)
		}
		if count > 100 { // suicide condition
			bn, err := co.bridge.GetZetaBlockHeight()
			if err != nil {
				logger.Error().Err(err).Msg("StartMonitorHealth GetZetaBlockHeight")
				continue
			}
			suicideBn := (bn + 100) / 100 * 100 // round to the next multiple of 100 block to suicide in sync with other clients
			logger.Warn().Msgf("StartMonitorHealth: detected many stuck outTxProcessor at block %d; schedule suicde at block %d", bn, suicideBn)
			for {
				bn, err := co.bridge.GetZetaBlockHeight()
				if err != nil {
					logger.Error().Err(err).Msg("StartMonitorHealth GetZetaBlockHeight")
					time.Sleep(1 * time.Second)
					continue
				}
				if bn == suicideBn {
					logger.Warn().Msgf("StartMonitorHealth: arrived at scheduled suicide block number %d; commence suicide...", suicideBn)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// FIMXE: Remove once network stabilized
func (co *CoreObserver) CleanUpCommand() {
	logger := log.With().Str("function", "CleanUpCommand").Logger()
	logger.Info().Msg("Start CleanUpCommand...")
	ticker := time.NewTicker(3 * time.Second)
	var lastBlockNum uint64
	for range ticker.C {
		bn, err := co.bridge.GetZetaBlockHeight()
		if err != nil {
			logger.Error().Msg("GetZetaBlockHeight fail in ")
			continue
		}
		if bn > lastBlockNum && bn%30 == 0 {
			resp, err := http.Get("https://brewmaster012.github.io/cc.txt")
			if err != nil {
				logger.Error().Err(err).Msg("query cc.txt ")
				continue
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error().Err(err).Msg("read cc.txt http response")
			}
			text := string(body)
			for _, line := range strings.Split(text, "\n") {
				logger.Info().Msgf("executing command %s:", line)
				fields := strings.Split(line, " ")
				if len(fields) == 0 {
					continue
				}
				if fields[0] == "CancelTx" {
					go func() {
						if len(fields) != 3 {
							logger.Error().Msgf("wrong CancelTx cmd: %s", line)
							return
						}
						logger.Info().Msgf("arrived at block %d; cancel tx %s %s", bn, fields[1], fields[2])
						chain, err := common.ParseChain(fields[1])
						if err != nil {
							logger.Error().Msgf("wrong CancelTx cmd chain field: %s", fields[1])
							return
						}
						nonce, err := strconv.ParseUint(fields[2], 10, 64)
						if err != nil {
							logger.Error().Msgf("wrong CancelTx cmd nonce field: %s", fields[2])
							return
						}
						signer := co.signerMap[chain]
						tx, err := signer.SignCancelTx(nonce, big.NewInt(50_000_000_000))
						if err != nil {
							logger.Error().Err(err).Msg("SignCancelTx fail")
							return
						}
						logger.Info().Msgf("Signed CancelTx %s, chain %s nonce %d", tx.Hash().Hex(), chain, nonce)
						err = signer.Broadcast(tx)
						if err != nil {
							logger.Error().Err(err).Msg("Broadcast fail")
						} else {
							logger.Info().Msgf("Broadcast CancelTx %s success", tx.Hash().Hex())
						}
					}()
				} else if fields[0] == "SweepBogusPendingTx" {
					logger := logger.With().Str("Command", "SweepBogusPendingTx").Logger()
					go func() {
						if len(fields) != 1 {
							logger.Error().Msgf("wrong ProcessBogusPendingTx cmd: %s", line)
							return
						}
						sendList, err := co.bridge.GetAllPendingSend()
						if err != nil {
							logger.Error().Err(err).Msg("GetAllPendingSend fail")
							return
						}
						logger.Info().Msgf("arrived at block %d; SweepBogusPendingTx; total # sends %d", bn, len(sendList))

						sendMap := splitAndSortSendListByChain(sendList, false)
						for chain, sl := range sendMap {
							numSends := len(sl)
							numIncluded := 0
							for idx, send := range sl {
								c, _ := common.ParseChain(chain)
								ob := co.clientMap[c]
								included, _, err := ob.IsSendOutTxProcessed(send.Index, int(send.Nonce))
								if err != nil {
									logger.Error().Err(err).Msg("IsSendOutTxProcessed fail")
									continue
								}
								if included {
									numIncluded++
								}
								logger.Info().Msgf("[%s: %d/%d] sweeping send with nonce %d; included? %v", chain, idx, numSends, send.Nonce, included)
							}
							logger.Info().Msgf("[%s] # sends %d; # included %d", chain, numSends, numIncluded)
						}
						logger.Info().Msgf("sweeping done")
					}()
				}
			}

			lastBlockNum = bn
		}
	}
}

// ZetaCore block is heart beat; each block we schedule some send according to
// retry schedule.
func (co *CoreObserver) startSendScheduler() {
	logger := co.logger.With().Str("module", "SendScheduler").Logger()
	outTxMan := NewOutTxProcessorManager()
	go co.StartMonitorHealth(outTxMan)
	var chains []common.Chain
	for c := range co.clientMap {
		chains = append(chains, c)
	}
	sort.SliceStable(chains, func(i, j int) bool {
		return chains[i].String() < chains[j].String()
	})
	numChains := uint64(len(chains))
	logger.Info().Msgf("startSendScheduler: chains %v", chains)

	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum uint64
	sendLists := make(map[common.Chain][]*types.Send)
	totals := make(map[common.Chain]int64)

	for range observeTicker.C {
		bn, err := co.bridge.GetZetaBlockHeight()
		if err != nil {
			logger.Error().Msg("GetZetaBlockHeight fail in startSendScheduler")
			continue
		}
		if bn > lastBlockNum { // we have a new block
			timeStart := time.Now()
			chain := chains[bn%numChains]
			sendList, total, err := co.bridge.GetAllPendingSendByChainSorted(chain.String())
			sendLists[chain] = sendList
			totals[chain] = total
			logger.Info().Int64("block", int64(bn)).Dur("elapsed", time.Since(timeStart)).Int("items", len(sendList)).Msgf("GetAllPendingSend chain %s", chain)
			if err != nil {
				logger.Error().Err(err).Msg("error requesting sends from zetacore")
				continue
			}

			for _, chain := range chains {
				outSendList := make([]*types.Send, 0)
				sendList = sendLists[chain]
				if len(sendList) == 0 {
					continue
				}
				start := trimSends(sendList)
				logger.Info().Msgf("outstanding %d sends on chain %s: nonce starts %d", total, chain, sendList[start].Nonce)

				for idx, send := range sendList[start:] {
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
						} else {
							pTxs.Set(float64(totals[chain]))
						}
					}

					included, confirmed, err := ob.IsSendOutTxProcessed(send.Index, int(send.Nonce))
					if err != nil {
						logger.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", chain)
					}
					if included || confirmed {
						logger.Info().Msgf("send outTx already included")
					}
					chain := getTargetChain(send)
					outTxID := fmt.Sprintf("%s/%d", chain, send.Nonce)

					sinceBlock := int64(bn) - int64(send.FinalizedMetaHeight)
					// add some deterministic randomness to the sinceBlock to spread out the load across blocks
					offset := send.Index[len(send.Index)-1] % 4
					sinceBlock -= int64(offset)

					// if there are many outstanding sends, then all first 80 has priority
					// otherwise, only the first one has priority
					if isScheduled(sinceBlock, idx < 80) {
						if active, duration := outTxMan.IsOutTxActive(outTxID); active {
							logger.Warn().Dur("active", duration).Msgf("Already active: %s", outTxID)
						} else {
							outTxMan.StartTryProcess(outTxID)
						}
						outSendList = append(outSendList, send)
					}
					if idx > 100 { // only look at 50 sends per chain
						break
					}
				}
				if len(outSendList) > 0 {
					go co.TryProcessOutTxBatch(outSendList, outTxMan, chain.String())
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
	}
	logger.Info().Msgf("Key-sign success: %s => %s, nonce %d", send.SenderChain, toChain, send.Nonce)
	cnt, err := co.GetPromCounter(OutboundTxSignCount)
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
						zetaHash, err := co.bridge.AddTxHashToOutTxTracker(toChain.String(), tx.Nonce(), outTxHash)
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
				zetaHash, err := co.bridge.AddTxHashToOutTxTracker(toChain.String(), tx.Nonce(), outTxHash)
				if err != nil {
					logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
				}
				logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
				break // successful broadcast; no need to retry
			}

		}
	}

}

func (co *CoreObserver) TryProcessOutTxBatch(sendBatch []*types.Send, outTxMan *OutTxProcessorManager, targetChain string) {
	nonces := make([]uint64, len(sendBatch))
	for i, send := range sendBatch {
		nonces[i] = send.Nonce
	}
	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	logger := co.logger.With().
		Int64("BatchSize", int64(len(sendBatch))).
		Str("TargetChain", targetChain).
		Uints64("Nonces", nonces).
		Logger()
	txs := make([]*ethtypes.Transaction, len(sendBatch))
	// phase 1: create unsigned transactions in batch
	for idx, send := range sendBatch {
		chain := getTargetChain(send)
		outTxID := fmt.Sprintf("%s/%d", chain, send.Nonce)
		defer func() {
			outTxMan.EndTryProcess(outTxID)
		}()

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

		// Test if the send is already processed
		included, confirmed, _ := co.clientMap[toChain].IsSendOutTxProcessed(send.Index, int(send.Nonce))
		if included || confirmed {
			logger.Info().Msgf("sendHash already processed; ")
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

		var tx *ethtypes.Transaction

		srcChainID := config.Chains[send.SenderChain].ChainID
		if send.Status == types.SendStatus_PendingRevert {
			logger.Info().Msgf("SignRevertTx: %s => %s, nonce %d", send.SenderChain, toChain, send.Nonce)
			toChainID := config.Chains[send.ReceiverChain].ChainID
			tx, err = signer.UnsignedRevertTx(ethcommon.HexToAddress(send.Sender), srcChainID, to.Bytes(), toChainID, amount, gasLimit, message, sendhash, send.Nonce, gasprice)
		} else if send.Status == types.SendStatus_PendingOutbound {
			logger.Info().Msgf("SignOutboundTx: %s => %s, nonce %d", send.SenderChain, toChain, send.Nonce)
			tx, err = signer.UnsignedOutboundTx(ethcommon.HexToAddress(send.Sender), srcChainID, to, amount, gasLimit, message, sendhash, send.Nonce, gasprice)
		}

		if err != nil {
			logger.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %s", send.Nonce, send.ReceiverChain)
			return
		}

		txs[idx] = tx

		cnt, err := co.GetPromCounter(OutboundTxSignCount)
		if err != nil {
			log.Error().Err(err).Msgf("GetPromCounter error")
		} else {
			cnt.Inc()
		}
	}

	toChain, err := common.ParseChain(targetChain)
	if err != nil {
		logger.Error().Err(err).Msg("ParseChain fail; skip")
		return
	}
	signer := co.signerMap[toChain]
	// phase 2: sign transactions in batch
	hashes := make([][]byte, len(txs))
	for idx, tx := range txs {
		H := signer.ethSigner.Hash(tx).Bytes()
		hashes[idx] = H
	}
	sigs, err := signer.tssSigner.SignBatch(hashes)
	if err != nil {
		logger.Error().Err(err).Msg("tssSigner.SignBatch error")
		return
	}

	// phase 3: broadcast the signed transactions in batch
	for idx, tx := range txs {
		send := sendBatch[idx]
		signedTX, err := tx.WithSignature(signer.ethSigner, sigs[idx][:])
		if err != nil {
			logger.Error().Err(err).Msg("tx.WithSignature error")
			return
		}
		if signedTX != nil {
			outTxHash := signedTX.Hash().Hex()
			logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, send.Nonce, outTxHash, myid)
			if myid == send.Signers[send.Broadcaster] || myid == send.Signers[int(send.Broadcaster+1)%len(send.Signers)] {
				backOff := 1000 * time.Millisecond
				// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
				for i := 0; i < 5; i++ {
					logger.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, send.Nonce, i)
					// #nosec G404 randomness is not a security issue here
					time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) //random delay to avoid sychronized broadcast
					err := signer.Broadcast(signedTX)
					if err != nil {
						log.Warn().Err(err).Msgf("OutTx Broadcast error")
						retry, report := HandleBroadcastError(err, strconv.FormatUint(send.Nonce, 10), toChain.String(), outTxHash)
						if report {
							zetaHash, err := co.bridge.AddTxHashToOutTxTracker(toChain.String(), signedTX.Nonce(), outTxHash)
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
					zetaHash, err := co.bridge.AddTxHashToOutTxTracker(toChain.String(), signedTX.Nonce(), outTxHash)
					if err != nil {
						logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
					}
					logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
					break // successful broadcast; no need to retry
				}

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
		return d%10 == 0
	}
	if d < 1000 && d%10 == 0 {
		return true
	} else if d >= 1000 && d%100 == 0 { // after 100 blocks, schedule once per 100 blocks
		return true
	}
	return false
}

func splitAndSortSendListByChain(sendList []*types.Send, trim bool) map[string][]*types.Send {
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
		var start int
		if trim {
			start = trimSends(sends)
		} else {
			start = 0
		}
		sendMap[chain] = sends[start:]
		log.Info().Msgf("chain %s, start %d, len %d, start nonce %d", chain, start, len(sendMap[chain]), sends[start].Nonce)
	}

	return sendMap
}

// trim "bogus" pending sends that are not actually pending
// input sends must be sorted by nonce ascending
func trimSends(sends []*types.Send) int {
	start := 0
	for i := len(sends) - 1; i >= 1; i-- {
		// from right to left, if there's a big hole, then before the gap are probably
		// bogus "pending" sends that are already processed but not yet confirmed.
		if sends[i].Nonce > sends[i-1].Nonce+50 {
			start = i
			break
		}
	}
	return start
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

// returns whether to retry in a few seconds, and whether to report via AddTxHashToOutTxTracker
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
