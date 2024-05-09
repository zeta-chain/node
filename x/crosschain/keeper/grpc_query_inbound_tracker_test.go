package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
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
