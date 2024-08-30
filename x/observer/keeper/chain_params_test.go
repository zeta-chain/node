package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_GetSupportedChainFromChainID(t *testing.T) {
	t.Run("return nil if chain not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// no core params list
		_, found := k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		require.False(t, found)

		// core params list but chain not in list
		setSupportedChain(ctx, *k, getValidEthChainIDWithIndex(t, 0))
		_, found = k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 1))
		require.False(t, found)

		// chain params list but chain not supported
		chainParams := sample.ChainParams(getValidEthChainIDWithIndex(t, 0))
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{chainParams},
		})
		_, found = k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		require.False(t, found)
	})

	t.Run("return chain if chain found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		chainID := getValidEthChainIDWithIndex(t, 0)
		setSupportedChain(ctx, *k, getValidEthChainIDWithIndex(t, 1), chainID)
		chain, found := k.GetSupportedChainFromChainID(ctx, chainID)
		require.True(t, found)
		require.EqualValues(t, chainID, chain.ChainId)
	})
}

func TestKeeper_GetChainParamsByChainID(t *testing.T) {
	t.Run("return false if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		_, found := k.GetChainParamsByChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		require.False(t, found)
	})

	t.Run("return true if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		chainParams := sample.ChainParams(getValidEthChainIDWithIndex(t, 0))
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{chainParams},
		})
		res, found := k.GetChainParamsByChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		require.True(t, found)
		require.Equal(t, chainParams, res)
	})

	t.Run("return false if chain not found in params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		chainParams := sample.ChainParams(getValidEthChainIDWithIndex(t, 0))
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{chainParams},
		})
		_, found := k.GetChainParamsByChainID(ctx, getValidEthChainIDWithIndex(t, 1))
		require.False(t, found)
	})
}
func TestKeeper_GetSupportedChains(t *testing.T) {
	t.Run("return empty list if no core params list", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		require.Empty(t, k.GetSupportedChains(ctx))
	})

	t.Run("return list containing supported chains", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		chainList := chains.ExternalChainList([]chains.Chain{})

		require.Greater(t, len(chainList), 5)
		supported1 := chainList[0]
		supported2 := chainList[1]
		unsupported := chainList[2]
		supported3 := chainList[3]
		supported4 := chainList[4]

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
