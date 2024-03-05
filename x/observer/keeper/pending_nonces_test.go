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
		k, ctx, _ := keepertest.ObserverKeeper(t)
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
		k, ctx, _ := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 42)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		offset := 10
		rst, pageRes, err := k.GetAllPendingNoncesPaginated(ctx, &query.PageRequest{Offset: uint64(offset), CountTotal: true})
		require.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		require.Subset(t, nonces, rst)
		require.Len(t, rst, len(nonces)-offset)
		require.Equal(t, len(nonces), int(pageRes.Total))
	})
	t.Run("Get all pending nonces ", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)
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
	})
}
