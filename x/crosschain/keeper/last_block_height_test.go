package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createNLastBlockHeight(keeper *Keeper, ctx sdk.Context, n int) []types.LastBlockHeight {
	items := make([]types.LastBlockHeight, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetLastBlockHeight(ctx, items[i])
	}
	return items
}

func TestLastBlockHeightGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNLastBlockHeight(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetLastBlockHeight(ctx, item.Index)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestLastBlockHeightRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNLastBlockHeight(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveLastBlockHeight(ctx, item.Index)
		_, found := keeper.GetLastBlockHeight(ctx, item.Index)
		require.False(t, found)
	}
}

func TestLastBlockHeightGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNLastBlockHeight(keeper, ctx, 10)
	require.Equal(t, items, keeper.GetAllLastBlockHeight(ctx))
}
