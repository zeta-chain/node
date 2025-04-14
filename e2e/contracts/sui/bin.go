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

//go:embed example.mv
var exampleBinary []byte

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

// TokenBytecodeBase64 gets the token binary encoded as base64 for deployment
func TokenBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(tokenBinary)
}

// ExampleBytecodeBase64 gets the example binary encoded as base64 for deployment
func ExampleBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(exampleBinary)
}
