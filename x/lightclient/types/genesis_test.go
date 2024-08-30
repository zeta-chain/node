package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestGenesisState_Validate(t *testing.T) {
	duplicatedHash := sample.Hash().Bytes()

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				BlockHeaderVerification: sample.BlockHeaderVerification(),
				BlockHeaders: []proofs.BlockHeader{
					sample.BlockHeader(sample.Hash().Bytes()),
					sample.BlockHeader(sample.Hash().Bytes()),
					sample.BlockHeader(sample.Hash().Bytes()),
				},
				ChainStates: []types.ChainState{
					sample.ChainState(chains.Ethereum.ChainId),
					sample.ChainState(chains.BitcoinMainnet.ChainId),
					sample.ChainState(chains.BscMainnet.ChainId),
				},
			},
			valid: true,
		},
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "duplicate block headers is invalid",
			genState: &types.GenesisState{
				BlockHeaders: []proofs.BlockHeader{
					sample.BlockHeader(sample.Hash().Bytes()),
					sample.BlockHeader(duplicatedHash),
					sample.BlockHeader(duplicatedHash),
				},
			},
			valid: false,
		},
		{
			desc: "duplicate chain state is invalid",
			genState: &types.GenesisState{
				ChainStates: []types.ChainState{
					sample.ChainState(chains.Ethereum.ChainId),
					sample.ChainState(chains.Ethereum.ChainId),
					sample.ChainState(chains.BscMainnet.ChainId),
				},
			},
			valid: false,
		},
		{
			desc: "invalid block header verification",
			genState: &types.GenesisState{
				BlockHeaderVerification: types.BlockHeaderVerification{
					HeaderSupportedChains: []types.HeaderSupportedChain{{ChainId: 1, Enabled: true}, {ChainId: 1, Enabled: true}},
				},
			},
			valid: false,
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
