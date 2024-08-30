package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestKeeper_GetBlockHeaderVerification(t *testing.T) {
	t.Run("can get all verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		bhv := sample.BlockHeaderVerification()

		k.SetBlockHeaderVerification(ctx, bhv)

		blockHeaderVerification, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.Len(t, blockHeaderVerification.HeaderSupportedChains, 2)
		require.Equal(t, bhv, blockHeaderVerification)
	})

	t.Run("return empty list when no flags are set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		blockHeaderVerification, found := k.GetBlockHeaderVerification(ctx)
		require.False(t, found)
		require.Len(t, blockHeaderVerification.HeaderSupportedChains, 0)
		require.Equal(t, types.BlockHeaderVerification{}, blockHeaderVerification)
	})
}

func TestKeeper_CheckVerificationFlagsEnabled(t *testing.T) {
	t.Run("can check verification flags with ethereum enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: true,
				},
			},
		})

		err := k.CheckBlockHeaderVerificationEnabled(ctx, chains.Ethereum.ChainId)
		require.NoError(t, err)

		err = k.CheckBlockHeaderVerificationEnabled(ctx, chains.BitcoinMainnet.ChainId)
		require.Error(t, err)
		require.ErrorContains(
			t,
			err,
			fmt.Sprintf("proof verification is disabled for chain %d", chains.BitcoinMainnet.ChainId),
		)

		err = k.CheckBlockHeaderVerificationEnabled(ctx, 1000)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("proof verification is disabled for chain %d", 1000))
	})

	t.Run("can check verification flags with bitcoin enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
			},
		})

		err := k.CheckBlockHeaderVerificationEnabled(ctx, chains.Ethereum.ChainId)
		require.Error(t, err)
		require.ErrorContains(
			t,
			err,
			fmt.Sprintf("proof verification is disabled for chain %d", chains.Ethereum.ChainId),
		)

		err = k.CheckBlockHeaderVerificationEnabled(ctx, chains.BitcoinMainnet.ChainId)
		require.NoError(t, err)

		err = k.CheckBlockHeaderVerificationEnabled(ctx, 1000)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("proof verification is disabled for chain %d", 1000))
	})

	t.Run("check returns false if flag is not set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		err := k.CheckBlockHeaderVerificationEnabled(ctx, chains.Ethereum.ChainId)
		require.ErrorContains(t, err, "proof verification is disabled for all chains")

		err = k.CheckBlockHeaderVerificationEnabled(ctx, chains.BitcoinMainnet.ChainId)
		require.ErrorContains(t, err, "proof verification is disabled for all chains")
	})

	t.Run("check returns false is flag is disabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: false,
				},
			},
		})

		err := k.CheckBlockHeaderVerificationEnabled(ctx, chains.Ethereum.ChainId)
		require.ErrorContains(
			t,
			err,
			fmt.Sprintf("proof verification is disabled for chain %d", chains.Ethereum.ChainId),
		)
	})
}
