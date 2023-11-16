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
	t.Run("should add aborted gas and erc20 amounts", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(rand.Uint64())
		k.SetZetaAccounting(ctx, types.ZetaAccounting{
			AbortedZetaAmount:  originalAmount,
			AbortedGasAmount:   originalAmount,
			AbortedErc20Amount: originalAmount,
		})
		val, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedZetaAmount)
		addAmount := sdkmath.NewUint(rand.Uint64())
		k.AddZetaAbortedAmount(ctx, addAmount)
		k.AddErc20AbortedAmount(ctx, addAmount)
		k.AddGasAbortedAmount(ctx, addAmount)
		val, found = k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Add(addAmount), val.AbortedZetaAmount)
		require.Equal(t, originalAmount.Add(addAmount), val.AbortedErc20Amount)
		require.Equal(t, originalAmount.Add(addAmount), val.AbortedGasAmount)
	})

}
