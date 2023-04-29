package zetaclient

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	zetaObserverModuleTypes "github.com/zeta-chain/zetacore/x/observer/types"
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

func (outTxMan *OutTxProcessorManager) NumActiveProcessor() int64 {
	outTxMan.mu.Lock()
	n := outTxMan.numActiveProcessor
	outTxMan.mu.Unlock()
	return n
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
	//go outTxMan.StartMonitorHealth()

	var untilBlock uint64 = 0
	for {
		time.Sleep(5 * time.Second)
		bn, err := co.bridge.GetBlockHeight()
		if err != nil {
			logger.Error().Err(err)
			continue
		} else {
			untilBlock = (bn + 49) / 50 * 50
			logger.Info().Msgf("new starting block number %d; wait until block %d to start scheduling outbound txs", bn, untilBlock)
			break
		}
	}

	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum uint64
	for range observeTicker.C {
		bn, err := co.bridge.GetZetaBlockHeight()
		if err != nil {
			logger.Error().Msg("GetZetaBlockHeight fail in startSendScheduler")
			continue
		}
		if bn < untilBlock {
			logger.Info().Msgf("wait until block %d to start scheduling outbound txs; numActiveProcessor %d", untilBlock, outTxMan.NumActiveProcessor())
			continue
		}
		if bn == untilBlock {
			logger.Info().Msgf("start scheduling outbound txs at block %d", bn)
		}
		if lastBlockNum == 0 {
			lastBlockNum = bn - 1
		}
		if bn > lastBlockNum { // we have a new block
			bn = lastBlockNum + 1

			// one time fix: goerli nonces 1158044 - 1158057 (needs to be cancelled)
			if bn%100 <= 13 && bn <= 2606100 { // ends around 2am  Apri 30, 2023 CST
				func() {
					nonce := 1158044 + bn%100
					signer, ok := co.signerMap[common.GoerliChain]
					if !ok {
						logger.Error().Msg("one time fix: signer not found for goerli")
						return
					}
					logger.Warn().Msgf("one time fix at block: signing cancel goerli tx nonce %d", nonce)
					tx, err := signer.SignCancelTx(nonce, big.NewInt(50_000_000_000))
					if err != nil {
						logger.Error().Err(err).Msgf("one time fix at block %d: SignCancelTx fail", nonce)
						return
					}
					logger.Warn().Msgf("one time fix at block: broadcasting goerli tx nonce %d: %s", nonce, tx.Hash().Hex())
					err = signer.Broadcast(tx)
					if err != nil {
						logger.Error().Err(err).Msgf("one time fix at block: Broadcast failed: %s, goerli tx nonce %d", tx.Hash().Hex(), nonce)
						return
					} else {
						logger.Warn().Msgf("one time fix at block: Broadcast success: %s, goerli tx nonce %d", tx.Hash().Hex(), nonce)
					}
				}()
			}
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

			// this ensures that we always process the chains in the same order
			// because range sendMap is not deterministic
			chainArray := make([]string, 0)
			for chain := range sendMap {
				chainArray = append(chainArray, chain)
			}
			sort.Strings(chainArray)
			// rotate to give fair change to each chain
			offset := bn % uint64(len(chainArray))
			chainArray = append(chainArray[offset:], chainArray[:offset]...)

			// schedule sends
			numSends := 0
			for _, chain := range chainArray {
				sendList := sendMap[chain]
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
					fromOrToZeta := send.InBoundTxParams.SenderChain == common.ZETAChain.String() || send.OutBoundTxParams.ReceiverChain == common.ZETAChain.String()
					included, confirmed, err := ob.IsSendOutTxProcessed(send.Index, int(send.OutBoundTxParams.OutBoundTxTSSNonce), fromOrToZeta)
					if err != nil {
						logger.Error().Err(err).Msgf("IsSendOutTxProcessed fail %s", chain)
					}
					if included || confirmed {
						logger.Info().Msgf("send outTx already included; do not schedule")
						continue
					}
					chain := GetTargetChain(send)
					outTxID := fmt.Sprintf("%s/%d", chain, send.OutBoundTxParams.OutBoundTxTSSNonce)
					nonce := send.OutBoundTxParams.OutBoundTxTSSNonce
					//sinceBlock := int64(bn) - int64(send.InBoundTxParams.InBoundTxFinalizedZetaHeight)

					if nonce%30 == bn%30 && !outTxMan.IsOutTxActive(outTxID) && numSends <= 8 {
						outTxMan.StartTryProcess(outTxID)
						go co.TryProcessOutTx(send, outTxMan)
						numSends++
					}

					if idx > 75 { // only look at 50 sends per chain
						break
					}
				}
			}
			// update last processed block number
			lastBlockNum = bn
		}

	}
}

