package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"testing"
)

func TestKeeper_GetObserver(t *testing.T) {

	tt := []struct {
		name             string
		mapper           []*types.ObserverMapper
		assertChain      types.ObserverChain
		assertObsType    types.ObservationType
		assertObsListLen int
		isFound          bool
	}{
		{
			name:             "4 eth Observers",
			mapper:           types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_InBoundTx),
			assertChain:      types.ObserverChain_Eth,
			assertObsType:    types.ObservationType_InBoundTx,
			assertObsListLen: 4,
			isFound:          true,
		},
		{
			name: "Filter out from multiple mappers",
			mapper: append(append(types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_InBoundTx),
				types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_OutBoundTx)...),
				types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_OutBoundTx)...),
			assertChain:      types.ObserverChain_Eth,
			assertObsType:    types.ObservationType_InBoundTx,
			assertObsListLen: 4,
			isFound:          true,
		},
		{
			name: "No Observers of expected Observation Chain",
			mapper: append(append(types.CreateObserverMapperList(1, types.ObserverChain_Btc, types.ObservationType_InBoundTx),
				types.CreateObserverMapperList(1, types.ObserverChain_Polygon, types.ObservationType_OutBoundTx)...),
				types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_OutBoundTx)...),
			assertChain:      types.ObserverChain_Eth,
			assertObsType:    types.ObservationType_InBoundTx,
			assertObsListLen: 0,
			isFound:          false,
		},
		{
			name: "No Observers of expected Observation Type",
			mapper: append(append(types.CreateObserverMapperList(1, types.ObserverChain_Btc, types.ObservationType_InBoundTx),
				types.CreateObserverMapperList(1, types.ObserverChain_Polygon, types.ObservationType_OutBoundTx)...),
				types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_OutBoundTx)...),
			assertChain:      types.ObserverChain_Eth,
			assertObsType:    types.ObservationType_GasPrice,
			assertObsListLen: 0,
			isFound:          false,
		},
	}

	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			keeper, ctx := SetupKeeper(t)
			for _, mapper := range test.mapper {
				keeper.SetObserverMapper(ctx, mapper)
			}
			mapper, found := keeper.GetObserverMapper(ctx, test.assertChain, test.assertObsType.String())
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
		assertChain      string
		assertObsType    string
		assertObsListLen int
		isFound          bool
	}{
		{
			name:          "4 ETH InBoundTx Observers",
			mapper:        types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_InBoundTx),
			assertChain:   "EthChainObserver",
			assertObsType: "InBoundTx",
			isFound:       true,
		},
		{
			name:          "4 ETH OutBoundTx Observers",
			mapper:        types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_OutBoundTx),
			assertChain:   "EthChainObserver",
			assertObsType: "OutBoundTx",
			isFound:       true,
		},
		{
			name:          "4 BSC InBoundTx Observers",
			mapper:        types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_InBoundTx),
			assertChain:   "BscChainObserver",
			assertObsType: "InBoundTx",
			isFound:       true,
		},
		{
			name:          "4 POLYGON OutBoundTx Observers",
			mapper:        types.CreateObserverMapperList(1, types.ObserverChain_Polygon, types.ObservationType_OutBoundTx),
			assertChain:   "PolygonChainObserver",
			assertObsType: "OutBoundTx",
			isFound:       true,
		},
		{
			name: "Filter out from multiple mappers",
			mapper: append(append(types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_InBoundTx),
				types.CreateObserverMapperList(1, types.ObserverChain_Eth, types.ObservationType_OutBoundTx)...),
				types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_OutBoundTx)...),
			assertChain:   "EthChainObserver",
			assertObsType: "InBoundTx",
			isFound:       true,
		},
		{
			name: "No Observers of expected Observation Chain",
			mapper: append(append(types.CreateObserverMapperList(1, types.ObserverChain_Btc, types.ObservationType_InBoundTx),
				types.CreateObserverMapperList(1, types.ObserverChain_Polygon, types.ObservationType_OutBoundTx)...),
				types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_OutBoundTx)...),
			assertChain:   "EthChainObserver",
			assertObsType: "InBoundTx",
			isFound:       false,
		},
		{
			name: "No Observers of expected Observation Type",
			mapper: append(append(types.CreateObserverMapperList(1, types.ObserverChain_Btc, types.ObservationType_InBoundTx),
				types.CreateObserverMapperList(1, types.ObserverChain_Polygon, types.ObservationType_OutBoundTx)...),
				types.CreateObserverMapperList(1, types.ObserverChain_Bsc, types.ObservationType_OutBoundTx)...),
			assertChain:   "EthChainObserver",
			assertObsType: "GasPrice",
			isFound:       false,
		},
	}

	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			keeper, ctx := SetupKeeper(t)
			for _, mapper := range test.mapper {
				keeper.SetObserverMapper(ctx, mapper)
			}
			goCtx := sdk.WrapSDKContext(ctx)
			msg := &types.QueryObserversByChainAndTypeRequest{
				ObservationChain: test.assertChain,
				ObservationType:  test.assertObsType,
			}

			mapper, _ := keeper.ObserversByChainAndType(goCtx, msg)
			if test.isFound {
				assert.NotEqual(t, "", mapper)
			}

		})
	}
}
