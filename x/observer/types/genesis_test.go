package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGenesisState_Validate(t *testing.T) {
	invalidCoreParamsGen := types.DefaultGenesis()
	coreParams := types.GetCoreParams().CoreParams
	invalidCoreParamsGen.CoreParamsList.CoreParams = append(coreParams, coreParams[0])

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
			desc:     "invalid core params",
			genState: invalidCoreParamsGen,
			valid:    false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

	list := types.GetCoreParams()
	list.CoreParams = append(list.CoreParams, list.CoreParams[0])
}
