package fungible_test

//
//func TestGenesis(t *testing.T) {
//	genesisState := types.GenesisState{
//		Params: types.DefaultParams(),
//
//		ForeignCoinsList: []types.ForeignCoins{
//			{
//				Index: "0",
//			},
//			{
//				Index: "1",
//			},
//		},
//		SystemContract: &types.SystemContract{
//			SystemContract: "29",
//		},
//	}
//
//	k, ctx := keepertest.FungibleKeeper(t)
//	fungible.InitGenesis(ctx, *k, genesisState, authkeeper.AccountKeeper{})
//	got := fungible.ExportGenesis(ctx, *k)
//	require.NotNil(t, got)
//
//	nullify.Fill(&genesisState)
//	nullify.Fill(got)
//
//	require.ElementsMatch(t, genesisState.ForeignCoinsList, got.ForeignCoinsList)
//}
