package common

import (
	"os"
	"strings"
)

// ChainNetwork is to indicate which chain environment THORNode are working with
type ChainNetwork uint8

const (
	// TestNet network for test
	TestNet ChainNetwork = iota
	// MainNet network for main net
	MainNet
	// MockNet network for main net
	MockNet
)

// GetCurrentChainNetwork determinate what kind of network currently it is working with
func GetCurrentChainNetwork() ChainNetwork {
	if strings.EqualFold(os.Getenv("NET"), "mocknet") {
		return MockNet
	}
	if strings.EqualFold(os.Getenv("NET"), "testnet") {
		return TestNet
	}
	return MainNet
}

// Soft Equals check is mainnet == mainet, or (testnet/mocknet == testnet/mocknet)
func (net ChainNetwork) SoftEquals(net2 ChainNetwork) bool {
	if net == MainNet && net2 == MainNet {
		return true
	}
	if net != MainNet && net2 != MainNet {
		return true
	}
	return false
}
