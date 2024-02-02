package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGenesisState_Validate(t *testing.T) {
	invalidChainParamsGen := types.DefaultGenesis()
	chainParams := types.GetDefaultChainParams().ChainParams
	invalidChainParamsGen.ChainParamsList.ChainParams = append(chainParams, chainParams[0])

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
			desc:     "valid genesis state",
			genState: &types.GenesisState{},
			valid:    true,
		},
		{
			desc:     "invalid chain params",
			genState: invalidChainParamsGen,
			valid:    false,
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
