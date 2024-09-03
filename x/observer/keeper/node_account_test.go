package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// Keeper Tests

func createNNodeAccount(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.NodeAccount {
	items := make([]types.NodeAccount, n)
	for i := range items {
		items[i].Operator = fmt.Sprintf("%d", i)
		keeper.SetNodeAccount(ctx, items[i])
	}
	return items
}

func TestNodeAccountGet(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	items := createNNodeAccount(k, ctx, 10)
	for _, item := range items {
		rst, found := k.GetNodeAccount(ctx, item.Operator)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestNodeAccountRemove(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	items := createNNodeAccount(k, ctx, 10)
	for _, item := range items {
		k.RemoveNodeAccount(ctx, item.Operator)
		_, found := k.GetNodeAccount(ctx, item.Operator)
		require.False(t, found)
	}
}

func TestNodeAccountGetAll(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	items := createNNodeAccount(k, ctx, 10)
	require.Equal(t, items, k.GetAllNodeAccount(ctx))
}
