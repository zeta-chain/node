package utils

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// ScriptPKToAddress is a hex string for P2WPKH script
func ScriptPKToAddress(scriptPKHex string, params *chaincfg.Params) string {
	pkh, err := hex.DecodeString(scriptPKHex[4:])
	if err == nil {
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pkh, params)
		if err == nil {
			return addr.EncodeAddress()
		}
	}
	return ""
}

type infoLogger interface {
	Info(message string, args ...interface{})
}

type NoopLogger struct{}

func (nl NoopLogger) Info(_ string, _ ...interface{}) {}
