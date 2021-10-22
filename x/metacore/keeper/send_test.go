package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNSend(keeper *Keeper, ctx sdk.Context, n int) []types.Send {
	items := make([]types.Send, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetSend(ctx, items[i])
	}
	return items
}

func TestSendGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetSend(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestSendRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveSend(ctx, item.Index)
		_, found := keeper.GetSend(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestSendGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllSend(ctx))
}
