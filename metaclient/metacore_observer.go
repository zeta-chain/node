package metaclient

import (
	"encoding/hex"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
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
)

type CoreObserver struct {
	sendQueue  []*types.Send
	sendMap    map[string]*types.Send
	sendStatus map[string]TxStatus
	bridge     *MetachainBridge
	signerMap  map[common.Chain]*Signer
	lock       sync.Mutex
}

func NewCoreObserver(bridge *MetachainBridge, signerMap map[common.Chain]*Signer) *CoreObserver {
	co := CoreObserver{}
	co.bridge = bridge
	co.signerMap = signerMap
	co.sendQueue = make([]*types.Send, 0)
	co.sendMap = make(map[string]*types.Send)
	co.sendStatus = make(map[string]TxStatus)
	return &co
}

func (co *CoreObserver) MonitorCore() {
	myid := co.bridge.keys.GetSignerInfo().GetAddress().String()
	log.Info().Msgf("MonitorCore started by signer %s", myid)
	go func() {
		for {
			sendList, err := co.bridge.GetAllSend()
			if err != nil {
				fmt.Println("error requesting receives from metacore")
				time.Sleep(5 * time.Second)
				continue
			}
			for _, send := range sendList {
				if types.SendStatus_name[int32(send.Status)] == "Finalized" {
					if _, found := co.sendMap[send.Index]; !found {
						co.lock.Lock()
						log.Info().Msgf("New send queued with finalized block %d", send.FinalizedMetaHeight)
						co.sendMap[send.Index] = send
						co.sendQueue = append(co.sendQueue, send)
						co.sendStatus[send.Index] = Unprocessed
						co.lock.Unlock()
					}
				} else if types.SendStatus_name[int32(send.Status)] == "Mined" {
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
					to := ethcommon.HexToAddress(send.Receiver)
					toChain, err := common.ParseChain(send.ReceiverChain)
					if err != nil {
						log.Err(err).Msg("ParseChain fail; skip")
						time.Sleep(5 * time.Second)
						continue
					}
					signer := co.signerMap[toChain]
					message := []byte(send.Message)

					var gasLimit uint64 = 80000

					log.Info().Msgf("chain %s minting %d to %s", toChain, amount, to.Hex())
					if send.Signers[send.Broadcaster] == myid {
						sendHash, err := hex.DecodeString(send.Index[2:]) // remove the leading 0x
						if err != nil || len(sendHash) != 32 {
							log.Err(err).Msgf("decode sendHash %s error", send.Index)
						}
						var sendhash [32]byte
						copy(sendhash[:32], sendHash[:32])
						outTxHash, err := signer.MMint(amount, to, gasLimit, message, sendhash, send.Nonce)
						if err != nil {
							log.Err(err).Msg("singer %s error minting received transaction")
						}
						fmt.Printf("sendHash: %s, outTxHash %s signer %s\n", send.Index[:6], outTxHash, myid)
					}
					co.sendStatus[send.Index] = Pending // do not process this; other signers might already done it

					//for {
					//	tx, isPending, err := signer.client.TransactionByHash(context.Background(), ethcommon.HexToHash(outTxHash))
					//	if err != nil {
					//		log.Warn().Msgf("TransactionByHash %s err %s", outTxHash, err)
					//		time.Sleep(2*time.Second)
					//		continue
					//	}
					//	if !isPending {
					//		receipt, err := signer.client.TransactionReceipt(context.Background(), tx.Hash())
					//		if err != nil {
					//			log.Err(err).Msg("TransactionReceipt")
					//		}
					//		if receipt.Status == 1 { // success execution
					//			fmt.Printf("PostReceive %s %s %d\n", send.Index[:6], outTxHash[:6], receipt.BlockNumber.Uint64())
					//			metahash, err := co.bridge.PostReceiveConfirmation(send.Index, outTxHash, receipt.BlockNumber.Uint64(), "111")
					//			if err != nil {
					//				log.Err(err).Msgf("PostReceiveConfirmation metahash %s", metahash)
					//			}
					//		}
					//		break
					//	}
					//	time.Sleep(2*time.Second)
					//}

				}

			}
			time.Sleep(5 * time.Second)
		}

	}()

}
