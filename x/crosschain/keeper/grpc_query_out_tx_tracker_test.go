package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_OutTxTrackerAllByChain(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutTxTrackerAllByChain(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: 1,
		})
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: 2,
		})

		res, err := k.OutTxTrackerAllByChain(ctx, &types.QueryAllOutTxTrackerByChainRequest{
			Chain: 1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.OutTxTracker))
	})
}

func TestKeeper_OutTxTrackerAll(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutTxTrackerAll(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: 1,
		})

		res, err := k.OutTxTrackerAll(ctx, &types.QueryAllOutTxTrackerRequest{})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.OutTxTracker))
	})
}

func TestKeeper_OutTxTracker(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.OutTxTracker(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: 1,
			Nonce:   1,
		})

		res, err := k.OutTxTracker(ctx, &types.QueryGetOutTxTrackerRequest{
			ChainID: 1,
			Nonce:   1,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})
}
