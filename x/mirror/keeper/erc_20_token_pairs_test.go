package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/mirror/keeper"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

func createTestERC20TokenPairs(keeper *keeper.Keeper, ctx sdk.Context) types.ERC20TokenPairs {
	item := types.ERC20TokenPairs{}
	keeper.SetERC20TokenPairs(ctx, item)
	return item
}

func TestERC20TokenPairsGet(t *testing.T) {
	keeper, ctx := keepertest.MirrorKeeper(t)
	item := createTestERC20TokenPairs(keeper, ctx)
	rst, found := keeper.GetERC20TokenPairs(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestERC20TokenPairsRemove(t *testing.T) {
	keeper, ctx := keepertest.MirrorKeeper(t)
	createTestERC20TokenPairs(keeper, ctx)
	keeper.RemoveERC20TokenPairs(ctx)
	_, found := keeper.GetERC20TokenPairs(ctx)
	require.False(t, found)
}
