package zetaclient

import (
	"errors"
	"github.com/rs/zerolog"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"os"
	"sort"
	"time"

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
	logger     zerolog.Logger
}

func NewCoreObserver(bridge *ZetaCoreBridge, signerMap map[common.Chain]*Signer, clientMap map[common.Chain]*ChainObserver, metrics *metrics.Metrics, tss *TSS) *CoreObserver {
	co := CoreObserver{}
	co.logger = log.With().Str("module", "CoreOb").Logger()
	co.tss = tss
	co.bridge = bridge
	co.signerMap = signerMap
	co.sendQueue = make([]*types.Send, 0)

	co.clientMap = clientMap
	co.metrics = metrics

	err := metrics.RegisterCounter(OUTBOUND_TX_SIGN_COUNT, "number of outbound tx signed")
	if err != nil {
		co.logger.Error().Err(err).Msg("error registering counter")
	}

	co.sendNew = make(chan *types.Send)
	co.sendDone = make(chan *types.Send)
	MAX_SIGNERS := 24 // assuming each signer takes 100s to finish (have outTx included), then throughput is bounded by 100/100 = 1 tx/s
	co.signerSlots = make(chan bool, MAX_SIGNERS)
	for i := 0; i < MAX_SIGNERS; i++ {
		co.signerSlots <- true
	}
	co.shepherds = make(map[string]bool)

	logFile, err := os.OpenFile("zetacore_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		// Can we log an error before we have our logger? :)
		co.logger.Error().Err(err).Msgf("there was an error creating a logFile on zetacore")
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
	go co.ShepherdManager()

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
					co.logger.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
					continue
				}
				// Keygen succeed! Report TSS address
				co.logger.Info().Msgf("Keygen success! keygen response: %v...", res)
				err = co.tss.SetPubKey(res.PubKey)
				if err != nil {
					co.logger.Error().Msgf("SetPubKey fail")
					continue
				}

				for _, chain := range config.ChainsEnabled {
					_, err = co.bridge.SetTSS(chain, co.tss.Address().Hex(), co.tss.PubkeyInBech32)
					if err != nil {
						co.logger.Error().Err(err).Msgf("SetTSS fail %s", chain)
					}
				}

				// Keysign test: sanity test
				co.logger.Info().Msgf("test keysign...")
				TestKeysign(co.tss.PubkeyInBech32, co.tss.Server)
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

// startObserve retrieves the pending list of Sends from ZetaCore every 10s
// for each new send, it tries to launch a send shepherd.
// the send shepherd makes sure the send is settled on all chains.
func (co *CoreObserver) startObserve() {
	observeTicker := time.NewTicker(12 * time.Second)
	for range observeTicker.C {
		sendList, err := co.bridge.GetAllPendingSend()
		if err != nil {
			co.logger.Error().Err(err).Msg("error requesting sends from zetacore")
			continue
		}
		if len(sendList) > 0 {
			co.logger.Info().Msgf("#pending send: %d", len(sendList))
		}
		sendMap := splitAndSortSendListByChain(sendList)
		for chain, sends := range sendMap {
			if len(sends) > 0 {
				co.logger.Info().Msgf("#pending sends on chain %s: %d nonce range [%d,%d]", chain, len(sends), sends[0].Nonce, sends[len(sends)-1].Nonce)
			}
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
