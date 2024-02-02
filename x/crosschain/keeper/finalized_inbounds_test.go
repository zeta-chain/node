package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_IsFinalizedInbound(t *testing.T) {
	t.Run("check true for finalized inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxHash := sample.Hash().String()
		chainID := sample.Chain(5).ChainId
		eventIndex := sample.EventIndex()
		k.AddFinalizedInbound(ctx, intxHash, chainID, eventIndex)
		assert.True(t, k.IsFinalizedInbound(ctx, intxHash, chainID, eventIndex))
	})
	t.Run("check false for non-finalized inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxHash := sample.Hash().String()
		chainID := sample.Chain(5).ChainId
		eventIndex := sample.EventIndex()
		assert.False(t, k.IsFinalizedInbound(ctx, intxHash, chainID, eventIndex))
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
			assert.True(t, k.IsFinalizedInbound(ctx, txHashList[i], chainIdList[i], eventIndexList[i]))
		}
	})
}

func TestKeeper_AddFinalizedInbound(t *testing.T) {
	t.Run("check add finalized inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxHash := sample.Hash().String()
		chainID := sample.Chain(5).ChainId
		eventIndex := sample.EventIndex()
		k.AddFinalizedInbound(ctx, intxHash, chainID, eventIndex)
		assert.True(t, k.IsFinalizedInbound(ctx, intxHash, chainID, eventIndex))
	})
}

func TestKeeper_GetAllFinalizedInbound(t *testing.T) {
	t.Run("check empty list", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		list := k.GetAllFinalizedInbound(ctx)
		assert.Empty(t, list)
	})
	t.Run("check list", func(t *testing.T) {
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
		list := k.GetAllFinalizedInbound(ctx)
		assert.Equal(t, listSize, len(list))
		for i := 0; i < listSize; i++ {
			assert.Contains(t, list, types.FinalizedInboundKey(txHashList[i], chainIdList[i], eventIndexList[i]))
		}
		assert.NotContains(t, list, types.FinalizedInboundKey(sample.Hash().String(), sample.Chain(5).ChainId, sample.EventIndex()))
	})
}
