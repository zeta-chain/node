package sui

import (
	_ "embed"
	"encoding/base64"
)

//go:embed gateway.mv
var gatewayBinary []byte

//go:embed fake_usdc.mv
var fakeUSDC []byte

//go:embed evm.mv
var evmBinary []byte

//go:embed token.mv
var tokenBinary []byte

//go:embed connected.mv
var connectedBinary []byte

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

// ExampleTokenBytecodeBase64 gets the token binary encoded as base64 for deployment
func ExampleTokenBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(tokenBinary)
}

// ExampleConnectedBytecodeBase64 gets the connected binary encoded as base64 for deployment
func ExampleConnectedBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(connectedBinary)
}
