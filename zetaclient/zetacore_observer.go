package zetaclient

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/rs/zerolog"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"gitlab.com/thorchain/tss/go-tss/keygen"

	prom "github.com/prometheus/client_golang/prometheus"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"

	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
)

const (
	OUTBOUND_TX_SIGN_COUNT = "zetaclient_outbound_tx_sign_count"
)

type CoreObserver struct {
	sendQueue []*types.Send
	bridge    *ZetaCoreBridge
	signerMap map[common.Chain]*Signer
	clientMap map[common.Chain]*ChainObserver
	metrics   *metrics.Metrics
	tss       *TSS

	// channels for shepherd manager
	sendNew     chan *types.Send
	sendDone    chan *types.Send
	signerSlots chan bool
	shepherds   map[string]bool

	fileLogger *zerolog.Logger
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]*Signer, clientMap map[common.Chain]*ChainObserver, metrics *metrics.Metrics, tss *TSS) *CoreObserver {
	co := CoreObserver{}
	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap
	co.sendQueue = make([]*types.Send, 0)

	co.clientMap = clientMap
	co.metrics = metrics

	err := metrics.RegisterCounter(OUTBOUND_TX_SIGN_COUNT, "number of outbound tx signed")
	if err != nil {
		log.Error().Err(err).Msg("error registering counter")
	}

	co.sendNew = make(chan *types.Send)
	co.sendDone = make(chan *types.Send)
	MAX_SIGNERS := 100 // assuming each signer takes 100s to finish (have outTx included), then throughput is bounded by 100/100 = 1 tx/s
	co.signerSlots = make(chan bool, MAX_SIGNERS)
	for i := 0; i < MAX_SIGNERS; i++ {
		co.signerSlots <- true
	}
	co.shepherds = make(map[string]bool)

	logFile, err := os.OpenFile("zetacore_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		// Can we log an error before we have our logger? :)
		log.Error().Err(err).Msgf("there was an error creating a logFile on zetacore")
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	co.fileLogger = &fileLogger

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
	go co.startObserve()
	go co.shepherdManager()

	noKeygen := os.Getenv("NO_KEYGEN")
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
					log.Fatal().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
					continue
				}
				// Keygen succeed! Report TSS address
				log.Info().Msgf("Keygen success! keygen response: %v...", res)
				err = co.tss.SetPubKey(res.PubKey)
				if err != nil {
					log.Error().Msgf("SetPubKey fail")
					continue
				}

				for _, chain := range config.ChainsEnabled {
					_, err = co.bridge.SetTSS(chain, co.tss.Address().Hex(), co.tss.PubkeyInBech32)
					if err != nil {
						log.Error().Err(err).Msgf("SetTSS fail %s", chain)
					}
				}

				// Keysign test: sanity test
				log.Info().Msgf("test keysign...")
				TestKeysign(co.tss.PubkeyInBech32, co.tss.Server)
				log.Info().Msg("test keysign finished. exit keygen loop. ")

				for _, chain := range config.ChainsEnabled {
					err = co.clientMap[chain].PostNonceIfNotRecorded()
					if err != nil {
						log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain)
					}
				}

				return
			}
		}()
		return
	}
}

// startObserve retrieves the pending list of Sends from ZetaCore every 10s
// for each new send, it tries to launch a send shepherd.
// the send shepherd makes sure the send is settled on all chains.
func (co *CoreObserver) startObserve() {
	observeTicker := time.NewTicker(12 * time.Second)
	for range observeTicker.C {
		sendList, err := co.bridge.GetAllPendingSend()
		if err != nil {
			log.Error().Err(err).Msg("error requesting sends from zetacore")
			continue
		}
		if len(sendList) > 0 {
			log.Info().Msgf("#pending send: %d", len(sendList))
		}
		sort.Slice(sendList, func(i, j int) bool {
			return sendList[i].Nonce < sendList[j].Nonce
		})
		for _, send := range sendList {
			if send.Status == types.SendStatus_PendingOutbound || send.Status == types.SendStatus_PendingRevert {
				co.sendNew <- send
			} //else if send.Status == types.SendStatus_Mined || send.Status == types.SendStatus_Reverted || send.Status == types.SendStatus_Aborted {
		}
	}
}

func (co *CoreObserver) shepherdManager() {
	numShepherds := 0
	for {
		select {
		case send := <-co.sendNew:
			if _, ok := co.shepherds[send.Index]; !ok {
				log.Info().Msgf("shepherd manager: new send %s", send.Index)
				co.shepherds[send.Index] = true
				log.Info().Msg("waiting on a signer slot...")
				<-co.signerSlots
				log.Info().Msg("got a signer slot! spawn shepherd")
				go co.shepherdSend(send)
				numShepherds++
				log.Info().Msgf("new shepherd: %d shepherds in total", numShepherds)
			}
		case send := <-co.sendDone:
			delete(co.shepherds, send.Index)
			numShepherds--
			log.Info().Msgf("remove shepherd: %d shepherds left", numShepherds)
		}
	}
}

