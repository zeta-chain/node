package metaclient

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math/big"
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
	sendQueue []*types.Send
	sendMap   map[string]*types.Send
	sendStatus map[string]TxStatus
	bridge    *MetachainBridge
	signer    *Signer
}

func  NewCoreObserver(bridge *MetachainBridge, signer *Signer) *CoreObserver{
	co := CoreObserver{}
	co.bridge = bridge
	co.signer = signer
	co.sendQueue = make([]*types.Send, 0)
	co.sendMap = make(map[string]*types.Send)
	co.sendStatus = make(map[string]TxStatus)
	return &co
}

func (co *CoreObserver) MonitorCore() {
	go func() {
		sendList, err := co.bridge.GetAllSend()
		for {
			select {
			default:
				if err != nil {
					fmt.Println("error requesting receives from metacore")
					time.Sleep(5 * time.Second)
					continue
				}
				for _, send := range sendList {
					if types.SendStatus_name[int32(send.Status)] == "Finalized" {
						if _, found := co.sendMap[send.Index]; !found {
							log.Info().Msgf("New send queued with finalized block %d", send.FinalizedMetaHeight)
							co.sendMap[send.Index] = send
							co.sendQueue = append(co.sendQueue, send)
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}

	}()

	// Pull items from queue
	go func() {
		for {
			if len(co.sendQueue) > 0 {
				send := co.sendQueue[0]
				amount, ok := new(big.Int).SetString(send.MBurnt, 10)
				if !ok {
					fmt.Println("error converting MBurnt to big.Int")
					time.Sleep(5 * time.Second)
					continue
				}
				to := ethcommon.HexToAddress(send.Receiver)
				message := []byte(send.Message)

				// TODO: Eventually this should come from smart contract
				var gasLimit uint64 = 80000

				outTxHash, err := co.signer.MMint(amount, to, gasLimit, message)
				if err != nil {
					fmt.Println("error minting received transaction")
					time.Sleep(5 * time.Second)
					continue
				}
				co.sendStatus[send.Index] = Pending

				// TODO: We now have outTxHash and sendHash (from send)
				// How do we save this for use in observer?
				fmt.Println("sendHash: ", send.Index)
				fmt.Println("outTxHash: ", outTxHash)

			}
			time.Sleep(5*time.Second)
		}

	}()

}
