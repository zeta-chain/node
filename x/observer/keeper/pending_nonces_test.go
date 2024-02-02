package keeper_test

import (
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_PendingNoncesAll(t *testing.T) {
	t.Run("Get all pending nonces paginated by limit", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 10)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		rst, pageRes, err := k.GetAllPendingNoncesPaginated(ctx, &query.PageRequest{Limit: 10, CountTotal: true})
		assert.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		assert.Equal(t, nonces, rst)
		assert.Equal(t, len(nonces), int(pageRes.Total))
	})
	t.Run("Get all pending nonces paginated by offset", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 42)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		offset := 10
		rst, pageRes, err := k.GetAllPendingNoncesPaginated(ctx, &query.PageRequest{Offset: uint64(offset), CountTotal: true})
		assert.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		assert.Subset(t, nonces, rst)
		assert.Len(t, rst, len(nonces)-offset)
		assert.Equal(t, len(nonces), int(pageRes.Total))
	})
	t.Run("Get all pending nonces ", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		nonces := sample.PendingNoncesList(t, "sample", 10)
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i].ChainId < nonces[j].ChainId
		})
		for _, nonce := range nonces {
			k.SetPendingNonces(ctx, nonce)
		}
		rst, err := k.GetAllPendingNonces(ctx)
		assert.NoError(t, err)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].ChainId < rst[j].ChainId
		})
		assert.Equal(t, nonces, rst)
	})
}
