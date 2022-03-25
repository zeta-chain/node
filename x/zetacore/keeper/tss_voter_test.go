package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func createNTSSVoter(keeper *Keeper, ctx sdk.Context, n int) []types.TSSVoter {
	items := make([]types.TSSVoter, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetTSSVoter(ctx, items[i])
	}
	return items
}

func TestTSSVoterGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSSVoter(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTSSVoter(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestTSSVoterRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSSVoter(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTSSVoter(ctx, item.Index)
		_, found := keeper.GetTSSVoter(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestTSSVoterGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSSVoter(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTSSVoter(ctx))
}
