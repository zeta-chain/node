package keeper_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// createCctxWithNonceRange create in the store:
// mined cctx from nonce 0 to low
// pending cctx from low to high
// set pending nonces from low to higg
// return pending cctxs
func createCctxWithNonceRange(
	t *testing.T,
	ctx sdk.Context,
	k keeper.Keeper,
	lowPending int,
	highPending int,
	chainID int64,
	tss observertypes.TSS,
	zk keepertest.ZetaKeepers,
) (cctxs []*types.CrossChainTx) {
	for i := 0; i < lowPending; i++ {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", chainID, i))
		cctx.CctxStatus.Status = types.CctxStatus_OutboundMined
		cctx.InboundParams.SenderChainId = chainID
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetNonceToCctx(ctx, observertypes.NonceToCctx{
			ChainId:   chainID,
			Nonce:     int64(i),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})
	}
	for i := lowPending; i < highPending; i++ {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", chainID, i))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.SenderChainId = chainID
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetNonceToCctx(ctx, observertypes.NonceToCctx{
			ChainId:   chainID,
			Nonce:     int64(i),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})
		cctxs = append(cctxs, cctx)
	}
	zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
		ChainId:   chainID,
		NonceLow:  int64(lowPending),
		NonceHigh: int64(highPending),
		Tss:       tss.TssPubkey,
	})

	return
}

func TestKeeper_CctxListPending(t *testing.T) {
	t.Run("should fail for empty req", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.ListPendingCctx(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("should use max limit if limit is too high", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{Limit: keeper.MaxPendingCctxs + 1})
		require.ErrorContains(t, err, "tss not found")
	})

	t.Run("should fail if no TSS", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{Limit: 1})
		require.ErrorContains(t, err, "tss not found")
	})

	t.Run("should return empty list if no nonces", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		//  set TSS
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		_, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{Limit: 1})
		require.ErrorContains(t, err, "pending nonces not found")
	})

	t.Run("can retrieve pending cctx in range", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		cctxs := createCctxWithNonceRange(t, ctx, *k, 1000, 2000, chainID, tss, zk)

		res, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{ChainId: chainID, Limit: 100})
		require.NoError(t, err)
		require.Equal(t, 100, len(res.CrossChainTx))
		require.EqualValues(t, cctxs[0:100], res.CrossChainTx)
		require.EqualValues(t, uint64(1000), res.TotalPending)

		res, err = k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{ChainId: chainID})
		require.NoError(t, err)
		require.Equal(t, keeper.MaxPendingCctxs, len(res.CrossChainTx))
		require.EqualValues(t, cctxs[0:keeper.MaxPendingCctxs], res.CrossChainTx)
		require.EqualValues(t, uint64(1000), res.TotalPending)
	})

	t.Run("can retrieve pending cctx with range smaller than max", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		cctxs := createCctxWithNonceRange(t, ctx, *k, 1000, 1100, chainID, tss, zk)

		res, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{ChainId: chainID})
		require.NoError(t, err)
		require.Equal(t, 100, len(res.CrossChainTx))
		require.EqualValues(t, cctxs, res.CrossChainTx)
		require.EqualValues(t, uint64(100), res.TotalPending)
	})

	t.Run("can retrieve pending cctx with pending cctx below nonce low", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		cctxs := createCctxWithNonceRange(t, ctx, *k, 1000, 2000, chainID, tss, zk)

		// set some cctxs as pending below nonce
		cctx1, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("1337-940"))
		require.True(t, found)
		cctx1.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, cctx1)

		cctx2, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("1337-955"))
		require.True(t, found)
		cctx2.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, cctx2)

		res, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{ChainId: chainID, Limit: 100})
		require.NoError(t, err)
		require.Equal(t, 100, len(res.CrossChainTx))

		expectedCctxs := append([]*types.CrossChainTx{&cctx1, &cctx2}, cctxs[0:98]...)
		require.EqualValues(t, expectedCctxs, res.CrossChainTx)

		// pending nonce + 2
		require.EqualValues(t, uint64(1002), res.TotalPending)
	})

	t.Run("error if some before low nonce are missing", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		cctxs := createCctxWithNonceRange(t, ctx, *k, 1000, 2000, chainID, tss, zk)

		// set some cctxs as pending below nonce
		cctx1, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("1337-940"))
		require.True(t, found)
		cctx1.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, cctx1)

		cctx2, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("1337-955"))
		require.True(t, found)
		cctx2.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, cctx2)

		res, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{ChainId: chainID, Limit: 100})
		require.NoError(t, err)
		require.Equal(t, 100, len(res.CrossChainTx))

		expectedCctxs := append([]*types.CrossChainTx{&cctx1, &cctx2}, cctxs[0:98]...)
		require.EqualValues(t, expectedCctxs, res.CrossChainTx)

		// pending nonce + 2
		require.EqualValues(t, uint64(1002), res.TotalPending)
	})
}

