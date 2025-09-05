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

	// EnvBtcRPCSignet is the environment variable to enable signet for bitcoin rpc
	EnvBtcRPCSignet = "BTC_RPC_SIGNET"

	// EnvBtcRPCTestnet4 is the environment variable to enable testnet4 for bitcoin rpc
	EnvBtcRPCTestnet4 = "BTC_RPC_TESTNET4"

	// EnvTONRPC is the environment variable to enable ton rpc
	EnvTONRPC = "TON_RPC"

	// BTC_TEST_PK is the environment variable to configure bitcoin private key for live unit tests
	EnvBTCTestPK = "BTC_TEST_PK"
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

// IsEnvVariableSet returns true if the environment variable is set
func IsEnvVariableSet(envVar string) bool {
	value := os.Getenv(envVar)
	return value != ""
}
