package metaclientd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type CoreObserver struct {
	sendQueue []*types.Send
	bridge    *MetachainBridge
	signer    *Signer
}

func (co *CoreObserver) InitCoreObserver(bridge *MetachainBridge, signer *Signer) {
	co.bridge = bridge
	co.signer = signer
	co.sendQueue = make([]*types.Send, 0)
}

func (co *CoreObserver) MonitorCore() {
	// Pull from meta core and add to queue
	// TODO: Lock required?
	// TODO: Need some kind of waitgroup to prevent MonitorCore from
	// quitting?
	coreTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range coreTicker.C {
			sendList, err := co.bridge.GetAllSend()
			if err != nil {
				fmt.Println("error requesting receives from metacore")
				return
			}

			// Add sendList items to queue if status is finalized
			// TODO: extra check to make sure we don't double add?
			// ask @pwu
			for _, send := range sendList {
				if types.SendStatus_name[int32(send.Status)] == "Finalized" {
					co.sendQueue = append(co.sendQueue, send)
				}
			}
		}
	}()

	// Pull items from queue
	go func() {
		for range coreTicker.C {
			for len(co.sendQueue) > 0 {
				// Pull the top
				send := co.sendQueue[0]

				// Process
				amount, ok := new(big.Int).SetString(send.MBurnt, 10)
				if !ok {
					fmt.Println("error converting MBurnt to big.Int")
					return
				}
				to := ethcommon.HexToAddress(send.Receiver)
				message := []byte(send.Message)

				// Gas limit hard-coded to 80k for now
				// TODO: Eventually this should come from smart contract
				var gasLimit uint64 = 80000

				outTxHash, err := co.signer.MMint(
					amount,
					to,
					gasLimit,
					message,
				)
				if err != nil {
					fmt.Println("error minting received transaction")
					return
				}

				// TODO: We now have outTxHash and sendHash (from send)
				// How do we save this for use in observer?
				fmt.Println("sendHash: ", send.Index)
				fmt.Println("outTxHash: ", outTxHash)

				// Discard top
				co.sendQueue = co.sendQueue[1:]
			}
		}
	}()

}