func TestKeeper_ZetaAccounting(t *testing.T) {
	t.Run("should error if zeta accounting not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.ZetaAccounting(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return zeta accounting if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.SetZetaAccounting(ctx, types.ZetaAccounting{
			AbortedZetaAmount: sdkmath.NewUint(100),
		})

		res, err := k.ZetaAccounting(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, &types.QueryZetaAccountingResponse{
			AbortedZetaAmount: sdkmath.NewUint(100).String(),
		}, res)
	})
}

func TestKeeper_CctxByNonce(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.CctxByNonce(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
			ChainID: 1,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if nonce to cctx not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)

		res, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
			ChainID: chainID,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if crosschain tx not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		nonce := 1000
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", chainID, nonce))

		zk.ObserverKeeper.SetNonceToCctx(ctx, observertypes.NonceToCctx{
			ChainId:   chainID,
			Nonce:     int64(nonce),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})

		res, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
			ChainID: chainID,
			Nonce:   uint64(nonce),
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if crosschain tx found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		nonce := 1000
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", chainID, nonce))

		zk.ObserverKeeper.SetNonceToCctx(ctx, observertypes.NonceToCctx{
			ChainId:   chainID,
			Nonce:     int64(nonce),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})
		k.SetCrossChainTx(ctx, *cctx)

		res, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
			ChainID: chainID,
			Nonce:   uint64(nonce),
		})
		require.NoError(t, err)
		require.Equal(t, cctx, res.CrossChainTx)

		// ensure that LastUpdateTimestamp is set to current block time
		require.Equal(t, res.CrossChainTx.CctxStatus.LastUpdateTimestamp, ctx.BlockTime().Unix())
	})
}

func assertCctxIndexEqual(t *testing.T, expectedCctxs []*types.CrossChainTx, cctxs []*types.CrossChainTx) {
	t.Helper()
	require.Equal(t, len(expectedCctxs), len(cctxs), "slice lengths not equal")
	for i, expectedCctx := range expectedCctxs {
		require.Equal(t, expectedCctx.Index, cctxs[i].Index, "index mismatch at %v", i)
	}
}

func sortByIndex(l *types.CrossChainTx, r *types.CrossChainTx) int {
	return strings.Compare(l.Index, r.Index)
}

func TestKeeper_CctxAll(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{})
		require.NoError(t, err)
	})

	t.Run("default page size", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		_ = createCctxWithNonceRange(t, ctx, *k, 1000, 2000, chainID, tss, zk)

		res, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{})
		require.NoError(t, err)
		require.Len(t, res.CrossChainTx, keeper.DefaultPageSize)
	})

	t.Run("page size provided", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		_ = createCctxWithNonceRange(t, ctx, *k, 1000, 2000, chainID, tss, zk)
		testPageSize := 200

		res, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{
			Pagination: &query.PageRequest{
				Limit: uint64(testPageSize),
			},
		})
		require.NoError(t, err)
		require.Len(t, res.CrossChainTx, testPageSize)
	})

	t.Run("basic descending ordering", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		createdCctx := createCctxWithNonceRange(t, ctx, *k, 0, 10, chainID, tss, zk)

		res, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{})
		require.NoError(t, err)
		slices.Reverse(createdCctx)
		assertCctxIndexEqual(t, createdCctx, res.CrossChainTx)

		// also assert unordered query return same number of results
		resUnordered, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{
			Unordered: true,
		})
		require.NoError(t, err)
		require.Len(t, res.CrossChainTx, len(resUnordered.CrossChainTx))
	})

	t.Run("basic ascending ordering", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		createdCctx := createCctxWithNonceRange(t, ctx, *k, 0, 10, chainID, tss, zk)

		res, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{
			Pagination: &query.PageRequest{
				Reverse: true,
			},
		})
		require.NoError(t, err)
		assertCctxIndexEqual(t, createdCctx, res.CrossChainTx)

		// also assert unordered query return same number of results
		resUnordered, err := k.CctxAll(ctx, &types.QueryAllCctxRequest{
			Unordered: true,
		})
		require.NoError(t, err)
		require.Len(t, res.CrossChainTx, len(resUnordered.CrossChainTx))
	})
}
