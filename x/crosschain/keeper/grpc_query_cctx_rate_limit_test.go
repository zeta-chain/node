package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_CctxListPendingWithRateLimit(t *testing.T) {
	t.Run("should fail for empty req", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.CctxListPendingWithinRateLimit(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})
	t.Run("should fail if limit is too high", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.CctxListPendingWithinRateLimit(ctx, &types.QueryListCctxPendingWithRateLimitRequest{Limit: keeper.MaxPendingCctxs + 1})
		require.ErrorContains(t, err, "limit exceeds max limit of")
	})
	t.Run("should fail if no TSS", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.CctxListPendingWithinRateLimit(ctx, &types.QueryListCctxPendingWithRateLimitRequest{Limit: 1})
		require.ErrorContains(t, err, "tss not found")
	})
	t.Run("should return empty list if no nonces", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		//  set TSS
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		_, err := k.CctxListPendingWithinRateLimit(ctx, &types.QueryListCctxPendingWithRateLimitRequest{Limit: 1})
		require.ErrorContains(t, err, "pending nonces not found")
	})
}
