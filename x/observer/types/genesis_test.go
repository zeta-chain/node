package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestGenesisState_Validate(t *testing.T) {
	invalidChainParamsGen := types.DefaultGenesis()
	chainParams := types.GetDefaultChainParams().ChainParams
	invalidChainParamsGen.ChainParamsList.ChainParams = append(chainParams, chainParams[0])

	gsWithDuplicateNodeAccountList := types.DefaultGenesis()
	nodeAccount := sample.NodeAccount()
	gsWithDuplicateNodeAccountList.NodeAccountList = []*types.NodeAccount{nodeAccount, nodeAccount}

	gsWithDuplicateChainNonces := types.DefaultGenesis()
	chainNonce := sample.ChainNonces(0)
	gsWithDuplicateChainNonces.ChainNonces = []types.ChainNonces{chainNonce, chainNonce}

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
		{
			desc:     "invalid genesis state duplicate node account list",
			genState: gsWithDuplicateNodeAccountList,
			valid:    false,
		},
		{
			desc:     "invalid genesis state duplicate chain nonces",
			genState: gsWithDuplicateChainNonces,
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
}
