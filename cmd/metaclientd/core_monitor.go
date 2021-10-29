package metaclientd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func MonitorCore() {
	// create receive queue
	rxQueue := make([]*types.Receive, 0)

	bridge, err := createBridge()
	if err != nil {
		fmt.Println("error creating metachain bridge")
		return
	}

	// Pull from meta core and add to queue
	// TODO: Lock required?
	// TODO: Need some kind of waitgroup to prevent MonitorCore from
	// quitting?
	coreTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range coreTicker.C {
			rxList, err := bridge.GetAllReceive()
			if err != nil {
				fmt.Println("error requesting receives from metacore")
				return
			}

			// Add rxList items to queue
			for _, rx := range rxList {
				rxQueue = append(rxQueue, rx)
			}
		}
	}()

	// Pull items from queue
	go func() {
		for len(rxQueue) > 0 {
			// Pull the top
			rx := rxQueue[0]

			// Process

			// Discard top
			rxQueue = rxQueue[1:]
		}
	}()

}

// NOTE: Can we have a general "createBridge" function
// for the entire client instance?
func createBridge() (*MetachainBridge, error) {
	// TODO: How do we properly set these values?
	signerName := "alice"
	signerPass := "password"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &MetachainBridge{}, err
	}

	chainHomeFoler := filepath.Join(homeDir, ".metacore")

	kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		return &MetachainBridge{}, err
	}

	k := NewKeysWithKeybase(kb, signerName, signerPass)

	chainIP := "127.0.0.1"
	bridge, err := NewMetachainBridge(k, chainIP, "alice")
	if err != nil {
		return &MetachainBridge{}, err
	}

	return bridge, nil
}
