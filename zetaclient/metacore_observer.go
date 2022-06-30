package zetaclient

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
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
	bridge    *MetachainBridge
	signerMap map[common.Chain]*Signer
	clientMap map[common.Chain]*ChainObserver
	metrics   *metrics.Metrics
	tss       *TSS

	// channels for shepherd manager
	sendNew     chan *types.Send
	sendDone    chan *types.Send
	signerSlots chan bool
	shepherds   map[string]bool
}

func NewCoreObserver(bridge *MetachainBridge, signerMap map[common.Chain]*Signer, clientMap map[common.Chain]*ChainObserver, metrics *metrics.Metrics, tss *TSS) *CoreObserver {
	co := CoreObserver{}
	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap
	co.sendQueue = make([]*types.Send, 0)

	co.clientMap = clientMap
	co.metrics = metrics

	metrics.RegisterCounter(OUTBOUND_TX_SIGN_COUNT, "number of outbound tx signed")

	co.sendNew = make(chan *types.Send)
	co.sendDone = make(chan *types.Send)
	MAX_SIGNERS := 12 // assuming each signer takes 100s to finish, then throughput is bounded by 100tx/100s = 1tx/s
	co.signerSlots = make(chan bool, MAX_SIGNERS)
	for i := 0; i < MAX_SIGNERS; i++ {
		co.signerSlots <- true
	}
	co.shepherds = make(map[string]bool)

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
	go co.keygenObserve()
}

func (co *CoreObserver) keygenObserve() {
	observeTicker := time.NewTicker(2 * time.Second)
	for range observeTicker.C {
		kg, err := co.bridge.GetKeyGen()
		if err != nil {
			continue
		}
		bn, _ := co.bridge.GetMetaBlockHeight()
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

				_, err = co.bridge.SetTSS(common.GoerliChain, co.tss.Address().Hex(), co.tss.PubkeyInBech32)
				if err != nil {
					log.Error().Err(err).Msgf("SetTSS fail %s", common.GoerliChain)
				}
				_, err = co.bridge.SetTSS(common.BSCTestnetChain, co.tss.Address().Hex(), co.tss.PubkeyInBech32)
				if err != nil {
					log.Error().Err(err).Msgf("SetTSS fail %s", common.BSCTestnetChain)
				}
				_, err = co.bridge.SetTSS(common.MumbaiChain, co.tss.Address().Hex(), co.tss.PubkeyInBech32)
				if err != nil {
					log.Error().Err(err).Msgf("SetTSS fail %s", common.MumbaiChain)
				}
				_, err = co.bridge.SetTSS(common.RopstenChain, co.tss.Address().Hex(), co.tss.PubkeyInBech32)
				if err != nil {
					log.Error().Err(err).Msgf("SetTSS fail %s", common.RopstenChain)
				}

				// Keysign test: sanity test
				log.Info().Msgf("test keysign...")
				TestKeysign(co.tss.PubkeyInBech32, co.tss.Server)
				log.Info().Msg("test keysign finished. exit keygen loop. ")

				err = co.clientMap[common.GoerliChain].PostNonceIfNotRecorded()
				if err != nil {
					log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.GoerliChain)
				}
				err = co.clientMap[common.BSCTestnetChain].PostNonceIfNotRecorded()
				if err != nil {
					log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.BSCTestnetChain)
				}
				err = co.clientMap[common.MumbaiChain].PostNonceIfNotRecorded()
				if err != nil {
					log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.MumbaiChain)
				}
				err = co.clientMap[common.RopstenChain].PostNonceIfNotRecorded()
				if err != nil {
					log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.RopstenChain)
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

		for _, send := range sendList {
			log.Info().Msgf("#pending send: %d", len(sendList))
			if send.Status == types.SendStatus_PendingOutbound || send.Status == types.SendStatus_PendingRevert {
				co.sendNew <- send
			} //else if send.Status == types.SendStatus_Mined || send.Status == types.SendStatus_Reverted || send.Status == types.SendStatus_Aborted {
		}
	}
}

