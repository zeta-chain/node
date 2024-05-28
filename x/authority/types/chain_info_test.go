package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestDefaultChainInfo(t *testing.T) {
	t.Run("default is empty", func(t *testing.T) {
		chainInfo := types.DefaultChainInfo()
		require.Empty(t, chainInfo.Chains)
	})
}

func TestChainInfo_Validate(t *testing.T) {
	tests := []struct {
		name        string
		chainInfo   types.ChainInfo
		errContains string
	}{
		{
			name:      "empty is valid",
			chainInfo: types.ChainInfo{},
		},
		{
			name:      "valid chain info",
			chainInfo: sample.ChainInfo(),
		},
		{
			name: "invalid if chain is invalid",
			chainInfo: types.ChainInfo{
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
			errContains: "chain ID must be positive",
		},
		{
			name: "invalid if chain is not external",
			chainInfo: types.ChainInfo{
				Chains: []chains.Chain{
					{
						ChainId:     42,
						ChainName:   chains.ChainName_empty,
						Network:     chains.Network_optimism,
						NetworkType: chains.NetworkType_testnet,
						Vm:          chains.Vm_evm,
						Consensus:   chains.Consensus_op_stack,
						IsExternal:  false,
					},
				},
			},
			errContains: "not external",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chainInfo.Validate()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
