package keeper

// Not needed moved to params
//func TestKeeper_SupportedChains(t *testing.T) {
//	keeper, ctx := SetupKeeper(t)
//	list := []*common.Chain{
//		{
//			ChainName: common.ChainName_eth_mainnet,
//			ChainId:   1,
//		},
//		{
//			ChainName: common.ChainName_btc_mainnet,
//			ChainId:   2,
//		},
//		{
//			ChainName: common.ChainName_bsc_mainnet,
//			ChainId:   3,
//		},
//	}
//
//	keeper.SetSupportedChain(ctx, types.SupportedChains{ChainList: list})
//	getList, found := keeper.GetSupportedChains(ctx)
//	assert.True(t, found)
//	assert.Equal(t, list, getList.ChainList)
//}
