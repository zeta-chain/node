package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNTxin(keeper *Keeper, ctx sdk.Context, n int) []types.Txin {
	items := make([]types.Txin, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetTxin(ctx, items[i])
	}
	return items
}

func TestTxinGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxin(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTxin(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestTxinRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxin(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTxin(ctx, item.Index)
		_, found := keeper.GetTxin(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestTxinGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxin(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTxin(ctx))
}
