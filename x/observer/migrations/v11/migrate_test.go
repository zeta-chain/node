package v11_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	v11 "github.com/zeta-chain/node/x/observer/migrations/v11"
	"github.com/zeta-chain/node/x/observer/types"
)

// getTestChainParams returns a list of chain params for testing
func getTestChainParams() types.ChainParamsList {
	return types.ChainParamsList{
		ChainParams: []*types.ChainParams{
			sample.ChainParams(chains.Ethereum.ChainId),
			sample.ChainParams(chains.BscMainnet.ChainId),
			sample.ChainParams(chains.Amoy.ChainId),
			sample.ChainParams(chains.ArbitrumMainnet.ChainId),
		},
	}
}

func TestMigrateStore(t *testing.T) {
	t.Run("can migrate chain params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		testChainParams := getTestChainParams()
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// ensure we still have same number of chain params after migration
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, newChainParams.ChainParams, len(oldChainParams.ChainParams))

		for _, params := range newChainParams.ChainParams {
			// verify that chains that should have SkipBlockScan set to true
			if params.ChainId == chains.Amoy.ChainId || params.ChainId == chains.ArbitrumMainnet.ChainId {
				require.True(t, params.SkipBlockScan, "SkipBlockScan should be true for chain %d", params.ChainId)
			} else {
				require.False(t, params.SkipBlockScan, "SkipBlockScan should be false for chain %d", params.ChainId)
			}
		}
	})
}
