package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestKeeper_RateLimiterFlags(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.RateLimiterFlags(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if rate limiter flags not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.RateLimiterFlags(wctx, &types.QueryRateLimiterFlagsRequest{})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if rate limiter flags found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		flags := sample.RateLimiterFlags()
		k.SetRateLimiterFlags(ctx, flags)

		res, err := k.RateLimiterFlags(wctx, &types.QueryRateLimiterFlagsRequest{})

		require.NoError(t, err)
		require.Equal(t, &types.QueryRateLimiterFlagsResponse{
			RateLimiterFlags: flags,
		}, res)
	})
}
