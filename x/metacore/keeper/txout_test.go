package keeper

import (
	"testing"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func createNTxout(keeper *Keeper, ctx sdk.Context, n int) []types.Txout {
	items := make([]types.Txout, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Id = keeper.AppendTxout(ctx, items[i])
	}
	return items
}

func TestTxoutGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxout(keeper, ctx, 10)
	for _, item := range items {
		assert.Equal(t, item, keeper.GetTxout(ctx, item.Id))
	}
}

func TestTxoutExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxout(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasTxout(ctx, item.Id))
	}
}

func TestTxoutRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxout(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTxout(ctx, item.Id)
		assert.False(t, keeper.HasTxout(ctx, item.Id))
	}
}

func TestTxoutGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxout(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTxout(ctx))
}

func TestTxoutCount(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxout(keeper, ctx, 10)
	count := uint64(len(items))
	assert.Equal(t, count, keeper.GetTxoutCount(ctx))
}