// Once this function receives a Send, it will make sure that the send is processed and confirmed
// on external chains and ZetaCore.
// FIXME: make sure that ZetaCore is updated when the Send cannot be processed.
func (co *CoreObserver) shepherdSend(send *types.Send) {
	startTime := time.Now()
	confirmDone := make(chan bool, 1)
	coreSendDone := make(chan bool, 1)
	numQueries := 0
	keysignCount := 0

	defer func() {
		elapsedTime := time.Since(startTime)
		if keysignCount > 0 {
			log.Info().Msgf("shepherd stopped: numQueries %d; elapsed time %s; keysignCount %d", numQueries, elapsedTime, keysignCount)
			co.fileLogger.Info().Msgf("shepherd stopped: numQueries %d; elapsed time %s; keysignCount %d", numQueries, elapsedTime, keysignCount)
		}
		co.signerSlots <- true
		co.sendDone <- send
		confirmDone <- true
		coreSendDone <- true
	}()

	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	amount, ok := new(big.Int).SetString(send.ZetaMint, 0)
	if !ok {
		log.Error().Msg("error converting MBurnt to big.Int")
		return
	}

	var to ethcommon.Address
	var err error
	var toChain common.Chain
	if send.Status == types.SendStatus_PendingRevert {
		to = ethcommon.HexToAddress(send.Sender)
		toChain, err = common.ParseChain(send.SenderChain)
		log.Info().Msgf("Abort: reverting inbound")
	} else if send.Status == types.SendStatus_PendingOutbound {
		to = ethcommon.HexToAddress(send.Receiver)
		toChain, err = common.ParseChain(send.ReceiverChain)
	}
	if err != nil {
		log.Error().Err(err).Msg("ParseChain fail; skip")
		return
	}

	// Early return if the send is already processed
	included, confirmed, _ := co.clientMap[toChain].IsSendOutTxProcessed(send.Index, int(send.Nonce))
	if included || confirmed {
		log.Info().Msgf("sendHash %s already processed; exit signer", send.Index)
		return
	}

	signer := co.signerMap[toChain]
	message, err := base64.StdEncoding.DecodeString(send.Message)
	if err != nil {
		log.Err(err).Msgf("decode send.Message %s error", send.Message)
	}

	gasLimit := send.GasLimit
	if gasLimit < 50_000 {
		gasLimit = 50_000
	}

	log.Info().Msgf("chain %s minting %d to %s, nonce %d, finalized %d", toChain, amount, to.Hex(), send.Nonce, send.FinalizedMetaHeight)
	sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
	if err != nil || len(sendHash) != 32 {
		log.Err(err).Msgf("decode sendHash %s error", send.Index)
		return
	}
	var sendhash [32]byte
	copy(sendhash[:32], sendHash[:32])
	gasprice, ok := new(big.Int).SetString(send.GasPrice, 10)
	if !ok {
		log.Err(err).Msgf("cannot convert gas price  %s ", send.GasPrice)
		return
	}
	var tx *ethtypes.Transaction

	signloopDone := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-confirmDone:
				return
			default:
				included, confirmed, err := co.clientMap[toChain].IsSendOutTxProcessed(send.Index, int(send.Nonce))
				if err != nil {
					numQueries++
				}
				if included || confirmed {
					log.Info().Msgf("sendHash %s included; kill this shepherd", send.Index)
					signloopDone <- true
					return
				}
				time.Sleep(12 * time.Second)
			}
		}
	}()

	// watch ZetaCore /zeta-chain/send/<sendHash> endpoint; send coreSendDone when the state of the send is updated;
	// e.g. pendingOutbound->outboundMined; or pendingOutbound->pendingRevert
	go func() {
		for {
			select {
			case <-coreSendDone:
				return
			default:
				newSend, err := co.bridge.GetSendByHash(send.Index)
				if err != nil || send == nil {
					log.Info().Msgf("sendHash %s cannot be found in ZetaCore; kill the shepherd", send.Index)
					signloopDone <- true
				}
				if newSend.Status != send.Status {
					log.Info().Msgf("sendHash %s status changed to %s from %s; kill the shepherd", send.Index, newSend.Status, send.Status)
					signloopDone <- true
				}
				time.Sleep(12 * time.Second)
			}
		}
	}()

	// The following keysign loop tries to sign outbound tx until the following conditions are met:
	// 1. zetacore /zeta-chain/send/<sendHash> endpoint returns a changed status
	// 2. outTx is confirmed to be successfully or failed
	signTicker := time.NewTicker(time.Second)
	signInterval := 128 * time.Second // minimum gap between two keysigns
	lastSignTime := time.Unix(1, 0)