func (co *CoreObserver) shepherdManager() {
	for {
		select {
		case send := <-co.sendNew:
			if _, ok := co.shepherds[send.Index]; !ok {
				log.Info().Msgf("shepherd manager: new send %s", send.Index)
				co.shepherds[send.Index] = true
				log.Info().Msg("waiting on a signer slot...")
				<-co.signerSlots
				log.Info().Msg("got back a signer slot! spawn shepherd")
				go co.shepherdSend(send)
			}
		case send := <-co.sendDone:
			delete(co.shepherds, send.Index)
		}
	}
}

// Once this function receives a Send, it will make sure that the send is processed and confirmed
// on external chains and ZetaCore.
// FIXME: make sure that ZetaCore is updated when the Send cannot be processed.
func (co *CoreObserver) shepherdSend(send *types.Send) {
	done2 := make(chan bool, 1)
	defer func() {
		log.Info().Msg("Giving back a signer slot")
		co.signerSlots <- true
		co.sendDone <- send
		done2 <- true
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
	_, confirmed, err := co.clientMap[toChain].IsSendOutTxProcessed(send.Index)
	if err != nil {
		log.Error().Err(err).Msg("IsSendOutTxProcessed error")
	}
	if confirmed {
		log.Info().Msgf("sendHash %s already processed; skip it", send.Index)
		return
	}

	signer := co.signerMap[toChain]
	message, err := base64.StdEncoding.DecodeString(send.Message)
	if err != nil {
		log.Err(err).Msgf("decode send.Message %s error", send.Message)
	}

	gasLimit := send.GasLimit
	if gasLimit == 0 {
		gasLimit = 250_000
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

	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done2:
				return
			default:
				included, confirmed, err := co.clientMap[toChain].IsSendOutTxProcessed(send.Index)
				if err != nil {
					log.Err(err).Msg("IsSendOutTxProcessed error")
				}
				if confirmed {
					log.Info().Msgf("sendHash %s already confirmed; skip it", send.Index)
					done <- true
					return
				}
				if included {
					log.Info().Msgf("sendHash %s already included but not yet confirmed. Keep monitoring", send.Index)
				}
				time.Sleep(24 * time.Second)
			}
		}
	}()

	// The following signing loop tries to sign outbound tx every 32 seconds.
	signTicker := time.NewTicker(time.Second)
SIGNLOOP:
	for range signTicker.C {
		select {
		case <-done:
			log.Info().Msg("breaking SignOutBoundTx loop: outbound already processed")
			break SIGNLOOP
		default:
			if time.Now().Second()%32 == int(sendhash[0])%32 {
				included, confirmed, err := co.clientMap[toChain].IsSendOutTxProcessed(send.Index)
				if err != nil {
					log.Error().Err(err).Msg("IsSendOutTxProcessed error")
				}
				if included {
					log.Info().Msgf("sendHash %s already included but not yet confirmed. will revisit", send.Index)
					continue
				}
				if confirmed {
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
				}
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
						log.Info().Msgf("broadcasting tx %s to chain %s: mint amount %d, nonce %d", outTxHash, toChain, amount, send.Nonce)
						err = signer.Broadcast(tx)
						if err != nil {
							log.Err(err).Msgf("Broadcast error: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
						} else {
							log.Err(err).Msgf("Broadcast success: nonce %d chain %s outTxHash %s", send.Nonce, toChain, outTxHash)
						}
					}
					// if outbound tx fails, kill this shepherd, a new one will be later spawned.
					if success := co.clientMap[toChain].WatchTxHashWithTimeout(outTxHash, send.Index); !success {
						time.Sleep(15 * time.Second) // wait until the receive confirm is voted on zetacore
						return
					}
				}
			}
		}
	}
}
