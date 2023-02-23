package types

import (
	"encoding/hex"
)

func BytesToEthHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
