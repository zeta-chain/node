package v8_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	v8 "github.com/zeta-chain/zetacore/x/observer/migrations/v8"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func Test_MigrateStore(t *testing.T) {
	tt := []struct {
		name      string
		params    types.Params
		chainList []*chains.Chain
	}{
		{
			name:      "privnet params",
			params:    LegacyParamsForNetwork(chains.NetworkType_PRIVNET),
			chainList: chains.ChainListByNetworkType(chains.NetworkType_PRIVNET),
		},
		{
			name:      "testnet params",
			params:    LegacyParamsForNetwork(chains.NetworkType_TESTNET),
			chainList: chains.ChainListByNetworkType(chains.NetworkType_TESTNET),
		},
		{
			name:      "mainnet params",
			params:    LegacyParamsForNetwork(chains.NetworkType_MAINNET),
			chainList: chains.ChainListByNetworkType(chains.NetworkType_MAINNET),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, _, _ := keepertest.ObserverKeeper(t)
			k.SetParams(ctx, tc.params)

			chainListOld := make([]chains.Chain, len(tc.params.ObserverParams))
			for i, ob := range tc.params.ObserverParams {
				chainListOld[i] = *ob.Chain
			}

			err := v8.MigrateStore(ctx, k)
			require.NoError(t, err)

			params := k.GetParamsIfExists(ctx)
			chainListNew := make([]chains.Chain, len(params.ObserverParams))
			for i, ob := range params.ObserverParams {
				chainListNew[i] = *ob.Chain
			}
			chainListTest := make([]chains.Chain, len(tc.chainList))
			for i, chain := range tc.chainList {
				chainListTest[i] = *chain
			}
			require.NotEqual(t, chainListTest, chainListOld)
			require.Equal(t, chainListTest, chainListNew)
		})

	}

}

func LegacyParamsForNetwork(networkType chains.NetworkType) types.Params {
	chainList := chains.ChainListByNetworkType(networkType)
	observerParams := make([]*types.ObserverParams, len(chainList))
	for i, chain := range chainList {
		observerParams[i] = &types.ObserverParams{
			IsSupported: true,
			// Set chain-Id and chain-name only for legacy params
			Chain:                 &chains.Chain{ChainId: chain.ChainId, ChainName: chain.ChainName},
			BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
			MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000"),
		}
	}
	return types.NewParams(observerParams, types.DefaultAdminPolicy(), 100)
}