SIGNLOOP:
	for range signTicker.C {
		select {
		case <-signloopDone:
			log.Info().Msg("breaking SignOutBoundTx loop: outbound already processed")
			break SIGNLOOP
		default:
			if co.clientMap[toChain].MinNonce == int(send.Nonce) {
				log.Warn().Msgf("this signer is likely blocking subsequent txs! nonce %d", send.Nonce)
				signInterval = 32 * time.Second
			}
			tnow := time.Now()
			if tnow.Before(lastSignTime.Add(signInterval)) {
				continue
			}
			if tnow.Unix()%16 == int64(sendhash[0])%16 { // weakly sync the TSS signers
				included, confirmed, _ := co.clientMap[toChain].IsSendOutTxProcessed(send.Index, int(send.Nonce))
				if included || confirmed {
					log.Info().Msgf("sendHash %s already confirmed; skip it", send.Index)
					break SIGNLOOP
				}
				srcChainID := config.Chains[send.SenderChain].ChainID
				if send.Status == types.SendStatus_PendingRevert {
					log.Info().Msgf("SignRevertTx: %s => %s, nonce %d, sendHash %s", send.SenderChain, toChain, send.Nonce, send.Index)
					toChainID := config.Chains[send.ReceiverChain].ChainID
					tx, err = signer.SignRevertTx(ethcommon.HexToAddress(send.Sender), srcChainID, to.Bytes(), toChainID, amount, gasLimit, message, sendhash, send.Nonce, gasprice)
				} else if send.Status == types.SendStatus_PendingOutbound {
					log.Info().Msgf("SignOutboundTx: %s => %s, nonce %d, sendHash %s", send.SenderChain, toChain, send.Nonce, send.Index)
					tx, err = signer.SignOutboundTx(ethcommon.HexToAddress(send.Sender), srcChainID, to, amount, gasLimit, message, sendhash, send.Nonce, gasprice)
				}
				if err != nil {
					log.Warn().Err(err).Msgf("SignOutboundTx error: nonce %d chain %s", send.Nonce, send.ReceiverChain)
					continue
				}
				lastSignTime = time.Now()
				cnt, err := co.GetPromCounter(OUTBOUND_TX_SIGN_COUNT)
				if err != nil {
					log.Error().Err(err).Msgf("GetPromCounter error")
				} else {
					cnt.Inc()
				}

				// if tx is nil, maybe I'm not an active signer?
				if tx != nil {
					outTxHash := tx.Hash().Hex()
					log.Info().Msgf("on chain %s nonce %d, sendHash: %s, outTxHash %s signer %s", signer.chain, send.Nonce, send.Index[:6], outTxHash, myid)
					if myid == send.Signers[send.Broadcaster] || myid == send.Signers[int(send.Broadcaster+1)%len(send.Signers)] {
						backOff := 1000 * time.Millisecond
						for i := 0; i < 5; i++ { // retry loop: 1s, 2s, 4s, 8s, 16s
							log.Info().Msgf("broadcasting tx %s to chain %s: nonce %d, retry %d", outTxHash, toChain, send.Nonce, i)
							// #nosec G404 randomness is not a security issue here
							time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) //random delay to avoid sychronized broadcast
							err = signer.Broadcast(tx)
							// TODO: the following error handling is robust?
							if err == nil {
								log.Err(err).Msgf("Broadcast success: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
								co.fileLogger.Err(err).Msgf("Broadcast success: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
								break // break the retry loop
							} else if strings.Contains(err.Error(), "nonce too low") {
								log.Info().Msgf("nonce too low! this might be a unnecessary keysign. increase re-try interval and awaits outTx confirmation")
								co.fileLogger.Err(err).Msgf("Broadcast nonce too low: nonce %d chain %s outTxHash %s; increase re-try interval", send.Nonce, toChain, outTxHash)
								break
							} else if strings.Contains(err.Error(), "replacement transaction underpriced") {
								log.Err(err).Msgf("Broadcast replacement: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
								co.fileLogger.Err(err).Msgf("Broadcast replacement: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
								break
							} else if strings.Contains(err.Error(), "already known") { // this is error code from QuickNode
								log.Err(err).Msgf("Broadcast duplicates: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
								co.fileLogger.Err(err).Msgf("Broadcast duplicates: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
								break
							} else { // most likely an RPC error, such as timeout or being rate limited. Exp backoff retry
								log.Err(err).Msgf("Broadcast error: nonce %d chain %s outTxHash %s; retring...", send.Nonce, toChain, outTxHash)
								co.fileLogger.Err(err).Msgf("Broadcast error: nonce %d chain %s outTxHash %s; retrying...", send.Nonce, toChain, outTxHash)
								time.Sleep(backOff)
							}
							backOff *= 2
						}

					}
					// if outbound tx fails, kill this shepherd, a new one will be later spawned.
					co.clientMap[toChain].AddTxHashToWatchList(outTxHash, int(send.Nonce), send.Index)
					co.fileLogger.Info().Msgf("Keysign: %s => %s, nonce %d, outTxHash %s; keysignCount %d", send.SenderChain, toChain, send.Nonce, outTxHash, keysignCount)
					keysignCount++
					signInterval *= 2 // exponential backoff
				}
			}
		}
	}
}
