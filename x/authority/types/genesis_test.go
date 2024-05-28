package types_test

import (
	"github.com/zeta-chain/zetacore/pkg/chains"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestGenesisState_Validate(t *testing.T) {
	setConfig(t)

	tests := []struct {
		name        string
		gs          *types.GenesisState
		errContains string
	}{
		{
			name:        "default is valid",
			gs:          types.DefaultGenesis(),
			errContains: "",
		},
		{
			name: "valid genesis",
			gs: &types.GenesisState{
				Policies:  sample.Policies(),
				ChainInfo: sample.ChainInfo(42),
			},
			errContains: "",
		},
		{
			name: "invalid if policies is invalid",
			gs: &types.GenesisState{
				Policies: types.Policies{
					Items: []*types.Policy{
						{
							Address:    "invalid",
							PolicyType: types.PolicyType_groupEmergency,
						},
					},
				},
				ChainInfo: sample.ChainInfo(42),
			},
			errContains: "invalid address",
		},
		{
			name: "invalid if policies is invalid",
			gs: &types.GenesisState{
				Policies: sample.Policies(),
				ChainInfo: types.ChainInfo{
					Chains: []chains.Chain{
						{
							ChainId:     0,
							ChainName:   chains.ChainName_empty,
							Network:     chains.Network_optimism,
							NetworkType: chains.NetworkType_testnet,
							Vm:          chains.Vm_evm,
							Consensus:   chains.Consensus_op_stack,
							IsExternal:  true,
						},
					},
				},
			},
			errContains: "chain ID must be positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gs.Validate()
			if tt.errContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
