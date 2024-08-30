package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestBlockHeaderVerification_Validate(t *testing.T) {
	t.Run("should return nil if no duplicate chain id", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: 1, Enabled: true},
				{ChainId: 2, Enabled: true},
			}}
		require.NoError(t, bhv.Validate())
	})

	t.Run("should return error if duplicate chain id", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: 1, Enabled: true},
				{ChainId: 1, Enabled: true},
			}}
		require.Error(t, bhv.Validate())
	})
}
func TestBlockHeaderVerification_EnableChain(t *testing.T) {
	t.Run("should enable chain if chain not present", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		bhv.EnableChain(chains.BscMainnet.ChainId)
		require.True(t, bhv.IsChainEnabled(chains.BscMainnet.ChainId))
	})

	t.Run("should not enable chain is present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: chains.BscMainnet.ChainId, Enabled: false},
			}}
		bhv.EnableChain(chains.BscMainnet.ChainId)
		require.True(t, bhv.IsChainEnabled(chains.BscMainnet.ChainId))
	})
}

func TestBlockHeaderVerification_DisableChain(t *testing.T) {
	t.Run("should disable chain if chain not present", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		bhv.DisableChain(chains.BscMainnet.ChainId)
		require.False(t, bhv.IsChainEnabled(chains.BscMainnet.ChainId))
	})

	t.Run("should disable chain if chain present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: chains.BscMainnet.ChainId, Enabled: true},
			}}
		bhv.DisableChain(chains.BscMainnet.ChainId)
		require.False(t, bhv.IsChainEnabled(chains.BscMainnet.ChainId))
	})
}

func TestBlockHeaderVerification_IsChainEnabled(t *testing.T) {
	t.Run("should return true if chain is enabled", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		require.True(t, bhv.IsChainEnabled(1))
	})

	t.Run("should return false if chain is disabled", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{{ChainId: 1, Enabled: false}}}
		require.False(t, bhv.IsChainEnabled(1))
	})

	t.Run("should return false if chain is not present", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		require.False(t, bhv.IsChainEnabled(1000))
	})
}

func TestBlockHeaderVerification_GetEnabledChainIDList(t *testing.T) {
	t.Run("should return list of enabled chain ids", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		enabledChains := bhv.GetHeaderEnabledChainIDs()
		require.Len(t, enabledChains, 2)
		require.Contains(t, enabledChains, int64(1))
		require.Contains(t, enabledChains, int64(2))
	})

	t.Run("should return empty list if no chain is enabled", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: 1, Enabled: false},
				{ChainId: 2, Enabled: false},
			}}
		enabledChains := bhv.GetHeaderEnabledChainIDs()
		require.Len(t, enabledChains, 0)
	})

	t.Run("should return empty list if no chain is present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{}
		enabledChains := bhv.GetHeaderEnabledChainIDs()
		require.Len(t, enabledChains, 0)
	})
}

func TestBlockHeaderVerification_GetEnabledChainsList(t *testing.T) {
	t.Run("should return list of enabled chains", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		enabledChains := bhv.GetHeaderEnabledChains()
		require.Len(t, enabledChains, 2)
		require.Contains(t, enabledChains, types.HeaderSupportedChain{ChainId: 1, Enabled: true})
		require.Contains(t, enabledChains, types.HeaderSupportedChain{ChainId: 2, Enabled: true})
	})

	t.Run("should return empty list if no chain is enabled", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: 1, Enabled: false},
				{ChainId: 2, Enabled: false},
			}}
		enabledChains := bhv.GetHeaderEnabledChains()
		require.Len(t, enabledChains, 0)
	})

	t.Run("should return empty list if no chain is present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{}
		enabledChains := bhv.GetHeaderEnabledChains()
		require.Len(t, enabledChains, 0)
	})
}

func TestBlockHeaderVerification_GetSupportedChainsList(t *testing.T) {
	t.Run("should return list of supported chains", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		supportedChains := bhv.GetHeaderSupportedChainsList()
		require.Len(t, supportedChains, 2)
		require.Contains(t, supportedChains, types.HeaderSupportedChain{ChainId: 1, Enabled: true})
		require.Contains(t, supportedChains, types.HeaderSupportedChain{ChainId: 2, Enabled: true})
	})

	t.Run("should return empty list if no chain is present", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{}
		supportedChains := bhv.GetHeaderSupportedChainsList()
		require.Len(t, supportedChains, 0)
	})

	t.Run("should items even if chain is not enabled but still supported", func(t *testing.T) {
		bhv := types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{ChainId: 1, Enabled: false},
				{ChainId: 2, Enabled: false},
			}}
		supportedChains := bhv.GetHeaderSupportedChains()
		require.Len(t, supportedChains, 2)
		require.Contains(t, supportedChains, types.HeaderSupportedChain{ChainId: 1, Enabled: false})
		require.Contains(t, supportedChains, types.HeaderSupportedChain{ChainId: 2, Enabled: false})
	})
}
