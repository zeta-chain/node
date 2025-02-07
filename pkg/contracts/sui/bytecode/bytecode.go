// Package bytecode provides the full bytecode for the Sui gateway
package bytecode

import (
	_ "embed"
	"encoding/base64"
)

//go:embed gateway.mv
var gatewayBinary []byte

// GetEncodedGateway gets the gateway binary encoded as base64 for deployement
func GetEncodedGateway() string {
	return base64.StdEncoding.EncodeToString(gatewayBinary)
}
