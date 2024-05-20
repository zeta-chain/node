package v4_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	v4 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v4"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("test migrate store add zeta accounting", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		amountZeta := SetRandomCctx(ctx, *k)
		err := v4.MigrateStore(ctx, k.GetObserverKeeper(), k)
		require.NoError(t, err)
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, amountZeta, zetaAccounting.AbortedZetaAmount)
	})
	t.Run("test migrate store move tss from cross chain to observer", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		tss1 := sample.Tss()
		tss2 := sample.Tss()
		tss2.FinalizedZetaHeight = tss2.FinalizedZetaHeight + 10
		tss1Bytes := k.GetCodec().MustMarshal(&tss1)
		tss2Bytes := k.GetCodec().MustMarshal(&tss2)

		store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), observertypes.KeyPrefix(observertypes.TSSKey))
		store.Set([]byte{0}, tss1Bytes)

		store = prefix.NewStore(ctx.KVStore(k.GetStoreKey()), observertypes.KeyPrefix(observertypes.TSSHistoryKey))
		store.Set(observertypes.KeyPrefix(fmt.Sprintf("%d", tss1.FinalizedZetaHeight)), tss1Bytes)
		store.Set(observertypes.KeyPrefix(fmt.Sprintf("%d", tss2.FinalizedZetaHeight)), tss2Bytes)

		err := v4.MigrateStore(ctx, k.GetObserverKeeper(), k)
		require.NoError(t, err)

		tss, found := zk.ObserverKeeper.GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tss1, tss)
		tssHistory := k.GetObserverKeeper().GetAllTSS(ctx)
		require.Equal(t, 2, len(tssHistory))
		require.Equal(t, tss1, tssHistory[0])
		require.Equal(t, tss2, tssHistory[1])
	})
	t.Run("test migrate store move nonce from cross chain to observer", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		chainNonces := sample.ChainNoncesList(t, 10)
		pendingNonces := sample.PendingNoncesList(t, "sample", 10)
		nonceToCctxList := sample.NonceToCctxList(t, "sample", 10)
		store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), types.KeyPrefix(observertypes.ChainNoncesKey))
		for _, nonce := range chainNonces {
			store.Set([]byte(nonce.Index), k.GetCodec().MustMarshal(&nonce))
		}
		store = prefix.NewStore(ctx.KVStore(k.GetStoreKey()), types.KeyPrefix(observertypes.PendingNoncesKeyPrefix))
		for _, nonce := range pendingNonces {
			store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d", nonce.Tss, nonce.ChainId)), k.GetCodec().MustMarshal(&nonce))
		}
		store = prefix.NewStore(ctx.KVStore(k.GetStoreKey()), types.KeyPrefix(observertypes.NonceToCctxKeyPrefix))
		for _, nonce := range nonceToCctxList {
			store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonce.Tss, nonce.ChainId, nonce.Nonce)), k.GetCodec().MustMarshal(&nonce))
		}
		err := v4.MigrateStore(ctx, k.GetObserverKeeper(), k)
		require.NoError(t, err)
		pn, err := k.GetObserverKeeper().GetAllPendingNonces(ctx)
		require.NoError(t, err)
		sort.Slice(pn, func(i, j int) bool {
			return pn[i].ChainId < pn[j].ChainId
		})
		sort.Slice(pendingNonces, func(i, j int) bool {
			return pendingNonces[i].ChainId < pendingNonces[j].ChainId
		})
		require.Equal(t, pendingNonces, pn)
		require.Equal(t, chainNonces, k.GetObserverKeeper().GetAllChainNonces(ctx))
		require.Equal(t, nonceToCctxList, k.GetObserverKeeper().GetAllNonceToCctx(ctx))
	})
}

func TestSetBitcoinFinalizedInbound(t *testing.T) {
	t.Run("test setting finalized inbound for Bitcoin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// set some cctxs with Bitcoin and non-Bitcoin chains
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "0",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.GoerliChain.ChainId,
				ObservedHash:  "0xaaa",
			},
		})
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "1",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.BtcMainnetChain.ChainId,
				ObservedHash:  "0x111",
			},
		})
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "2",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.EthChain.ChainId,
				ObservedHash:  "0xbbb",
			},
		})
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "3",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.BtcTestNetChain.ChainId,
				ObservedHash:  "0x222",
			},
		})
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "4",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.BtcTestNetChain.ChainId,
				ObservedHash:  "0x333",
			},
		})
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "5",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.MumbaiChain.ChainId,
				ObservedHash:  "0xccc",
			},
		})
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index: "6",
			InboundParams: &types.InboundParams{
				SenderChainId: chains.BtcRegtestChain.ChainId,
				ObservedHash:  "0x444",
			},
		})

		// migration
		v4.SetBitcoinFinalizedInbound(ctx, k)

		// check finalized inbound
		require.False(t, k.IsFinalizedInbound(ctx, "0xaaa", chains.GoerliChain.ChainId, 0))
		require.False(t, k.IsFinalizedInbound(ctx, "0xbbb", chains.EthChain.ChainId, 0))
		require.False(t, k.IsFinalizedInbound(ctx, "0xccc", chains.MumbaiChain.ChainId, 0))
		require.True(t, k.IsFinalizedInbound(ctx, "0x111", chains.BtcMainnetChain.ChainId, 0))
		require.True(t, k.IsFinalizedInbound(ctx, "0x222", chains.BtcTestNetChain.ChainId, 0))
		require.True(t, k.IsFinalizedInbound(ctx, "0x333", chains.BtcTestNetChain.ChainId, 0))
		require.True(t, k.IsFinalizedInbound(ctx, "0x444", chains.BtcRegtestChain.ChainId, 0))

	})
}

func SetRandomCctx(ctx sdk.Context, k keeper.Keeper) sdkmath.Uint {
	totalZeta := sdkmath.ZeroUint()

	i := 0
	r := rand.New(rand.NewSource(9))
	for ; i < 10; i++ {
		amount := sdkmath.NewUint(uint64(r.Uint32()))
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index:      fmt.Sprintf("%d", i),
			CctxStatus: &types.Status{Status: types.CctxStatus_Aborted},
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Zeta,
			},
			OutboundParams: []*types.OutboundParams{{
				Amount:   amount,
				CoinType: coin.CoinType_Zeta,
			}},
		})
		totalZeta = totalZeta.Add(amount)
	}
	return totalZeta
}