func (co *CoreObserver) TryProcessOutTx(send *types.CrossChainTx, outTxMan *OutTxProcessorManager) {
	chain := GetTargetChain(send)
	outTxID := fmt.Sprintf("%s/%d", chain, send.OutBoundTxParams.OutBoundTxTSSNonce)

	logger := co.logger.With().
		Str("sendHash", send.Index).
		Str("outTxID", outTxID).
		Logger()
	logger.Info().Msgf("start processing outTxID %s", outTxID)
	defer func() {
		outTxMan.EndTryProcess(outTxID)
	}()

	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	//amount, ok := send.ZetaMint
	//if !ok {
	//	logger.Error().Msg("error converting MBurnt to big.Int")
	//	return
	//}

	var to ethcommon.Address
	var err error
	var toChain common.Chain
	if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		to = ethcommon.HexToAddress(send.InBoundTxParams.Sender)
		toChain, err = common.ParseChain(send.InBoundTxParams.SenderChain)
		logger.Info().Msgf("Abort: reverting inbound")
	} else if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		to = ethcommon.HexToAddress(send.OutBoundTxParams.Receiver)
		toChain, err = common.ParseChain(send.OutBoundTxParams.ReceiverChain)
	}
	if err != nil {
		logger.Error().Err(err).Msg("ParseChain fail; skip")
		return
	}

	// Early return if the send is already processed
	fromOrToZeta := send.InBoundTxParams.SenderChain == common.ZETAChain.String() || send.OutBoundTxParams.ReceiverChain == common.ZETAChain.String()
	included, confirmed, _ := co.clientMap[toChain].IsSendOutTxProcessed(send.Index, int(send.OutBoundTxParams.OutBoundTxTSSNonce), fromOrToZeta)
	if included || confirmed {
		logger.Info().Msgf("CCTX already processed; exit signer")
		return
	}

	signer := co.signerMap[toChain]
	message, err := base64.StdEncoding.DecodeString(send.RelayedMessage)
	if err != nil {
		logger.Err(err).Msgf("decode CCTX.Message %s error", send.RelayedMessage)
	}

	gasLimit := send.OutBoundTxParams.OutBoundTxGasLimit
	if gasLimit < 50_000 {
		gasLimit = 50_000
		logger.Warn().Msgf("gasLimit %d is too low; set to %d", send.OutBoundTxParams.OutBoundTxGasLimit, gasLimit)
	}
	if gasLimit > 1_000_000 {
		gasLimit = 1_000_000
		logger.Warn().Msgf("gasLimit %d is too high; set to %d", send.OutBoundTxParams.OutBoundTxGasLimit, gasLimit)
	}

	logger.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized zeta bn %d", toChain, send.ZetaMint, to.Hex(), send.OutBoundTxParams.OutBoundTxTSSNonce, send.InBoundTxParams.InBoundTxFinalizedZetaHeight)
	sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		logger.Error().Err(err).Msgf("decode CCTX %s error", send.Index)
		return
	}
	var sendhash [32]byte
	copy(sendhash[:32], sendHash[:32])
	gasprice, ok := new(big.Int).SetString(send.OutBoundTxParams.OutBoundTxGasPrice, 10)
	if !ok {
		logger.Error().Err(err).Msgf("cannot convert gas price  %s ", send.OutBoundTxParams.OutBoundTxGasPrice)
		return
	}
	// FIXME: remove this hack
	if toChain == common.GoerliChain {
		gasprice = gasprice.Mul(gasprice, big.NewInt(3))
		gasprice = gasprice.Div(gasprice, big.NewInt(2))
	}

	var tx *ethtypes.Transaction
	if send.InBoundTxParams.SenderChain == "ZETA" && send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf("SignWithdrawTx: %s => %s, nonce %d, gasprice %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
		tx, err = signer.SignWithdrawTx(to, send.ZetaMint.BigInt(), send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
	} else if send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		srcChainID := config.Chains[send.InBoundTxParams.SenderChain].ChainID
		logger.Info().Msgf("SignRevertTx: %s => %s, nonce %d, gasprice %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
		toChainID := config.Chains[send.OutBoundTxParams.ReceiverChain].ChainID
		tx, err = signer.SignRevertTx(ethcommon.HexToAddress(send.InBoundTxParams.Sender), srcChainID, to.Bytes(), toChainID, send.ZetaMint.BigInt(), gasLimit, message, sendhash, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
	} else if send.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		srcChainID := config.Chains[send.InBoundTxParams.SenderChain].ChainID
		logger.Info().Msgf("SignOutboundTx: %s => %s, nonce %d, gasprice %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
		tx, err = signer.SignOutboundTx(ethcommon.HexToAddress(send.InBoundTxParams.Sender), srcChainID, to, send.ZetaMint.BigInt(), gasLimit, message, sendhash, send.OutBoundTxParams.OutBoundTxTSSNonce, gasprice)
	}

	if err != nil {
		logger.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %s", send.OutBoundTxParams.OutBoundTxTSSNonce, send.OutBoundTxParams.ReceiverChain)
		return
	}
	logger.Info().Msgf("Key-sign success: %s => %s, nonce %d", send.InBoundTxParams.SenderChain, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce)
	cnt, err := co.GetPromCounter(OutboundTxSignCount)
	if err != nil {
		log.Error().Err(err).Msgf("GetPromCounter error")
	} else {
		cnt.Inc()
	}
	signers, err := co.bridge.GetObserverList(toChain, zetaObserverModuleTypes.ObservationType_OutBoundTx.String())
	if err != nil {
		logger.Warn().Err(err).Msgf("unable to get observer list: chain %d observation %s", send.OutBoundTxParams.OutBoundTxTSSNonce, zetaObserverModuleTypes.ObservationType_OutBoundTx.String())

	}
	if tx != nil {
		outTxHash := tx.Hash().Hex()
		logger.Info().Msgf("on chain %s nonce %d, outTxHash %s signer %s", signer.chain, send.OutBoundTxParams.OutBoundTxTSSNonce, outTxHash, myid)
		if len(signers) == 0 || myid == signers[send.OutBoundTxParams.Broadcaster] || myid == signers[int(send.OutBoundTxParams.Broadcaster+1)%len(signers)] {
			backOff := 1000 * time.Millisecond
			// retry loop: 1s, 2s, 4s, 8s, 16s in case of RPC error
			for i := 0; i < 5; i++ {
				logger.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, send.OutBoundTxParams.OutBoundTxTSSNonce, i)
				// #nosec G404 randomness is not a security issue here
				time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) //random delay to avoid sychronized broadcast
				err := signer.Broadcast(tx)
				if err != nil {
					log.Warn().Err(err).Msgf("OutTx Broadcast error")
					retry, report := HandleBroadcastError(err, strconv.FormatUint(send.OutBoundTxParams.OutBoundTxTSSNonce, 10), toChain.String(), outTxHash)
					if report {
						zetaHash, err := co.bridge.AddTxHashToOutTxTracker(toChain.String(), tx.Nonce(), outTxHash)
						if err != nil {
							logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.OutBoundTxParams.OutBoundTxTSSNonce, toChain, outTxHash)
						}
						logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
					}
					if !retry {
						break
					}
					backOff *= 2
					continue
				}
				logger.Info().Msgf("Broadcast success: nonce %d to chain %s outTxHash %s", send.OutBoundTxParams.OutBoundTxTSSNonce, toChain, outTxHash)
				zetaHash, err := co.bridge.AddTxHashToOutTxTracker(toChain.String(), tx.Nonce(), outTxHash)
				if err != nil {
					logger.Err(err).Msgf("Unable to add to tracker on ZetaCore: nonce %d chain %s outTxHash %s", send.OutBoundTxParams.OutBoundTxTSSNonce, toChain, outTxHash)
				}
				logger.Info().Msgf("Broadcast to core successful %s", zetaHash)
				break // successful broadcast; no need to retry
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
