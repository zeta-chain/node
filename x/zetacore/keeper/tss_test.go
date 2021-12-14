package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func createNTSS(keeper *Keeper, ctx sdk.Context, n int) []types.TSS {
	items := make([]types.TSS, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetTSS(ctx, items[i])
	}
	return items
}

func TestTSSGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSS(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTSS(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestTSSRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSS(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTSS(ctx, item.Index)
		_, found := keeper.GetTSS(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestTSSGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSS(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTSS(ctx))
}
