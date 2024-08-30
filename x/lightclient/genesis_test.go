package lightclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestGenesis(t *testing.T) {
	t.Run("can import and export genesis", func(t *testing.T) {
		genesisState := types.GenesisState{
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
		}

		// Init and export
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		lightclient.InitGenesis(ctx, *k, genesisState)
		got := lightclient.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after init and export
		nullify.Fill(&genesisState)
		nullify.Fill(got)
		require.Equal(t, genesisState, *got)
	})

	t.Run("can export genesis with empty state", func(t *testing.T) {
		// Export genesis with empty state
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		got := lightclient.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after export
		expected := types.GenesisState{
			BlockHeaderVerification: types.DefaultBlockHeaderVerification(),
			BlockHeaders:            []proofs.BlockHeader(nil),
			ChainStates:             []types.ChainState(nil),
		}
		require.Equal(t, expected, *got)
		require.Equal(
			t,
			expected.BlockHeaderVerification.HeaderSupportedChains,
			got.BlockHeaderVerification.HeaderSupportedChains,
		)
	})
}
