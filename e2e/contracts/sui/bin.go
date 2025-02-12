package sui

import (
	_ "embed"
	"encoding/base64"
)

//go:embed gateway.mv
var gatewayBinary []byte

// GatewayBytecodeBase64 gets the gateway binary encoded as base64 for deployment
func GatewayBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(gatewayBinary)
}
