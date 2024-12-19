package common

import (
	"os"
	"strconv"
)

const (
	// EnvEnableLiveTest is the environment variable to enable zetaclient live tests
	EnvEnableLiveTest = "ENABLE_LIVE_TEST"

	// EnvBtcRPCMainnet is the environment variable to enable mainnet for bitcoin rpc
	EnvBtcRPCMainnet = "BTC_RPC_MAINNET"

	// EnvBtcRPCTestnet is the environment variable to enable testnet for bitcoin rpc
	EnvBtcRPCTestnet = "BTC_RPC_TESTNET"
)

// LiveTestEnabled returns true if live tests are enabled
func LiveTestEnabled() bool {
	value := os.Getenv(EnvEnableLiveTest)

	// parse flag
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return enabled
}
