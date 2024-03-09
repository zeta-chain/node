package v6_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	v6 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v6"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("sucessfull migrate cctx from v14 to v15", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		v14cctxList := make([]*types.CrossChainTxV14, 10)
		for i := 0; i < 10; i++ {
			v14cctxList[i] = sample.CrossChainTxV14(t, fmt.Sprintf("%d-%s", i, "v14"))
			SetCrossChainTxV14(ctx, *v14cctxList[i], k)
		}
		err := v6.MigrateStore(ctx, k)
		require.NoError(t, err)
		cctxListv15 := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxListv15, 10)
		sort.Slice(cctxListv15, func(i, j int) bool {
			return cctxListv15[i].Index < cctxListv15[j].Index
		})
		sort.Slice(v14cctxList, func(i, j int) bool {
			return v14cctxList[i].Index < v14cctxList[j].Index
		})
		for i := 0; i < 10; i++ {
			require.Equal(t, v14cctxList[i].Index, cctxListv15[i].Index)
			require.Equal(t, v14cctxList[i].Creator, cctxListv15[i].Creator)
			require.Equal(t, v14cctxList[i].ZetaFees, cctxListv15[i].ZetaFees)
			require.Equal(t, v14cctxList[i].RelayedMessage, cctxListv15[i].RelayedMessage)
			require.Equal(t, v14cctxList[i].CctxStatus, cctxListv15[i].CctxStatus)
			require.Equal(t, v14cctxList[i].InboundTxParams.Sender, cctxListv15[i].InboundTxParams.Sender)
			require.Equal(t, v14cctxList[i].InboundTxParams.SenderChainId, cctxListv15[i].InboundTxParams.SenderChainId)
			require.Equal(t, v14cctxList[i].InboundTxParams.TxOrigin, cctxListv15[i].InboundTxParams.TxOrigin)
			require.Equal(t, v14cctxList[i].InboundTxParams.Asset, cctxListv15[i].InboundTxParams.Asset)
			require.Equal(t, v14cctxList[i].InboundTxParams.Amount, cctxListv15[i].InboundTxParams.Amount)
			require.Equal(t, v14cctxList[i].InboundTxParams.InboundTxObservedHash, cctxListv15[i].InboundTxParams.InboundTxObservedHash)
			require.Equal(t, v14cctxList[i].InboundTxParams.InboundTxObservedExternalHeight, cctxListv15[i].InboundTxParams.InboundTxObservedExternalHeight)
			require.Equal(t, v14cctxList[i].InboundTxParams.InboundTxBallotIndex, cctxListv15[i].InboundTxParams.InboundTxBallotIndex)
			require.Equal(t, v14cctxList[i].InboundTxParams.InboundTxFinalizedZetaHeight, cctxListv15[i].InboundTxParams.InboundTxFinalizedZetaHeight)
			require.Equal(t, v14cctxList[i].InboundTxParams.CoinType, cctxListv15[i].CoinType)
			require.Len(t, v14cctxList[i].OutboundTxParams, len(cctxListv15[i].OutboundTxParams))
			for j := 0; j < len(cctxListv15[i].OutboundTxParams); j++ {
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].Receiver, cctxListv15[i].OutboundTxParams[j].Receiver)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].ReceiverChainId, cctxListv15[i].OutboundTxParams[j].ReceiverChainId)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].Amount, cctxListv15[i].OutboundTxParams[j].Amount)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxTssNonce, cctxListv15[i].OutboundTxParams[j].OutboundTxTssNonce)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxGasLimit, cctxListv15[i].OutboundTxParams[j].OutboundTxGasLimit)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxGasPrice, cctxListv15[i].OutboundTxParams[j].OutboundTxGasPrice)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxHash, cctxListv15[i].OutboundTxParams[j].OutboundTxHash)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxBallotIndex, cctxListv15[i].OutboundTxParams[j].OutboundTxBallotIndex)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxObservedExternalHeight, cctxListv15[i].OutboundTxParams[j].OutboundTxObservedExternalHeight)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxGasUsed, cctxListv15[i].OutboundTxParams[j].OutboundTxGasUsed)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].OutboundTxEffectiveGasPrice, cctxListv15[i].OutboundTxParams[j].OutboundTxEffectiveGasPrice)
				require.Equal(t, v14cctxList[i].OutboundTxParams[j].CoinType, cctxListv15[i].CoinType)
			}
		}
	})
}

func SetCrossChainTxV14(ctx sdk.Context, cctx types.CrossChainTxV14, k *keeper.Keeper) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), p)
	b := k.GetCodec().MustMarshal(&cctx)
	store.Set(types.KeyPrefix(cctx.Index), b)
}
