package authority_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestGenesis(t *testing.T) {
	t.Run("valid genesis", func(t *testing.T) {
		genesisState := types.GenesisState{
			Policies:          sample.Policies(),
			AuthorizationList: sample.AuthorizationList("sample"),
			ChainInfo:         sample.ChainInfo(42),
		}

		// Init
		k, ctx := keepertest.AuthorityKeeper(t)
		authority.InitGenesis(ctx, *k, genesisState)

		// Check policy is set
		policies, found := k.GetPolicies(ctx)
		require.True(t, found)
		require.Equal(t, genesisState.Policies, policies)

		// Check authorization list is set
		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, genesisState.AuthorizationList, authorizationList)

		// Check chain info is set
		chainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Equal(t, genesisState.ChainInfo, chainInfo)

		// Export
		got := authority.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after init and export
		nullify.Fill(&genesisState)
		nullify.Fill(got)
		require.Equal(t, genesisState, *got)
	})
}
