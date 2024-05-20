package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_OutboundTrackerAllByChain(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutboundTrackerAllByChain(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
		})
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 2,
		})

		res, err := k.OutboundTrackerAllByChain(ctx, &types.QueryAllOutboundTrackerByChainRequest{
			Chain: 1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.OutboundTracker))
	})
}

func TestKeeper_OutboundTrackerAll(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutboundTrackerAll(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
		})

		res, err := k.OutboundTrackerAll(ctx, &types.QueryAllOutboundTrackerRequest{})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.OutboundTracker))
	})
}

func TestKeeper_OutboundTracker(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutboundTracker(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutboundTracker(ctx, &types.QueryGetOutboundTrackerRequest{
			ChainID: 1,
			Nonce:   1,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
			Nonce:   1,
		})

		res, err := k.OutboundTracker(ctx, &types.QueryGetOutboundTrackerRequest{
			ChainID: 1,
			Nonce:   1,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})
}
