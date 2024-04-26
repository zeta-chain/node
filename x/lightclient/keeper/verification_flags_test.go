package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestKeeper_GetVerificationFlags(t *testing.T) {
	t.Run("can get and set verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		ethChainId := chains.EthChain.ChainId

		vf, found := k.GetVerificationFlags(ctx, ethChainId)
		require.False(t, found)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: ethChainId,
			Enabled: true,
		})
		vf, found = k.GetVerificationFlags(ctx, ethChainId)
		require.True(t, found)
		require.True(t, vf.Enabled)
	})
}

func TestKeeper_CheckVerificationFlagsEnabled(t *testing.T) {
	t.Run("can check verification flags with ethereum enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chains.EthChain.ChainId,
			Enabled: true,
		})

		err := k.CheckVerificationFlagsEnabled(ctx, chains.EthChain.ChainId)
		require.NoError(t, err)

		err = k.CheckVerificationFlagsEnabled(ctx, chains.BtcMainnetChain.ChainId)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("proof verification not enabled for,chain id: %d", chains.BtcMainnetChain.ChainId))

		err = k.CheckVerificationFlagsEnabled(ctx, 1000)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("proof verification not enabled for,chain id: %d", 1000))
	})

	t.Run("can check verification flags with bitcoin enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chains.BtcMainnetChain.ChainId,
			Enabled: true,
		})

		err := k.CheckVerificationFlagsEnabled(ctx, chains.EthChain.ChainId)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("proof verification not enabled for,chain id: %d", chains.EthChain.ChainId))

		err = k.CheckVerificationFlagsEnabled(ctx, chains.BtcMainnetChain.ChainId)
		require.NoError(t, err)

		err = k.CheckVerificationFlagsEnabled(ctx, 1000)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("proof verification not enabled for,chain id: %d", 1000))
	})
}
