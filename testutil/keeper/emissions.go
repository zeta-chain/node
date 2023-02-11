package keeper

//func EmissionsKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
//	storeKey := sdk.NewKVStoreKey(types.StoreKey)
//	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
//
//	db := tmdb.NewMemDB()
//	stateStore := store.NewCommitMultiStore(db)
//	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
//	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
//	require.NoError(t, stateStore.LoadLatestVersion())
//
//	registry := codectypes.NewInterfaceRegistry()
//	cdc := codec.NewProtoCodec(registry)
//
//	paramsSubspace := typesparams.NewSubspace(cdc,
//		types.Amino,
//		storeKey,
//		memStoreKey,
//		"EmissionsParams",
//	)
//	k := keeper.NewKeeper(
//		cdc,
//		storeKey,
//		memStoreKey,
//		paramsSubspace,
//	)
//
//	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
//
//	// Initialize params
//	k.SetParams(ctx, types.DefaultParams())
//
//	return k, ctx
//}
