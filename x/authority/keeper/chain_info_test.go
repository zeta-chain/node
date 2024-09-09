package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestKeeper_SetChainInfo(t *testing.T) {
	k, ctx := keepertest.AuthorityKeeper(t)
	chainInfo := sample.ChainInfo(42)

	_, found := k.GetChainInfo(ctx)
	require.False(t, found)

	k.SetChainInfo(ctx, chainInfo)

	// Check policy is set
	got, found := k.GetChainInfo(ctx)
	require.True(t, found)
	require.Equal(t, chainInfo, got)

	// Can set policies again
	newChainInfo := sample.ChainInfo(84)
	require.NotEqual(t, chainInfo, newChainInfo)
	k.SetChainInfo(ctx, newChainInfo)
	got, found = k.GetChainInfo(ctx)
	require.True(t, found)
	require.Equal(t, newChainInfo, got)
}

func TestKeeper_GetChainList(t *testing.T) {
	k, ctx := keepertest.AuthorityKeeper(t)

	// Empty list
	list := k.GetAdditionalChainList(ctx)
	require.Empty(t, list)

	// Set chain info
	chainInfo := sample.ChainInfo(42)
	k.SetChainInfo(ctx, chainInfo)

	// Check list
	list = k.GetAdditionalChainList(ctx)
	require.Equal(t, chainInfo.Chains, list)
}
