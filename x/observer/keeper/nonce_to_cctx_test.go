package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetNonceToCctx(t *testing.T) {
	t.Run("Get nonce to cctx", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 1)
		for _, n := range nonceToCctxList {
			k.SetNonceToCctx(ctx, n)
		}
		for _, n := range nonceToCctxList {
			rst, found := k.GetNonceToCctx(ctx, n.Tss, n.ChainId, n.Nonce)
			assert.True(t, found)
			assert.Equal(t, n, rst)
		}
	})
	t.Run("Get nonce to cctx not found", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 1)
		for _, n := range nonceToCctxList {
			k.SetNonceToCctx(ctx, n)
		}
		_, found := k.GetNonceToCctx(ctx, "not_found", 1, 1)
		assert.False(t, found)
	})
	t.Run("Get all nonce to cctx", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 10)
		for _, n := range nonceToCctxList {
			k.SetNonceToCctx(ctx, n)
		}
		rst := k.GetAllNonceToCctx(ctx)
		assert.Equal(t, nonceToCctxList, rst)
	})
}
