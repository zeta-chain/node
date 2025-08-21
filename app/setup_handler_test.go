package app

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func Test_CreateUpgrades(t *testing.T) {

	addErc20ModuleUpgrade := upgradeTrackerItem{
		index: 1752528615,
		storeUpgrade: &storetypes.StoreUpgrades{
			Added: []string{erc20types.ModuleName},
		},
	}
	tests := []struct {
		name    string
		chainID string
		result  []upgradeTrackerItem
	}{
		{
			name:    "mainnet chain ID",
			chainID: "zetachain_7000-1",
			result:  []upgradeTrackerItem{addErc20ModuleUpgrade},
		},
		{
			name:    "testnet chain ID",
			chainID: "zetachain_7001-1",
			result:  []upgradeTrackerItem{},
		},
		{
			name:    "invalid chain ID",
			chainID: "zetachain-7001-1",
			result:  []upgradeTrackerItem{},
		},
		{
			name:    "empty chain ID",
			chainID: "",
			result:  []upgradeTrackerItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createUpgrades(tt.chainID)
			require.Equal(t, tt.result, result, "createUpgrades(%s) = %v, want %v", tt.chainID, result, tt.result)
		})
	}
}
