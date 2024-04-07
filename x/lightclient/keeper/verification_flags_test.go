package keeper_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
	"testing"
)

func TestKeeper_GetVerificationFlags(t *testing.T) {
	t.Run("can get and set verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		vf, found := k.GetVerificationFlags(ctx)
		require.False(t, found)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: true,
		})
		vf, found = k.GetVerificationFlags(ctx)
		require.True(t, found)
		require.Equal(t, types.VerificationFlags{
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: true,
		}, vf)
	})
}

func TestKeeper_CheckVerificationFlagsEnabled(t *testing.T) {
	t.Run("can check verification flags with ethereum enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: false,
		})

		err := k.CheckVerificationFlagsEnabled(ctx, chains.EthChain().ChainId)
		require.NoError(t, err)

		err = k.CheckVerificationFlagsEnabled(ctx, chains.BtcMainnetChain().ChainId)
		require.Error(t, err)
		require.ErrorContains(t, err, "proof verification not enabled for bitcoin")

		err = k.CheckVerificationFlagsEnabled(ctx, 1001)
		require.Error(t, err)
		require.ErrorContains(t, err, "doesn't support block header verification")
	})

	t.Run("can check verification flags with bitcoin enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: true,
		})

		err := k.CheckVerificationFlagsEnabled(ctx, chains.EthChain().ChainId)
		require.NoError(t, err)

		err = k.CheckVerificationFlagsEnabled(ctx, chains.BtcMainnetChain().ChainId)
		require.Error(t, err)
		require.ErrorContains(t, err, "proof verification not enabled for evm")

		err = k.CheckVerificationFlagsEnabled(ctx, 1001)
		require.Error(t, err)
		require.ErrorContains(t, err, "doesn't support block header verification")
	})
}
