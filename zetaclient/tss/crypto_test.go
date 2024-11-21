package tss

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd"
	"github.com/zeta-chain/node/pkg/chains"
)

func TestPubKey(t *testing.T) {
	cmd.SetupCosmosConfig()

	t.Run("Invalid", func(t *testing.T) {
		_, err := NewPubKeyFromBech32("")
		require.ErrorContains(t, err, "empty bech32 address")
	})

	t.Run("Valid", func(t *testing.T) {
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

		assert.Equal(t, sample, pk.Bech32String())
		assert.Equal(t, "0x70e967acfcc17c3941e87562161406d41676fd83", strings.ToLower(addrEVM.Hex()))
		assert.Equal(t, "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y", addrBTC.String())

		// Check that NewPubKeyFromECDSA works
		pk2, err := NewPubKeyFromECDSA(*pk.ecdsaPubKey)
		require.NoError(t, err)
		require.Equal(t, pk.Bech32String(), pk2.Bech32String())
	})
}
