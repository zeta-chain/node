package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"strconv"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNForeignCoins(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ForeignCoins {
	items := make([]types.ForeignCoins, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetForeignCoins(ctx, items[i])
	}
	return items
}

//func TestForeignCoinsGet(t *testing.T) {
//	keeper, ctx := keepertest.FungibleKeeper(t)
//	items := createNForeignCoins(keeper, ctx, 10)
//	for _, item := range items {
//		rst, found := keeper.GetForeignCoins(ctx,
//			item.Index,
//		)
//		require.True(t, found)
//		require.Equal(t,
//			nullify.Fill(&item),
//			nullify.Fill(&rst),
//		)
//	}
//}
//
//func TestForeignCoinsGetAll(t *testing.T) {
//	keeper, ctx := keepertest.FungibleKeeper(t)
//	items := createNForeignCoins(keeper, ctx, 10)
//	require.ElementsMatch(t,
//		nullify.Fill(items),
//		nullify.Fill(keeper.GetAllForeignCoins(ctx)),
//	)
//}
