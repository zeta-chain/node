package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestChainNoncesGet(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	items := sample.ChainNoncesList(10)
	for _, item := range items {
		k.SetChainNonces(ctx, item)
	}
	for _, item := range items {
		rst, found := k.GetChainNonces(ctx, item.ChainId)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestChainNoncesRemove(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	items := sample.ChainNoncesList(10)
	for _, item := range items {
		k.SetChainNonces(ctx, item)
	}
	for _, item := range items {
		k.RemoveChainNonces(ctx, item.ChainId)
		_, found := k.GetChainNonces(ctx, item.ChainId)
		require.False(t, found)
	}
}

func TestChainNoncesGetAll(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	items := sample.ChainNoncesList(10)
	for _, item := range items {
		k.SetChainNonces(ctx, item)
	}
	require.Equal(t, items, k.GetAllChainNonces(ctx))
}
