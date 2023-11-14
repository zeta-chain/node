package keeper_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_AddAbortedZetaAmount(t *testing.T) {

	t.Run("should add aborted zeta amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(rand.Uint64())
		k.SetAbortedZetaAmount(ctx, types.AbortedZetaAmount{
			originalAmount,
		})
		val, found := k.GetAbortedZetaAmount(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.Amount)
		addAmount := sdkmath.NewUint(rand.Uint64())
		k.AddAbortedZetaAmount(ctx, addAmount)
		val, found = k.GetAbortedZetaAmount(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Add(addAmount), val.Amount)
	})

	t.Run("cant find aborted amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		val, found := k.GetAbortedZetaAmount(ctx)
		require.False(t, found)
		require.Equal(t, types.AbortedZetaAmount{}, val)
	})

	t.Run("add very high zeta amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		highAmount := sdkmath.NewUintFromString("100000000000000000000000000000000000000000000000")
		k.SetAbortedZetaAmount(ctx, types.AbortedZetaAmount{
			highAmount,
		})
		val, found := k.GetAbortedZetaAmount(ctx)
		require.True(t, found)
		require.Equal(t, highAmount, val.Amount)
	})

}
