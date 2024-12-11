package v3_test

import (
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/migrations/v3"
	"github.com/zeta-chain/node/x/fungible/types"
	"testing"
)

func TestMigrateStore(t *testing.T) {
	tests := []struct {
		name         string
		assetList    []string
		expectedList []string
	}{
		{
			name:         "no asset to update",
			assetList:    []string{},
			expectedList: []string{},
		},
		{
			name: "assets to update",
			assetList: []string{
				"",
				"0x5a4f260a7d716c859a2736151cb38b9c58c32c64", // lowercase
				"",
				"0xc0ffee254729296a45a3885639AC7E10F9d54979", // checksum
				"",
				"",
				"Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr",
				"BrS9iNMC3y8J4QTmCz8VrGrYepdoxXYvKxcDMiixwLn5",
				"0x999999CF1046E68E36E1AA2E0E07105EDDD1F08E", // uppcase
			},
			expectedList: []string{
				"",
				"0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
				"",
				"0xc0ffee254729296a45a3885639AC7E10F9d54979",
				"",
				"",
				"Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr",
				"BrS9iNMC3y8J4QTmCz8VrGrYepdoxXYvKxcDMiixwLn5",
				"0x999999cf1046e68e36E1aA2E0E07105eDDD1f08E",
			},
		},
	}

	for _, tt := range tests {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		// Arrange

		// set sample foreign coins
		expectedForeignCoins := make([]types.ForeignCoins, len(tt.assetList))
		for i, asset := range tt.assetList {
			expectedForeignCoins[i] = sample.ForeignCoins(t, sample.EthAddress().Hex())
			expectedForeignCoins[i].Asset = asset
			k.SetForeignCoins(ctx, expectedForeignCoins[i])
		}

		// update for expected list
		for i := range tt.assetList {
			expectedForeignCoins[i].Asset = tt.expectedList[i]
		}

		// Act
		err := v3.MigrateStore(ctx, k)
		require.NoError(t, err)

		// Assert
		actualForeignCoins := k.GetAllForeignCoins(ctx)
		require.ElementsMatch(t, expectedForeignCoins, actualForeignCoins)
	}

}
