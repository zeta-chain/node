package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
    
	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNReceive(keeper *Keeper, ctx sdk.Context, n int) []types.Receive {
	items := make([]types.Receive, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetReceive(ctx, items[i])
	}
	return items
}

func TestReceiveGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNReceive(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetReceive(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestReceiveRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNReceive(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveReceive(ctx, item.Index)
		_, found := keeper.GetReceive(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestReceiveGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNReceive(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllReceive(ctx))
}
