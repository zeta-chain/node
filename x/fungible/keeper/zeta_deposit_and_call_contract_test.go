package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
    "github.com/zeta-chain/zetacore/testutil/nullify"
)

func createTestZetaDepositAndCallContract(keeper *keeper.Keeper, ctx sdk.Context) types.ZetaDepositAndCallContract {
	item := types.ZetaDepositAndCallContract{}
	keeper.SetZetaDepositAndCallContract(ctx, item)
	return item
}

func TestZetaDepositAndCallContractGet(t *testing.T) {
	keeper, ctx := keepertest.FungibleKeeper(t)
	item := createTestZetaDepositAndCallContract(keeper, ctx)
	rst, found := keeper.GetZetaDepositAndCallContract(ctx)
    require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestZetaDepositAndCallContractRemove(t *testing.T) {
	keeper, ctx := keepertest.FungibleKeeper(t)
	createTestZetaDepositAndCallContract(keeper, ctx)
	keeper.RemoveZetaDepositAndCallContract(ctx)
    _, found := keeper.GetZetaDepositAndCallContract(ctx)
    require.False(t, found)
}
