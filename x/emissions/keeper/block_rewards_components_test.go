package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	emissionskeeper "github.com/zeta-chain/zetacore/x/emissions/keeper"
)

func TestKeeper_CalculateFixedValidatorRewards(t *testing.T) {
	tt := []struct {
		name                 string
		blockTimeInSecs      string
		expectedBlockRewards sdk.Dec
	}{
		{
			name:                 "Block Time 5.7",
			blockTimeInSecs:      "5.7",
			expectedBlockRewards: sdk.MustNewDecFromStr("9620949074074074074.074070733466756687"),
		},
		{
			name:                 "Block Time 6",
			blockTimeInSecs:      "6",
			expectedBlockRewards: sdk.MustNewDecFromStr("10127314814814814814.814814814814814815"),
		},
		{
			name:                 "Block Time 3",
			blockTimeInSecs:      "3",
			expectedBlockRewards: sdk.MustNewDecFromStr("5063657407407407407.407407407407407407"),
		},
		{
			name:                 "Block Time 2",
			blockTimeInSecs:      "2",
			expectedBlockRewards: sdk.MustNewDecFromStr("3375771604938271604.938271604938271605"),
		},
		{
			name:                 "Block Time 8",
			blockTimeInSecs:      "8",
			expectedBlockRewards: sdk.MustNewDecFromStr("13503086419753086419.753086419753086420"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			blockRewards, err := emissionskeeper.CalculateFixedValidatorRewards(tc.blockTimeInSecs)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBlockRewards, blockRewards)
		})
	}
}
