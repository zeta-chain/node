package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	emissionskeeper "github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestKeeper_CalculateFixedValidatorRewards(t *testing.T) {
	tt := []struct {
		name            string
		blockTimeInSecs string
	}{
		{
			name:            "Block Time 5.7",
			blockTimeInSecs: "5.7",
		},
		{
			name:            "Block Time 6",
			blockTimeInSecs: "6",
		},
		{
			name:            "Block Time 3",
			blockTimeInSecs: "3",
		},
		{
			name:            "Block Time 2",
			blockTimeInSecs: "2",
		},
		{
			name:            "Block Time 8",
			blockTimeInSecs: "8",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			blockRewards, err := emissionskeeper.CalculateFixedValidatorRewards(tc.blockTimeInSecs)
			assert.NoError(t, err)
			avgBlockTime, err := sdk.NewDecFromStr(tc.blockTimeInSecs)
			assert.NoError(t, err)
			azetaAmountTotalRewards, err := sdk.NewDecFromStr("210000000000000000000000000")
			numberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
			numberOfBlocksTotal := numberOfBlocksInAMonth.Mul(sdk.NewDec(12)).Mul(sdk.NewDec(types.EmissionScheduledYears))
			constantRewardPerBlock := azetaAmountTotalRewards.Quo(numberOfBlocksTotal)
			assert.Equal(t, constantRewardPerBlock, blockRewards)
		})
	}
}
