package sui

import (
	_ "embed"
	"encoding/base64"
	"os"

	"github.com/stretchr/testify/require"
)

//go:embed gateway.mv
var gatewayBinary []byte

//go:embed fake_usdc.mv
var fakeUSDC []byte

//go:embed evm.mv
var evmBinary []byte

//go:embed arbi_call_token.mv
var arbiCallTokenBinary []byte

//go:embed arbi_call_connected.mv
var arbiCallConnectedBinary []byte

// GatewayBytecodeBase64 gets the gateway binary encoded as base64 for deployment
func GatewayBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(gatewayBinary)
}

// FakeUSDCBytecodeBase64 gets the fake USDC binary encoded as base64 for deployment
func FakeUSDCBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(fakeUSDC)
}

// EVMBytecodeBase64 gets the EVM binary encoded as base64 for deployment
func EVMBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(evmBinary)
}

// ExampleArbiCallFungibleTokenBytecodeBase64 gets the arbitrary call example package's fungible token binary encoded as base64 for deployment
func ExampleArbiCallFungibleTokenBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(arbiCallTokenBinary)
}

// ExampleArbiCallConnectedBytecodeBase64 gets the arbitrary call example package's connected binary encoded as base64 for deployment
func ExampleArbiCallConnectedBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(arbiCallConnectedBinary)
}

// ReadMoveBinaryBase64 reads a given move binary file and returns it as base64 encoded string
func ReadMoveBinaryBase64(t require.TestingT, binaryName string) string {
	binaryBytes, err := os.ReadFile(binaryName)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(binaryBytes)
}
