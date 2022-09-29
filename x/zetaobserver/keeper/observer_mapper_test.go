package keeper

import (
	"github.com/stretchr/testify/assert"
	testdata "github.com/zeta-chain/zetacore/x/zetaobserver/testing"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"testing"
)

func TestKeeper_GetObserver(t *testing.T) {

	tt := []struct {
		name             string
		mapper           []types.ObserverMapper
		assertChain      types.ObserverChain
		assertObsType    types.ObservationType
		assertObsListLen int
		isFound          bool
	}{
		{
			name:             "4 eth Observers",
			mapper:           testdata.CreateObserverMapperList(1, types.ObserverChain_EthChainObserver, types.ObservationType_InboundTx),
			assertChain:      types.ObserverChain_EthChainObserver,
			assertObsType:    types.ObservationType_InboundTx,
			assertObsListLen: 4,
			isFound:          true,
		},
		{
			name: "Filter out from multiple mappers",
			mapper: append(append(testdata.CreateObserverMapperList(1, types.ObserverChain_EthChainObserver, types.ObservationType_InboundTx),
				testdata.CreateObserverMapperList(1, types.ObserverChain_EthChainObserver, types.ObservationType_OutBoundTx)...),
				testdata.CreateObserverMapperList(1, types.ObserverChain_BscChainObserver, types.ObservationType_OutBoundTx)...),
			assertChain:      types.ObserverChain_EthChainObserver,
			assertObsType:    types.ObservationType_InboundTx,
			assertObsListLen: 4,
			isFound:          true,
		},
		{
			name: "No Observers of expected Observation Chain",
			mapper: append(append(testdata.CreateObserverMapperList(1, types.ObserverChain_BTCChainObserver, types.ObservationType_InboundTx),
				testdata.CreateObserverMapperList(1, types.ObserverChain_PolygonChainObserver, types.ObservationType_OutBoundTx)...),
				testdata.CreateObserverMapperList(1, types.ObserverChain_BscChainObserver, types.ObservationType_OutBoundTx)...),
			assertChain:      types.ObserverChain_EthChainObserver,
			assertObsType:    types.ObservationType_InboundTx,
			assertObsListLen: 0,
			isFound:          false,
		},
		{
			name: "No Observers of expected Observation Type",
			mapper: append(append(testdata.CreateObserverMapperList(1, types.ObserverChain_BTCChainObserver, types.ObservationType_InboundTx),
				testdata.CreateObserverMapperList(1, types.ObserverChain_PolygonChainObserver, types.ObservationType_OutBoundTx)...),
				testdata.CreateObserverMapperList(1, types.ObserverChain_BscChainObserver, types.ObservationType_OutBoundTx)...),
			assertChain:      types.ObserverChain_EthChainObserver,
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
			mapper, found := keeper.GetObserverMapper(ctx, test.assertChain.String(), test.assertObsType.String())
			assert.Equal(t, test.isFound, found)
			if test.isFound {
				assert.Equal(t, test.assertObsListLen, len(mapper.ObserverList))
			}

		})
	}

}
