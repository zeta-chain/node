package keeper_test

import (
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_PendingNoncesAll(t *testing.T) {
	t.Run("Get all pending nonces paginated by limit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 10)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		rst, pageRes, err := k.GetAllPendingNoncesPaginated(ctx, &query.PageRequest{Limit: 10, CountTotal: true})
		require.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		require.Equal(t, nonces, rst)
		require.Equal(t, len(nonces), int(pageRes.Total))
	})
	t.Run("Get all pending nonces paginated by offset", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 42)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		offset := 10
		rst, pageRes, err := k.GetAllPendingNoncesPaginated(
			ctx,
			&query.PageRequest{Offset: uint64(offset), CountTotal: true},
		)
		require.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		require.Subset(t, nonces, rst)
		require.Len(t, rst, len(nonces)-offset)
		require.Equal(t, len(nonces), int(pageRes.Total))
	})
	t.Run("Get all pending nonces ", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 10)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		rst, err := k.GetAllPendingNonces(ctx)
		require.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		require.Equal(t, nonces, rst)

		k.RemovePendingNonces(ctx, nonces[0])
		rst, err = k.GetAllPendingNonces(ctx)
		require.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		require.Equal(t, nonces[1:], rst)
	})
}

func TestKeeper_SetTssAndUpdateNonce(t *testing.T) {
	t.Run("should set tss and update nonces", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		_, found := k.GetTSS(ctx)
		require.False(t, found)
		pendingNonces, err := k.GetAllPendingNonces(ctx)
		require.NoError(t, err)
		require.Empty(t, pendingNonces)
		chainNonces := k.GetAllChainNonces(ctx)
		require.NoError(t, err)
		require.Empty(t, chainNonces)

		tss := sample.Tss()
		// core params list but chain not in list
		setSupportedChain(ctx, *k, getValidEthChainIDWithIndex(t, 0))
		k.SetTssAndUpdateNonce(ctx, tss)

		_, found = k.GetTSS(ctx)
		require.True(t, found)
		pendingNonces, err = k.GetAllPendingNonces(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(pendingNonces))
		chainNonces = k.GetAllChainNonces(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(chainNonces))
	})
}

func TestKeeper_RemoveFromPendingNonces(t *testing.T) {
	t.Run("should remove from pending nonces", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 10)
		tss := sample.Tss()
		// make nonces and pubkey deterministic for test
		for i := range nonces {
			nonces[i].NonceLow = int64(i)
			nonces[i].NonceHigh = nonces[i].NonceLow + 3
			nonces[i].Tss = tss.TssPubkey
		}
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}

		k.RemoveFromPendingNonces(ctx, tss.TssPubkey, 1, 1)
		pendingNonces, err := k.GetAllPendingNonces(ctx)
		require.NoError(t, err)
		nonceUpdated := false
		for _, pn := range pendingNonces {
			if pn.ChainId == 1 {
				require.Equal(t, int64(2), pn.NonceLow)
				nonceUpdated = true
			}
		}
		require.True(t, nonceUpdated)
	})
}
