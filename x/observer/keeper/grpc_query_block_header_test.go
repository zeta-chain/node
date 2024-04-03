package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetAllBlockHeaders(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetAllBlockHeaders(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		bh := proofs.BlockHeader{
			Height:     1,
			Hash:       sample.Hash().Bytes(),
			ParentHash: sample.Hash().Bytes(),
			ChainId:    1,
			Header:     proofs.HeaderData{},
		}
		k.SetBlockHeader(ctx, bh)

		res, err := k.GetAllBlockHeaders(wctx, &types.QueryAllBlockHeaderRequest{})
		require.NoError(t, err)
		require.Equal(t, &bh, res.BlockHeaders[0])
	})
}

func TestKeeper_GetBlockHeaderByHash(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetBlockHeaderByHash(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetBlockHeaderByHash(wctx, &types.QueryGetBlockHeaderByHashRequest{
			BlockHash: sample.Hash().Bytes(),
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		bh := proofs.BlockHeader{
			Height:     1,
			Hash:       sample.Hash().Bytes(),
			ParentHash: sample.Hash().Bytes(),
			ChainId:    1,
			Header:     proofs.HeaderData{},
		}
		k.SetBlockHeader(ctx, bh)

		res, err := k.GetBlockHeaderByHash(wctx, &types.QueryGetBlockHeaderByHashRequest{
			BlockHash: bh.Hash,
		})
		require.NoError(t, err)
		require.Equal(t, &bh, res.BlockHeader)
	})
}

func TestKeeper_GetBlockHeaderStateByChain(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetBlockHeaderStateByChain(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetBlockHeaderStateByChain(wctx, &types.QueryGetBlockHeaderStateRequest{
			ChainId: 1,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if block header state is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		bhs := types.BlockHeaderState{
			ChainId: 1,
		}
		k.SetBlockHeaderState(ctx, bhs)

		res, err := k.GetBlockHeaderStateByChain(wctx, &types.QueryGetBlockHeaderStateRequest{
			ChainId: 1,
		})
		require.NoError(t, err)
		require.Equal(t, &bhs, res.BlockHeaderState)
	})
}
