package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNTxinVoter(keeper *Keeper, ctx sdk.Context, n int) []types.TxinVoter {
	items := make([]types.TxinVoter, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetTxinVoter(ctx, items[i])
	}
	return items
}

func TestTxinVoterGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxinVoter(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTxinVoter(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestTxinVoterRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxinVoter(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTxinVoter(ctx, item.Index)
		_, found := keeper.GetTxinVoter(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestTxinVoterGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTxinVoter(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTxinVoter(ctx))
}
