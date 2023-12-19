package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetObserver(t *testing.T) {

	tt := []struct {
		name             string
		mapper           []*types.ObserverMapper
		assertChain      *common.Chain
		assertObsListLen int
		isFound          bool
	}{
		{
			name: "4 eth Observers",
			mapper: types.CreateObserverMapperList(1, common.Chain{
				ChainName: common.ChainName_eth_mainnet,
				ChainId:   1,
			}),
			assertChain: &common.Chain{
				ChainName: common.ChainName_eth_mainnet,
				ChainId:   1,
			},
			assertObsListLen: 4,
			isFound:          true,
		},
		{
			name: "Filter out from multiple mappers",
			mapper: append(append(types.CreateObserverMapperList(1, common.Chain{
				ChainName: common.ChainName_eth_mainnet,
				ChainId:   1,
			}),
				types.CreateObserverMapperList(1, common.Chain{
					ChainName: common.ChainName_eth_mainnet,
					ChainId:   1,
				})...),
				types.CreateObserverMapperList(1, common.Chain{
					ChainName: common.ChainName_bsc_mainnet,
					ChainId:   2,
				})...),
			assertChain: &common.Chain{
				ChainName: common.ChainName_eth_mainnet,
				ChainId:   1,
			},
			assertObsListLen: 4,
			isFound:          true,
		},
		{
			name: "No Observers of expected Observation Chain",
			mapper: append(append(types.CreateObserverMapperList(1, common.Chain{
				ChainName: common.ChainName_btc_mainnet,
				ChainId:   3,
			}),
				types.CreateObserverMapperList(1, common.Chain{
					ChainName: common.ChainName_polygon_mainnet,
					ChainId:   4,
				})...),
				types.CreateObserverMapperList(1, common.Chain{
					ChainName: common.ChainName_bsc_mainnet,
					ChainId:   5,
				})...),
			assertChain: &common.Chain{
				ChainName: common.ChainName_eth_mainnet,
				ChainId:   1,
			},
			assertObsListLen: 0,
			isFound:          false,
		},
		{
			name: "No Observers of expected Observation Type",
			mapper: append(append(types.CreateObserverMapperList(1, common.Chain{
				ChainName: common.ChainName_btc_mainnet,
				ChainId:   3,
			}),
				types.CreateObserverMapperList(1, common.Chain{
					ChainName: common.ChainName_polygon_mainnet,
					ChainId:   4,
				})...),
				types.CreateObserverMapperList(1, common.Chain{
					ChainName: common.ChainName_bsc_mainnet,
					ChainId:   5,
				})...),
			assertChain: &common.Chain{
				ChainName: common.ChainName_eth_mainnet,
				ChainId:   1,
			},
			assertObsListLen: 0,
			isFound:          false,
		},
	}

	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			k, ctx := keepertest.ObserverKeeper(t)
			for _, mapper := range test.mapper {
				k.SetObserverMapper(ctx, mapper)
			}
			mapper, found := k.GetObserverMapper(ctx, test.assertChain)
			assert.Equal(t, test.isFound, found)
			if test.isFound {
				assert.Equal(t, test.assertObsListLen, len(mapper.ObserverList))
			}

		})
	}
}

func TestKeeper_ObserversByChainAndType(t *testing.T) {
	tt := []struct {
		name             string
		mapper           []*types.ObserverMapper
		assertChain      common.ChainName
		assertObsListLen int
		isFound          bool
	}{
		{
			name:        "4 ETH InBoundTx Observers",
			mapper:      types.CreateObserverMapperList(1, common.GoerliChain()),
			assertChain: common.ChainName_goerli_localnet,
			isFound:     true,
		},
		{
			name:        "4 BTC InBoundTx Observers",
			mapper:      types.CreateObserverMapperList(1, common.BtcRegtestChain()),
			assertChain: common.ChainName_btc_regtest,
			isFound:     true,
		},
		{
			name: "Filter out from multiple mappers",
			mapper: append(append(types.CreateObserverMapperList(1, common.GoerliChain()),
				types.CreateObserverMapperList(1, common.ZetaPrivnetChain())...)),
			assertChain: common.ChainName_goerli_localnet,
			isFound:     true,
		},
		{
			name: "No Observers of expected Observation Chain",
			mapper: append(append(types.CreateObserverMapperList(1, common.GoerliChain()),
				types.CreateObserverMapperList(1, common.ZetaPrivnetChain())...)),
			assertChain: common.ChainName_btc_regtest,
			isFound:     false,
		},
	}

	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			k, ctx := keepertest.ObserverKeeper(t)

			for _, mapper := range test.mapper {
				k.SetObserverMapper(ctx, mapper)
			}
			goCtx := sdk.WrapSDKContext(ctx)
			msg := &types.QueryObserversByChainRequest{
				ObservationChain: test.assertChain.String(),
			}

			mapper, _ := k.ObserversByChain(goCtx, msg)
			if test.isFound {
				assert.NotEqual(t, "", mapper)
			}

		})
	}
}

func TestKeeper_GetAllObserverAddresses(t *testing.T) {
	mappers := append(append(types.CreateObserverMapperList(1, common.Chain{
		ChainName: common.ChainName_btc_mainnet,
		ChainId:   3,
	}),
		types.CreateObserverMapperList(1, common.Chain{
			ChainName: common.ChainName_polygon_mainnet,
			ChainId:   4,
		})...),
		types.CreateObserverMapperList(1, common.Chain{
			ChainName: common.ChainName_bsc_mainnet,
			ChainId:   5,
		})...)
	k, ctx := keepertest.ObserverKeeper(t)
	for _, mapper := range mappers {
		k.SetObserverMapper(ctx, mapper)
	}
	addresses := k.GetAllObserverAddresses(ctx)
	assert.Equal(t, 4, len(addresses))
}
