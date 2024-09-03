package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/types"
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
				Policies:          sample.Policies(),
				ChainInfo:         sample.ChainInfo(42),
				AuthorizationList: sample.AuthorizationList("testMessage"),
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
				ChainInfo:         sample.ChainInfo(42),
				AuthorizationList: sample.AuthorizationList("testMessage"),
			},
			errContains: "invalid address",
		},
		{
			name: "invalid if chainInfo is invalid",
			gs: &types.GenesisState{
				Policies: sample.Policies(),
				ChainInfo: types.ChainInfo{
					Chains: []chains.Chain{
						{
							ChainId:     0,
							Network:     chains.Network_optimism,
							NetworkType: chains.NetworkType_testnet,
							Vm:          chains.Vm_evm,
							Consensus:   chains.Consensus_op_stack,
							IsExternal:  true,
							Name:        "foo",
						},
					},
				},
				AuthorizationList: sample.AuthorizationList("testMessage"),
			},
			errContains: "chain ID must be positive",
		},
		{
			name: "invalid if authorization list is invalid",
			gs: &types.GenesisState{
				Policies:  sample.Policies(),
				ChainInfo: sample.ChainInfo(42),
				AuthorizationList: types.AuthorizationList{
					Authorizations: []types.Authorization{
						{
							MsgUrl:           fmt.Sprintf("/zetachain/%d%s", 0, "testMessage"),
							AuthorizedPolicy: types.PolicyType_groupEmergency,
						},
						{
							MsgUrl:           fmt.Sprintf("/zetachain/%d%s", 0, "testMessage"),
							AuthorizedPolicy: types.PolicyType_groupAdmin,
						},
						{
							MsgUrl:           fmt.Sprintf("/zetachain/%d%s", 0, "testMessage"),
							AuthorizedPolicy: types.PolicyType_groupOperational,
						},
					},
				},
			},
			errContains: "duplicate message url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gs.Validate()
			if tt.errContains != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
