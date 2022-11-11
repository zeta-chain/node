package keeper

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

func TestKeeper_SupportedChains(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	list := []*types.Chain{
		{
			ChainName: types.ChainName_Eth,
			ChainId:   1,
		},
		{
			ChainName: types.ChainName_Btc,
			ChainId:   2,
		},
		{
			ChainName: types.ChainName_BscMainnet,
			ChainId:   3,
		},
	}

	keeper.SetSupportedChain(ctx, types.SupportedChains{ChainList: list})
	getList, found := keeper.GetSupportedChains(ctx)
	assert.True(t, found)
	assert.Equal(t, list, getList.ChainList)
}
