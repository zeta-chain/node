package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestBlockHeaderVerification_EnableChain(t *testing.T) {
	t.Run("should enable chain if chain not present", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		bhv.EnableChain(chains.BscMainnetChain.ChainId)
		require.True(t, bhv.IsChainEnabled(chains.BscMainnetChain.ChainId))
	})

	t.Run("should not enable chain is present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			EnabledChains: []types.EnabledChain{{ChainId: chains.BscMainnetChain.ChainId, Enabled: false}}}
		bhv.EnableChain(chains.BscMainnetChain.ChainId)
		require.True(t, bhv.IsChainEnabled(chains.BscMainnetChain.ChainId))
	})
}

func TestBlockHeaderVerification_DisableChain(t *testing.T) {
	t.Run("should disable chain if chain not present", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		bhv.DisableChain(chains.BscMainnetChain.ChainId)
		require.False(t, bhv.IsChainEnabled(chains.BscMainnetChain.ChainId))
	})

	t.Run("should disable chain if chain present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			EnabledChains: []types.EnabledChain{{ChainId: chains.BscMainnetChain.ChainId, Enabled: true}}}
		bhv.DisableChain(chains.BscMainnetChain.ChainId)
		require.False(t, bhv.IsChainEnabled(chains.BscMainnetChain.ChainId))
	})
}

func TestBlockHeaderVerification_IsChainEnabled(t *testing.T) {
	t.Run("should return true if chain is enabled", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		require.True(t, bhv.IsChainEnabled(1))
	})

	t.Run("should return false if chain is disabled", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			EnabledChains: []types.EnabledChain{{ChainId: 1, Enabled: false}}}
		require.False(t, bhv.IsChainEnabled(1))
	})

	t.Run("should return false if chain is not present", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		require.False(t, bhv.IsChainEnabled(1000))
	})
}
