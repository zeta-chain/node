package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func createNSendVoter(keeper *Keeper, ctx sdk.Context, n int) []types.SendVoter {
	items := make([]types.SendVoter, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetSendVoter(ctx, items[i])
	}
	return items
}

func TestSendVoterGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSendVoter(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetSendVoter(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestSendVoterRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSendVoter(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveSendVoter(ctx, item.Index)
		_, found := keeper.GetSendVoter(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestSendVoterGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSendVoter(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllSendVoter(ctx))
}
