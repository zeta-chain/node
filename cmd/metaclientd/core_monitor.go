package metaclientd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type CoreMonitor struct {
	rxQueue []*types.Receive
	bridge  *MetachainBridge
	signer  *Signer
}

func (cm *CoreMonitor) InitCoreMonitor(bridge *MetachainBridge, signer *Signer) {
	cm.bridge = bridge
	cm.signer = signer
	cm.rxQueue = make([]*types.Receive, 0)
}

func (cm *CoreMonitor) MonitorCore() {
	// Pull from meta core and add to queue
	// TODO: Lock required?
	// TODO: Need some kind of waitgroup to prevent MonitorCore from
	// quitting?
	coreTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range coreTicker.C {
			rxList, err := cm.bridge.GetAllReceive()
			if err != nil {
				fmt.Println("error requesting receives from metacore")
				return
			}

			// Add rxList items to queue
			for _, rx := range rxList {
				cm.rxQueue = append(cm.rxQueue, rx)
			}
		}
	}()

	// Pull items from queue
	go func() {
		for len(cm.rxQueue) > 0 {
			// Pull the top
			rx := cm.rxQueue[0]

			// TODO: How to pull the data below off rx
			fmt.Println(rx)

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

			// Discard top
			cm.rxQueue = cm.rxQueue[1:]
		}
	}()

}
