package keeper

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"testing"
)

func TestKeeper_SupportedChains(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	list := []string{
		"GANACHE",
		"GOERILI",
	}
	keeper.SetSupportedChain(ctx, types.SupportedChains{ChainList: list})
	getList, found := keeper.GetSupportedChains(ctx)
	assert.True(t, found)
	assert.Equal(t, list, getList.ChainList)
}
