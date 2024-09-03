package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestKeeper_HeaderSupportedChains(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.HeaderSupportedChains(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return empty set if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, _ := k.HeaderSupportedChains(wctx, &types.QueryHeaderSupportedChainsRequest{})
		require.Len(t, res.HeaderSupportedChains, 0)
	})

	t.Run("should return if block header state is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		bhv := sample.BlockHeaderVerification()
		k.SetBlockHeaderVerification(ctx, bhv)

		res, err := k.HeaderSupportedChains(wctx, &types.QueryHeaderSupportedChainsRequest{})
		require.NoError(t, err)
		require.Equal(t, bhv.HeaderSupportedChains, res.HeaderSupportedChains)
	})
}

func TestKeeper_HeaderEnabledChains(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.HeaderEnabledChains(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return empty set if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, _ := k.HeaderEnabledChains(wctx, &types.QueryHeaderEnabledChainsRequest{})
		require.Len(t, res.HeaderEnabledChains, 0)
	})

	t.Run("should return if block header state is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		bhv := sample.BlockHeaderVerification()
		k.SetBlockHeaderVerification(ctx, bhv)

		res, err := k.HeaderEnabledChains(wctx, &types.QueryHeaderEnabledChainsRequest{})
		require.NoError(t, err)
		require.Equal(t, bhv.GetHeaderEnabledChains(), res.HeaderEnabledChains)
	})
}
