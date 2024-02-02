package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				ForeignCoinsList: []types.ForeignCoins{
					{
						Zrc20ContractAddress: "0",
					},
					{
						Zrc20ContractAddress: "1",
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicated foreignCoins",
			genState: &types.GenesisState{
				ForeignCoinsList: []types.ForeignCoins{
					{
						Zrc20ContractAddress: "0",
					},
					{
						Zrc20ContractAddress: "0",
					},
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
