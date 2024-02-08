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
			AbortedZetaAmount: originalAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedZetaAmount)
		addAmount := sdkmath.NewUint(rand.Uint64())
		k.AddZetaAbortedAmount(ctx, addAmount)
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
			AbortedZetaAmount: highAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, highAmount, val.AbortedZetaAmount)
	})

}

func TestKeeper_RemoveZetaAbortedAmount(t *testing.T) {
	t.Run("should remove aborted zeta amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUintFromString("100000000000000000000000000000000000000000000000")
		k.SetZetaAccounting(ctx, types.ZetaAccounting{
			AbortedZetaAmount: originalAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedZetaAmount)
		removeAmount := originalAmount.Sub(sdkmath.NewUintFromString("10000000000000000000000000000000000000000000000"))
		err := k.RemoveZetaAbortedAmount(ctx, removeAmount)
		require.NoError(t, err)
		val, found = k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Sub(removeAmount), val.AbortedZetaAmount)
	})
	t.Run("fail remove aborted zeta amount if accounting not set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		err := k.RemoveZetaAbortedAmount(ctx, sdkmath.OneUint())
		require.ErrorIs(t, err, types.ErrUnableToFindZetaAccounting)
	})
	t.Run("fail remove aborted zeta amount if insufficient amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(100)
		k.SetZetaAccounting(ctx, types.ZetaAccounting{
			AbortedZetaAmount: originalAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedZetaAmount)
		removeAmount := originalAmount.Add(sdkmath.NewUint(500))
		err := k.RemoveZetaAbortedAmount(ctx, removeAmount)
		require.ErrorIs(t, err, types.ErrInsufficientZetaAmount)
		val, found = k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedZetaAmount)
	})
}
