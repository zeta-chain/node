package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
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
		require.Equal(t, chainState, res.ChainState[0])
	})

	t.Run("can run paginated queries", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		nbItems := 5

		items := make([]types.ChainState, nbItems)
		for i := range items {
			items[i] = sample.ChainState(int64(i))
			k.SetChainState(ctx, items[i])
		}

		request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllChainStateRequest {
			return &types.QueryAllChainStateRequest{
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
				resp, err := k.ChainStateAll(wctx, request(nil, uint64(i), uint64(step), false))
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.ChainState), step)
				require.Subset(t,
					nullify.Fill(items),
					nullify.Fill(resp.ChainState),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < nbItems; i += step {
				resp, err := k.ChainStateAll(wctx, request(next, 0, uint64(step), false))
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.ChainState), step)
				require.Subset(t,
					nullify.Fill(items),
					nullify.Fill(resp.ChainState),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			resp, err := k.ChainStateAll(wctx, request(nil, 0, 0, true))
			require.NoError(t, err)
			require.Equal(t, nbItems, int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(items),
				nullify.Fill(resp.ChainState),
			)
		})
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
