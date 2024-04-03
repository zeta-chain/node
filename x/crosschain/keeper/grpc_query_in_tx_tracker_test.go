package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_InTxTrackerAllByChain(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  1,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})
	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  2,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})

	res, err := k.InTxTrackerAllByChain(ctx, &types.QueryAllInTxTrackerByChainRequest{
		ChainId: 1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.InTxTracker))
}

func TestKeeper_InTxTrackerAll(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  1,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})
	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  2,
		TxHash:   sample.Hash().Hex(),
		CoinType: coin.CoinType_Gas,
	})

	res, err := k.InTxTrackerAll(ctx, &types.QueryAllInTxTrackersRequest{})
	require.NoError(t, err)
	require.Equal(t, 2, len(res.InTxTracker))
}
