package types

import (
	"encoding/hex"
)

// EthHexToBytes converts an Ethereum hex string to bytes
func BytesToEthHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
