package keeper

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"testing"
)

func TestKeeper_SupportedChains(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	list := []types.ObserverChain{
		types.ObserverChain_EthChainObserver,
		types.ObserverChain_BscChainObserver,
		types.ObserverChain_BtcChainObserver,
	}
	keeper.SetSupportedChain(ctx, types.SupportedChains{ChainList: list})
	getList, found := keeper.GetSupportedChains(ctx)
	assert.True(t, found)
	assert.Equal(t, list, getList.ChainList)
}
