package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNTxoutConfirmation(keeper *Keeper, ctx sdk.Context, n int) []types.TxoutConfirmation {
	items := make([]types.TxoutConfirmation, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetTxoutConfirmation(ctx, items[i])
	}
	return items
}

func TestTxoutConfirmationGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxoutConfirmation(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTxoutConfirmation(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestTxoutConfirmationRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxoutConfirmation(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTxoutConfirmation(ctx, item.Index)
		_, found := keeper.GetTxoutConfirmation(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestTxoutConfirmationGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxoutConfirmation(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTxoutConfirmation(ctx))
}
