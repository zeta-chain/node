//go:build drain

package main

import (
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestOperatorPubKeyFingerprint(t *testing.T) {
	t.Run("rejects all-zero placeholder", func(t *testing.T) {
		_, err := operatorPubKeyFingerprint(make([]byte, 33))
		require.Error(t, err)
	})

	t.Run("returns fingerprint for a valid key", func(t *testing.T) {
		priv, err := ethcrypto.GenerateKey()
		require.NoError(t, err)
		fp, err := operatorPubKeyFingerprint(ethcrypto.CompressPubkey(&priv.PublicKey))
		require.NoError(t, err)
		require.Equal(t, ethcrypto.PubkeyToAddress(priv.PublicKey).Hex(), fp)
	})

	t.Run("rejects garbage", func(t *testing.T) {
		_, err := operatorPubKeyFingerprint([]byte{1, 2, 3})
		require.Error(t, err)
	})
}

func TestDrainWindowFromEnv(t *testing.T) {
	require.Equal(t, int64(drainWindow), drainWindowFromEnv())

	t.Setenv(envDrainWindow, "500")
	require.EqualValues(t, 500, drainWindowFromEnv())

	t.Setenv(envDrainWindow, "not-a-number")
	require.Equal(t, int64(drainWindow), drainWindowFromEnv())
}
