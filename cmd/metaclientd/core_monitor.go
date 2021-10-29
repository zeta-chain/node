package metaclientd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type CoreMonitor struct {
	sendQueue []*types.Send
	bridge    *MetachainBridge
	signer    *Signer
}

func (cm *CoreMonitor) InitCoreMonitor(bridge *MetachainBridge, signer *Signer) {
	cm.bridge = bridge
	cm.signer = signer
	cm.sendQueue = make([]*types.Send, 0)
}

func (cm *CoreMonitor) MonitorCore() {
	// Pull from meta core and add to queue
	// TODO: Lock required?
	// TODO: Need some kind of waitgroup to prevent MonitorCore from
	// quitting?
	coreTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range coreTicker.C {
			sendList, err := cm.bridge.GetAllSend()
			if err != nil {
				fmt.Println("error requesting receives from metacore")
				return
			}

			// Add rxList items to queue
			for _, send := range sendList {
				cm.sendQueue = append(cm.sendQueue, send)
			}
		}
	}()

	// Pull items from queue
	go func() {
		for len(cm.sendQueue) > 0 {
			// Pull the top
			send := cm.sendQueue[0]

			// TODO: How to pull the data below off send
			fmt.Println(send)

			// Process
			var amount *big.Int
			var to ethcommon.Address
			var gasLimit uint64
			var message []byte

			outTxHash, err := cm.signer.MMint(
				amount,
				to,
				gasLimit,
				message,
			)
			if err != nil {
				fmt.Println("error minting received transaction")
				return
			}

			// TODO: What to do with outTxHash now?
			fmt.Println(outTxHash)

			// TODO: We now have outTxHash and sendHash (from send)
			// How do we save this for use in observer?

			// Discard top
			cm.sendQueue = cm.sendQueue[1:]
		}
	}()

}
