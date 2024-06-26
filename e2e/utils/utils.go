package utils

import (
	"context"
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"
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

type testingKey struct{}

// WithTesting allows to store a testing.T instance in the context
func WithTesting(ctx context.Context, t require.TestingT) context.Context {
	return context.WithValue(ctx, testingKey{}, t)
}

// TestingFromContext extracts require.TestingT from the context or panics.
func TestingFromContext(ctx context.Context) require.TestingT {
	t, ok := ctx.Value(testingKey{}).(require.TestingT)
	if !ok {
		panic("context missing require.TestingT key")
	}

	return t
}
