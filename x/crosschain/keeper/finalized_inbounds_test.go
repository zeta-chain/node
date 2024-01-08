package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_IsFinalizedInbound(t *testing.T) {
	t.Run("check true for finalized inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxHash := sample.Hash().String()
		chainID := sample.Chain(5).ChainId
		eventIndex := sample.EventIndex()
		k.AddFinalizedInbound(ctx, intxHash, chainID, eventIndex)
		require.True(t, k.IsFinalizedInbound(ctx, intxHash, chainID, eventIndex))
	})
	t.Run("check false for non-finalized inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxHash := sample.Hash().String()
		chainID := sample.Chain(5).ChainId
		eventIndex := sample.EventIndex()
		require.False(t, k.IsFinalizedInbound(ctx, intxHash, chainID, eventIndex))
	})
	t.Run("check true for finalized inbound list", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		listSize := 1000
		txHashList := make([]string, listSize)
		chainIdList := make([]int64, listSize)
		eventIndexList := make([]uint64, listSize)
		for i := 0; i < listSize; i++ {
			txHashList[i] = sample.Hash().String()
			chainIdList[i] = sample.Chain(5).ChainId
			eventIndexList[i] = sample.EventIndex()
			k.AddFinalizedInbound(ctx, txHashList[i], chainIdList[i], eventIndexList[i])
		}
		for i := 0; i < listSize; i++ {
			require.True(t, k.IsFinalizedInbound(ctx, txHashList[i], chainIdList[i], eventIndexList[i]))
		}
	})
}
