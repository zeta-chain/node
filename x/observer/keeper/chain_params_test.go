package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetSupportedChainFromChainID(t *testing.T) {
	t.Run("return nil if chain not found", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)

		// no core params list
		require.Nil(t, k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0)))

		// core params list but chain not in list
		setSupportedChain(ctx, *k, getValidEthChainIDWithIndex(t, 0))
		require.Nil(t, k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 1)))

		// chain params list but chain not supported
		chainParams := sample.ChainParams(getValidEthChainIDWithIndex(t, 0))
		chainParams.IsActive = false
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{chainParams},
		})
		require.Nil(t, k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0)))
	})

	t.Run("return chain if chain found", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		chainID := getValidEthChainIDWithIndex(t, 0)
		setSupportedChain(ctx, *k, getValidEthChainIDWithIndex(t, 1), chainID)
		chain := k.GetSupportedChainFromChainID(ctx, chainID)
		require.NotNil(t, chain)
		require.EqualValues(t, chainID, chain.ChainId)
	})
}

func TestKeeper_GetSupportedChains(t *testing.T) {
	t.Run("return empty list if no core params list", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		require.Empty(t, k.GetSupportedChains(ctx))
	})

	t.Run("return list containing supported chains", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)

		require.Greater(t, len(common.ExternalChainList()), 5)
		supported1 := common.ExternalChainList()[0]
		supported2 := common.ExternalChainList()[1]
		unsupported := common.ExternalChainList()[2]
		supported3 := common.ExternalChainList()[3]
		supported4 := common.ExternalChainList()[4]

		var chainParamsList []*types.ChainParams
		chainParamsList = append(chainParamsList, sample.ChainParamsSupported(supported1.ChainId))
		chainParamsList = append(chainParamsList, sample.ChainParamsSupported(supported2.ChainId))
		chainParamsList = append(chainParamsList, sample.ChainParams(unsupported.ChainId))
		chainParamsList = append(chainParamsList, sample.ChainParamsSupported(supported3.ChainId))
		chainParamsList = append(chainParamsList, sample.ChainParamsSupported(supported4.ChainId))

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: chainParamsList,
		})

		supportedChains := k.GetSupportedChains(ctx)

		require.Len(t, supportedChains, 4)
		require.EqualValues(t, supported1.ChainId, supportedChains[0].ChainId)
		require.EqualValues(t, supported2.ChainId, supportedChains[1].ChainId)
		require.EqualValues(t, supported3.ChainId, supportedChains[2].ChainId)
		require.EqualValues(t, supported4.ChainId, supportedChains[3].ChainId)
	})
}
