package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// SetCctxAndNonceToCctxAndInTxHashToCctx does the following things in one function:
// 1. set the cctx in the store
// 2. set the mapping inTxHash -> cctxIndex , one inTxHash can be connected to multiple cctxindex
// 3. set the mapping nonce => cctx

func (k Keeper) SetCctxAndNonceToCctxAndInTxHashToCctx(ctx sdk.Context, send types.CrossChainTx) {
	k.SetCrossChainTx(ctx, send)

	// set mapping inTxHash -> cctxIndex
	in, _ := k.GetInTxHashToCctx(ctx, send.InboundTxParams.InboundTxObservedHash)
	in.InTxHash = send.InboundTxParams.InboundTxObservedHash
	found := false
	for _, cctxIndex := range in.CctxIndex {
		if cctxIndex == send.Index {
			found = true
			break
		}
	}
	if !found {
		in.CctxIndex = append(in.CctxIndex, send.Index)
	}
	k.SetInTxHashToCctx(ctx, in)

	tss, found := k.GetTSS(ctx)
	if !found {
		return
	}
	// set mapping nonce => cctxIndex
	if send.CctxStatus.Status == types.CctxStatus_PendingOutbound || send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		k.SetNonceToCctx(ctx, types.NonceToCctx{
			ChainId: send.GetCurrentOutTxParam().ReceiverChainId,
			// #nosec G701 always in range
			Nonce:     int64(send.GetCurrentOutTxParam().OutboundTxTssNonce),
			CctxIndex: send.Index,
			Tss:       tss.TssPubkey,
		})
	}
	if send.CctxStatus.Status == types.CctxStatus_Aborted {
		k.AddZetaAccounting(ctx, send.GetCurrentOutTxParam().Amount)
	}
}

// SetCrossChainTx set a specific send in the store from its index
func (k Keeper) SetCrossChainTx(ctx sdk.Context, send types.CrossChainTx) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&send)
	store.Set(types.KeyPrefix(send.Index), b)
}

// GetCrossChainTx returns a send from its index
func (k Keeper) GetCrossChainTx(ctx sdk.Context, index string) (val types.CrossChainTx, found bool) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllCrossChainTx(ctx sdk.Context) (list []types.CrossChainTx) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// RemoveCrossChainTx removes a send from the store
func (k Keeper) RemoveCrossChainTx(ctx sdk.Context, index string) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	store.Delete(types.KeyPrefix(index))
}

// Queries

func (k Keeper) CreateNewCCTX(ctx sdk.Context, msg *types.MsgVoteOnObservedInboundTx, index string, tssPubkey string, s types.CctxStatus, senderChain, receiverChain *common.Chain) types.CrossChainTx {
	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &types.InboundTxParams{
		Sender:                          msg.Sender,
		SenderChainId:                   senderChain.ChainId,
		TxOrigin:                        msg.TxOrigin,
		Asset:                           msg.Asset,
		Amount:                          msg.Amount,
		CoinType:                        msg.CoinType,
		InboundTxObservedHash:           msg.InTxHash,
		InboundTxObservedExternalHeight: msg.InBlockHeight,
		InboundTxFinalizedZetaHeight:    0,
		InboundTxBallotIndex:            index,
	}

	outBoundParams := &types.OutboundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChainId:                  receiverChain.ChainId,
		OutboundTxHash:                   "",
		OutboundTxTssNonce:               0,
		OutboundTxGasLimit:               msg.GasLimit,
		OutboundTxGasPrice:               "",
		OutboundTxBallotIndex:            "",
		OutboundTxObservedExternalHeight: 0,
		CoinType:                         msg.CoinType, // FIXME: is this correct?
		Amount:                           sdk.NewUint(0),
		TssPubkey:                        tssPubkey,
	}
	status := &types.Status{
		Status:              s,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaFees:         math.ZeroUint(),
		RelayedMessage:   msg.Message,
		CctxStatus:       status,
		InboundTxParams:  inboundParams,
		OutboundTxParams: []*types.OutboundTxParams{outBoundParams},
	}
	return newCctx
}
