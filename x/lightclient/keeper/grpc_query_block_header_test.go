package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
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

	t.Run("can run paginated queries", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		nbItems := 5

		items := make([]proofs.BlockHeader, nbItems)
		for i := range items {
			items[i] = sample.BlockHeader(sample.Hash().Bytes())
			k.SetBlockHeader(ctx, items[i])
		}

		request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllBlockHeaderRequest {
			return &types.QueryAllBlockHeaderRequest{
				Pagination: &query.PageRequest{
					Key:        next,
					Offset:     offset,
					Limit:      limit,
					CountTotal: total,
				},
			}
		}
		t.Run("ByOffset", func(t *testing.T) {
			step := 2
			for i := 0; i < nbItems; i += step {
				resp, err := k.BlockHeaderAll(wctx, request(nil, uint64(i), uint64(step), false))
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.BlockHeaders), step)
				require.Subset(t,
					nullify.Fill(items),
					nullify.Fill(resp.BlockHeaders),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < nbItems; i += step {
				resp, err := k.BlockHeaderAll(wctx, request(next, 0, uint64(step), false))
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.BlockHeaders), step)
				require.Subset(t,
					nullify.Fill(items),
					nullify.Fill(resp.BlockHeaders),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			resp, err := k.BlockHeaderAll(wctx, request(nil, 0, 0, true))
			require.NoError(t, err)
			require.Equal(t, nbItems, int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(items),
				nullify.Fill(resp.BlockHeaders),
			)
		})
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
