package tss

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
)

func TestPubKey(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		cases := []struct {
			name   string
			input  string
			errMsg string
		}{
			{"empty string", "", "empty bech32 address"},
			{"invalid prefix", "invalid1addwnpepq...", "unable to GetPubKeyFromBech32"},
			{"malformed bech32", "zetapub1invalid", "decoding bech32 failed"},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				_, err := NewPubKeyFromBech32(tt.input)
				require.ErrorContains(t, err, tt.errMsg)
			})
		}
	})

	t.Run("Valid NewPubKeyFromBech32", func(t *testing.T) {
		// ARRANGE
		const sample = `zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc`

		// ACT
		pk, err := NewPubKeyFromBech32(sample)

		// ASSERT
		require.NoError(t, err)
		assert.NotEmpty(t, pk)

		addrEVM := pk.AddressEVM()
		addrBTC, err := pk.AddressBTC(chains.BitcoinMainnet.ChainId)
		require.NoError(t, err)

		expectedPkScript, err := txscript.PayToAddrScript(addrBTC)
		require.NoError(t, err)
		pkScript, err := pk.BTCPayToAddrScript(chains.BitcoinMainnet.ChainId)
		require.NoError(t, err)
		assert.Equal(t, expectedPkScript, pkScript)

		assert.Equal(t, sample, pk.Bech32String())
		assert.Equal(t, "0x70e967acfcc17c3941e87562161406d41676fd83", strings.ToLower(addrEVM.Hex()))
		assert.Equal(t, "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y", addrBTC.String())

		// Check that NewPubKeyFromECDSA works
		pk2, err := NewPubKeyFromECDSA(*pk.ecdsaPubKey)
		require.NoError(t, err)
		require.Equal(t, pk.Bech32String(), pk2.Bech32String())
	})

	t.Run("Valid NewPubKeyFromECDSAHexString", func(t *testing.T) {
		// ARRANGE
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)

		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		evmAddr := crypto.PubkeyToAddress(pk.PublicKey)

		// ACT
		actual, err := NewPubKeyFromECDSAHexString(pubKeyHex)

		// ASSERT
		require.NoError(t, err)
		assert.Equal(t, evmAddr, actual.AddressEVM())
		assert.True(t, strings.HasPrefix(actual.Bech32String(), "zetapub"))

		t.Run("With 0x prefix", func(t *testing.T) {
			// ACT
			actual2, err := NewPubKeyFromECDSAHexString("0x" + pubKeyHex)

			// ASSERT
			require.NoError(t, err)
			assert.Equal(t, actual.Bech32String(), actual2.Bech32String())
		})
	})
}
