package app

//func Test_CreateUpgrades(t *testing.T) {
//
//	addErc20ModuleUpgrade := upgradeTrackerItem{
//		index: 1752528615,
//		storeUpgrade: &storetypes.StoreUpgrades{
//			Added: []string{erc20types.ModuleName},
//		},
//	}
//	tests := []struct {
//		name    string
//		chainID string
//		result  []upgradeTrackerItem
//		panic   bool
//	}{
//		{
//			name:    "mainnet chain ID",
//			chainID: "zetachain_7000-1",
//			result:  []upgradeTrackerItem{addErc20ModuleUpgrade},
//			panic:   false,
//		},
//		{
//			name:    "testnet chain ID",
//			chainID: "zetachain_7001-1",
//			result:  []upgradeTrackerItem{},
//			panic:   false,
//		},
//		{
//			name:    "invalid chain ID",
//			chainID: "zetachain-7001-1",
//			result:  []upgradeTrackerItem{},
//			panic:   true,
//		},
//		{
//			name:    "empty chain ID",
//			chainID: "",
//			result:  []upgradeTrackerItem{},
//			panic:   false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.panic {
//				require.Panics(t, func() {
//					createUpgrades(tt.chainID)
//				}, "createUpgrades(%s) should panic", tt.chainID)
//				return
//			} else {
//				require.NotPanics(t, func() {
//					createUpgrades(tt.chainID)
//				}, "createUpgrades(%s) should not panic", tt.chainID)
//			}
//
//			result := createUpgrades(tt.chainID)
//			require.Equal(t, tt.result, result, "createUpgrades(%s) = %v, want %v", tt.chainID, result, tt.result)
//		})
//	}
//}
