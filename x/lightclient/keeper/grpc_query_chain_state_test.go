package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestKeeper_ChainStateAll(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.ChainStateAll(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		chainState := sample.ChainState(42)
		k.SetChainState(ctx, chainState)

		res, err := k.ChainStateAll(wctx, &types.QueryAllChainStateRequest{})
		require.NoError(t, err)
		require.Equal(t, &chainState, res.ChainState[0])
	})
}

func TestKeeper_ChainState(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.ChainState(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.ChainState(wctx, &types.QueryGetChainStateRequest{
			ChainId: 1,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header state is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		chainState := sample.ChainState(42)
		k.SetChainState(ctx, chainState)

		res, err := k.ChainState(wctx, &types.QueryGetChainStateRequest{
			ChainId: 42,
		})
		require.NoError(t, err)
		require.Equal(t, &chainState, res.ChainState)
	})
}
