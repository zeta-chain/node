package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func createTestTxList(keeper *Keeper, ctx sdk.Context) types.TxList {
	item := types.TxList{
		Creator: "any",
	}
	keeper.SetTxList(ctx, item)
	return item
}

func TestTxListGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	item := createTestTxList(keeper, ctx)
	rst, found := keeper.GetTxList(ctx)
	assert.True(t, found)
	assert.Equal(t, item, rst)
}
func TestTxListRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	createTestTxList(keeper, ctx)
	keeper.RemoveTxList(ctx)
	_, found := keeper.GetTxList(ctx)
	assert.False(t, found)
}
