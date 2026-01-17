package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestKeeper_InboundTrackerAllByChain(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	k.SetInboundTracker(ctx, types.InboundTracker{
		ChainId:  1,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})
	k.SetInboundTracker(ctx, types.InboundTracker{
		ChainId:  2,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})

	res, err := k.InboundTrackerAllByChain(ctx, &types.QueryAllInboundTrackerByChainRequest{
		ChainId: 1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.InboundTracker))
}

func TestKeeper_InboundTrackerAll(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	k.SetInboundTracker(ctx, types.InboundTracker{
		ChainId:  1,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})
	k.SetInboundTracker(ctx, types.InboundTracker{
		ChainId:  2,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})

	res, err := k.InboundTrackerAll(ctx, &types.QueryAllInboundTrackersRequest{})
	require.NoError(t, err)
	require.Equal(t, 2, len(res.InboundTracker))
}

func TestKeeper_InboundTracker(t *testing.T) {
	t.Run("successfully get inbound tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		hash := sample.Hash().Hex()
		chainID := chains.GoerliLocalnet.ChainId
		k.SetInboundTracker(ctx, types.InboundTracker{
			ChainId:  chainID,
			TxHash:   hash,
			CoinType: coin.CoinType_Gas,
		})

		res, err := k.InboundTracker(ctx, &types.QueryInboundTrackerRequest{
			ChainId: chainID,
			TxHash:  hash,
		})
		require.NoError(t, err)
		require.NotNil(t, res.InboundTracker)
	})

	t.Run("inbound tracker not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		hash := sample.Hash().Hex()
		chainID := chains.GoerliLocalnet.ChainId

		res, err := k.InboundTracker(ctx, &types.QueryInboundTrackerRequest{
			ChainId: chainID,
			TxHash:  hash,
		})
		require.ErrorContains(t, err, "not found")
		require.Nil(t, res)
	})
}
