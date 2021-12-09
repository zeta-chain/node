package metaclient

import (
	"encoding/hex"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/Meta-Protocol/metacore/metaclient/config"
	"github.com/rs/zerolog/log"
	"math/big"
	"sync"
	"time"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type TxStatus int64

const (
	Unprocessed TxStatus = iota
	Pending
	Mined
	Confirmed
	Error
)

type CoreObserver struct {
	sendQueue  []*types.Send
	sendMap    map[string]*types.Send // send.Index => send
	sendStatus map[string]TxStatus    // send.Index => status
	recvMap    map[string]*types.Receive
	recvStatus map[string]TxStatus
	bridge     *MetachainBridge
	signerMap  map[common.Chain]*Signer
	clientMap  map[common.Chain]*ChainObserver
	lock       sync.Mutex
}

func NewCoreObserver(bridge *MetachainBridge, signerMap map[common.Chain]*Signer, clientMap map[common.Chain]*ChainObserver) *CoreObserver {
	co := CoreObserver{}
	co.bridge = bridge
	co.signerMap = signerMap
	co.sendQueue = make([]*types.Send, 0)
	co.sendMap = make(map[string]*types.Send)
	co.sendStatus = make(map[string]TxStatus)
	co.recvMap = make(map[string]*types.Receive)
	co.recvStatus = make(map[string]TxStatus)
	co.clientMap = clientMap
	return &co
}

func (co *CoreObserver) MonitorCore() {
	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	log.Info().Msgf("MonitorCore started by signer %s", myid)
	go func() {
		for {
			zetaHeight, err := co.bridge.GetMetaBlockHeight()
			if err != nil {
				log.Warn().Err(err).Msgf("GetMetaBlockHeight error")
				continue
			}

			sendList, err := co.bridge.GetAllPendingSend()
			if err != nil {
				fmt.Println("error requesting sends from metacore")
				time.Sleep(5 * time.Second)
				continue
			}

			for _, send := range sendList {
				if send.Status == types.SendStatus_Finalized || send.Status == types.SendStatus_Abort {
					oldSend, found := co.sendMap[send.Index]
					if !found || oldSend.Status != send.Status {
						co.lock.Lock()
						if !found {
							log.Debug().Msgf("New send queued with finalized block %d", send.FinalizedMetaHeight)
						}
						if found && oldSend.Status != send.Status {
							log.Debug().Msgf("Old send status updated from %s to %s", types.SendStatus_name[int32(oldSend.Status)], types.SendStatus_name[int32(send.Status)])
						}
						co.sendMap[send.Index] = send
						co.sendQueue = append(co.sendQueue, send)
						co.sendStatus[send.Index] = Unprocessed
						co.lock.Unlock()
					} else {
						status, found := co.sendStatus[send.Index]
						if !found {
							log.Error().Msgf("status of send: %s not found", send.Index)
							continue
						}
						// the send is not successfully process; re-process is needed
						if zetaHeight-send.FinalizedMetaHeight > config.TIMEOUT_THRESHOLD_FOR_RETRY &&
							zetaHeight-send.FinalizedMetaHeight < 2*config.TIMEOUT_THRESHOLD_FOR_RETRY &&
							status == Pending {
							log.Warn().Msgf("Timeout send: sendHash %s; re-processs...", send.Index)
							co.lock.Lock()
							co.sendStatus[send.Index] = Unprocessed
							co.lock.Unlock()
						}

					}

				} else if send.Status == types.SendStatus_Mined {
					if send, found := co.sendMap[send.Index]; found {
						co.lock.Lock()
						if co.sendStatus[send.Index] != Mined {
							log.Info().Msgf("Send status changed to Mined")
						}
						co.sendStatus[send.Index] = Mined
						co.lock.Unlock()
					}
				}
			}
			time.Sleep(5 * time.Second)

		}

	}()

	go func() {
		for {
			recvList, err := co.bridge.GetAllReceive()
			if err != nil {
				fmt.Println("error requesting receives from metacore")
				time.Sleep(5 * time.Second)
				continue
			}
			for _, recv := range recvList {
				if recv.Status == common.ReceiveStatus_Created {
					if _, found := co.recvMap[recv.Index]; !found {
						co.lock.Lock()
						log.Debug().Msgf("New recv created with finalized block")
						co.recvMap[recv.Index] = recv
						co.recvStatus[recv.Index] = Unprocessed
						chain, err := common.ParseChain(recv.Chain)
						if err != nil {
							fmt.Printf("recv chain invalid: %s\n", recv.Chain)
							continue
						}
						co.clientMap[chain].AddTxToWatchList(recv.OutTxHash, recv.SendHash)
						co.lock.Unlock()
					}
				}
			}
			time.Sleep(5 * time.Second)

		}

	}()

	// Pull items from queue
	go func() {
		for {
			if len(co.sendQueue) > 0 {
				for _, send := range co.sendQueue {
					if co.sendStatus[send.Index] != Unprocessed {
						continue
					}
					amount, ok := new(big.Int).SetString(send.MMint, 10)
					if !ok {
						log.Error().Msg("error converting MBurnt to big.Int")
						time.Sleep(5 * time.Second)
						continue
					}

					var to ethcommon.Address
					var err error
					var toChain common.Chain
					if send.Status == types.SendStatus_Abort {
						to = ethcommon.HexToAddress(send.Sender)
						toChain, err = common.ParseChain(send.SenderChain)
						log.Info().Msgf("Abort: reverting inbound")
					} else {
						to = ethcommon.HexToAddress(send.Receiver)
						toChain, err = common.ParseChain(send.ReceiverChain)
					}
					if err != nil {
						log.Err(err).Msg("ParseChain fail; skip")
						time.Sleep(5 * time.Second)
						continue
					}
					signer := co.signerMap[toChain]
					message := []byte(send.Message)

					var gasLimit uint64 = 90_000

					log.Info().Msgf("chain %s minting %d to %s", toChain, amount, to.Hex())
					sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
					if err != nil || len(sendHash) != 32 {
						log.Err(err).Msgf("decode sendHash %s error", send.Index)
					}
					var sendhash [32]byte
					copy(sendhash[:32], sendHash[:32])
					gasprice, ok := new(big.Int).SetString(send.GasPrice, 10)
					if !ok {
						log.Err(err).Msgf("cannot convert gas price  %s ", send.GasPrice)
					}
					tx, err := signer.SignOutboundTx(amount, to, gasLimit, message, sendhash, send.Nonce, gasprice)
					if err != nil {
						log.Err(err).Msgf("MMint error: nonce %d", send.Nonce)
						co.sendStatus[send.Index] = Error // do not process this; other signers might already done it
						continue
					}
					outTxHash := tx.Hash().Hex()
					fmt.Printf("sendHash: %s, outTxHash %s signer %s\n", send.Index[:6], outTxHash, myid)
					if send.Signers[send.Broadcaster] == myid {
						err := signer.Broadcast(tx)
						if err != nil {
							log.Err(err).Msgf("Broadcast error: nonce %d", send.Nonce)
						}
					}
					co.sendStatus[send.Index] = Pending // do not process this; other signers might already done it
					_, err = co.bridge.PostReceiveConfirmation(send.Index, outTxHash, 0, amount.String(), common.ReceiveStatus_Created, send.ReceiverChain)
					if err != nil {
						log.Err(err).Msgf("PostReceiveConfirmation of just created receive")
					}
					co.clientMap[toChain].AddTxToWatchList(outTxHash, send.Index)

				}

			}
			time.Sleep(5 * time.Second)
		}

	}()

}
