package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func createNSend(keeper *Keeper, ctx sdk.Context, n int) []types.CrossChainTx {
	items := make([]types.CrossChainTx, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].InBoundTxParams = &types.InBoundTxParams{
			Sender:                   fmt.Sprintf("%d", i),
			SenderChain:              fmt.Sprintf("%d", i),
			InBoundTxObservedHash:    fmt.Sprintf("%d", i),
			InBoundTxObservedHeight:  uint64(i),
			InBoundTxFinalizedHeight: uint64(i),
		}
		items[i].OutBoundTxParams = &types.OutBoundTxParams{
			Receiver:               fmt.Sprintf("%d", i),
			ReceiverChain:          fmt.Sprintf("%d", i),
			Broadcaster:            uint64(i),
			OutBoundTxHash:         fmt.Sprintf("%d", i),
			OutBoundTxTSSNonce:     uint64(i),
			OutBoundTxGasLimit:     uint64(i),
			OutBoundTxGasPrice:     fmt.Sprintf("%d", i),
			OutBoundTXReceiveIndex: fmt.Sprintf("%d", i),
		}
		items[i].CctxStatus = &types.Status{
			Status:              types.CctxStatus_PendingInbound,
			StatusMessage:       "any",
			LastUpdateTimestamp: 0,
		}
		items[i].ZetaBurnt = sdk.OneUint()
		items[i].ZetaMint = sdk.OneUint()
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetCrossChainTx(ctx, items[i])
	}
	return items
}

func TestSendGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 1)
	for _, item := range items {
		rst, found := keeper.GetCrossChainTx(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestSendRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveCrossChainTx(ctx, item.Index)
		_, found := keeper.GetCrossChainTx(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestSendGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllCrossChainTx(ctx))
}
