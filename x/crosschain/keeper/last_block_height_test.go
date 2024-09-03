package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func createNLastBlockHeight(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.LastBlockHeight {
	items := make([]types.LastBlockHeight, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetLastBlockHeight(ctx, items[i])
	}
	return items
}

func TestLastBlockHeightGet(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNLastBlockHeight(k, ctx, 10)
	for _, item := range items {
		rst, found := k.GetLastBlockHeight(ctx, item.Index)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestLastBlockHeightRemove(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNLastBlockHeight(k, ctx, 10)
	for _, item := range items {
		k.RemoveLastBlockHeight(ctx, item.Index)
		_, found := k.GetLastBlockHeight(ctx, item.Index)
		require.False(t, found)
	}
}

func TestLastBlockHeightGetAll(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNLastBlockHeight(k, ctx, 10)
	require.Equal(t, items, k.GetAllLastBlockHeight(ctx))
}
