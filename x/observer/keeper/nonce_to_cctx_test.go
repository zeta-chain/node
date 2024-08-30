package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_GetNonceToCctx(t *testing.T) {
	t.Run("Get nonce to cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 1)
		for _, n := range nonceToCctxList {
			k.SetNonceToCctx(ctx, n)
		}
		for _, n := range nonceToCctxList {
			rst, found := k.GetNonceToCctx(ctx, n.Tss, n.ChainId, n.Nonce)
			require.True(t, found)
			require.Equal(t, n, rst)
		}

		for _, n := range nonceToCctxList {
			k.RemoveNonceToCctx(ctx, n)
		}
		for _, n := range nonceToCctxList {
			_, found := k.GetNonceToCctx(ctx, n.Tss, n.ChainId, n.Nonce)
			require.False(t, found)
		}
	})
	t.Run("test nonce to cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		k.SetNonceToCctx(ctx, types.NonceToCctx{
			ChainId:   1337,
			Nonce:     0,
			CctxIndex: "0x705b88814b2a049e75b591fd80595c53f3bd9ddfb67ad06aa6965ed91023ee9a",
			Tss:       "zetapub1addwnpepq0akz8ene4z2mg3tghamr0m5eg3eeuqtjcfamkh5ecetua9u0pcyvjeyerd",
		})
		_, found := k.GetNonceToCctx(
			ctx,
			"zetapub1addwnpepq0akz8ene4z2mg3tghamr0m5eg3eeuqtjcfamkh5ecetua9u0pcyvjeyerd",
			1337,
			0,
		)
		require.True(t, found)
	})
	t.Run("Get nonce to cctx not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 1)
		for _, n := range nonceToCctxList {
			k.SetNonceToCctx(ctx, n)
		}
		_, found := k.GetNonceToCctx(ctx, "not_found", 1, 1)
		require.False(t, found)
	})
	t.Run("Get all nonce to cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 10)
		for _, n := range nonceToCctxList {
			k.SetNonceToCctx(ctx, n)
		}
		rst := k.GetAllNonceToCctx(ctx)
		require.Equal(t, nonceToCctxList, rst)
	})
}
