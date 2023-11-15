package keeper_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_AddZetaAccounting(t *testing.T) {

	t.Run("should add aborted zeta amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(rand.Uint64())
		k.SetZetaAccounting(ctx, types.ZetaAccounting{
			originalAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedZetaAmount)
		addAmount := sdkmath.NewUint(rand.Uint64())
		k.AddZetaAccounting(ctx, addAmount)
		val, found = k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Add(addAmount), val.AbortedZetaAmount)
	})

	t.Run("cant find aborted amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		val, found := k.GetZetaAccounting(ctx)
		require.False(t, found)
		require.Equal(t, types.ZetaAccounting{}, val)
	})

	t.Run("add very high zeta amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		highAmount := sdkmath.NewUintFromString("100000000000000000000000000000000000000000000000")
		k.SetZetaAccounting(ctx, types.ZetaAccounting{
			highAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, highAmount, val.AbortedZetaAmount)
	})

}
