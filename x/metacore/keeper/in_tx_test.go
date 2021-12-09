package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNInTx(keeper *Keeper, ctx sdk.Context, n int) []types.InTx {
	items := make([]types.InTx, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetInTx(ctx, items[i])
	}
	return items
}

func TestInTxGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNInTx(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetInTx(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestInTxRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNInTx(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveInTx(ctx, item.Index)
		_, found := keeper.GetInTx(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestInTxGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNInTx(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllInTx(ctx))
}
