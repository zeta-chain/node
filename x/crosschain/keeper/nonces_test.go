package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// Keeper Tests
func createNChainNonces(keeper *Keeper, ctx sdk.Context, n int) []types.ChainNonces {
	items := make([]types.ChainNonces, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetChainNonces(ctx, items[i])
	}
	return items
}

func TestChainNoncesGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNChainNonces(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetChainNonces(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestChainNoncesRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNChainNonces(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveChainNonces(ctx, item.Index)
		_, found := keeper.GetChainNonces(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestChainNoncesGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNChainNonces(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllChainNonces(ctx))
}
