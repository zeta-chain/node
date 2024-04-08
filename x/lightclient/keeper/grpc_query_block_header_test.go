package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestKeeper_BlockHeaderAll(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlockHeaderAll(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		bh := sample.BlockHeader(sample.Hash().Bytes())
		k.SetBlockHeader(ctx, bh)

		res, err := k.BlockHeaderAll(wctx, &types.QueryAllBlockHeaderRequest{})
		require.NoError(t, err)
		require.Equal(t, bh, res.BlockHeaders[0])
	})
}

func TestKeeper_BlockHeader(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlockHeader(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlockHeader(wctx, &types.QueryGetBlockHeaderRequest{
			BlockHash: sample.Hash().Bytes(),
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		bh := sample.BlockHeader(sample.Hash().Bytes())
		k.SetBlockHeader(ctx, bh)

		res, err := k.BlockHeader(wctx, &types.QueryGetBlockHeaderRequest{
			BlockHash: bh.Hash,
		})
		require.NoError(t, err)
		require.Equal(t, &bh, res.BlockHeader)
	})
}
