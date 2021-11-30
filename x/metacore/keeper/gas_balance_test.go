package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNGasBalance(keeper *Keeper, ctx sdk.Context, n int) []types.GasBalance {
	items := make([]types.GasBalance, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetGasBalance(ctx, items[i])
	}
	return items
}

func TestGasBalanceGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasBalance(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGasBalance(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestGasBalanceRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasBalance(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGasBalance(ctx, item.Index)
		_, found := keeper.GetGasBalance(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestGasBalanceGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasBalance(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllGasBalance(ctx))
}
