package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func createNNodeAccount(keeper *Keeper, ctx sdk.Context, n int) []types.NodeAccount {
	items := make([]types.NodeAccount, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetNodeAccount(ctx, items[i])
	}
	return items
}

func TestNodeAccountGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetNodeAccount(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestNodeAccountRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveNodeAccount(ctx, item.Index)
		_, found := keeper.GetNodeAccount(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestNodeAccountGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllNodeAccount(ctx))
}
