package v6

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// crosschainKeeper is an interface to prevent cyclic dependency
type crosschainKeeper interface {
	GetStoreKey() storetypes.StoreKey
	GetCodec() codec.Codec
	GetAllCrossChainTx(ctx sdk.Context) []types.CrossChainTx

	SetCrossChainTx(ctx sdk.Context, cctx types.CrossChainTx)
	AddFinalizedInbound(ctx sdk.Context, inboundTxHash string, senderChainID int64, height uint64)

	SetZetaAccounting(ctx sdk.Context, accounting types.ZetaAccounting)
}

// MigrateStore migrates the x/crosschain module state from the consensus version 4 to 5
// It resets the aborted zeta amount to use the inbound tx amount instead in situations where the outbound cctx is never created.
func MigrateStore(ctx sdk.Context, crosschainKeeper crosschainKeeper) error {
	tmpctx, commit := ctx.CacheContext()
	cctxListV14 := GetV14CCTX(tmpctx, crosschainKeeper)
	for _, cctx := range cctxListV14 {
		OutBoundParamsV15 := make([]*types.OutboundTxParams, len(cctx.OutboundTxParams))
		for j, outBoundParams := range cctx.OutboundTxParams {
			OutBoundParamsV15[j] = &types.OutboundTxParams{
				Receiver:                         outBoundParams.Receiver,
				ReceiverChainId:                  outBoundParams.ReceiverChainId,
				Amount:                           outBoundParams.Amount,
				OutboundTxTssNonce:               outBoundParams.OutboundTxTssNonce,
				OutboundTxGasLimit:               outBoundParams.OutboundTxGasLimit,
				OutboundTxGasPrice:               outBoundParams.OutboundTxGasPrice,
				OutboundTxHash:                   outBoundParams.OutboundTxHash,
				OutboundTxBallotIndex:            outBoundParams.OutboundTxBallotIndex,
				OutboundTxObservedExternalHeight: outBoundParams.OutboundTxObservedExternalHeight,
				OutboundTxGasUsed:                outBoundParams.OutboundTxGasUsed,
				OutboundTxEffectiveGasPrice:      outBoundParams.OutboundTxEffectiveGasPrice,
				OutboundTxEffectiveGasLimit:      outBoundParams.OutboundTxEffectiveGasLimit,
				TssPubkey:                        outBoundParams.TssPubkey,
				TxFinalizationStatus:             outBoundParams.TxFinalizationStatus,
			}
		}

		cctxV15 := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Sender:                          cctx.InboundTxParams.Sender,
				SenderChainId:                   cctx.InboundTxParams.SenderChainId,
				TxOrigin:                        cctx.InboundTxParams.TxOrigin,
				Asset:                           cctx.InboundTxParams.Asset,
				Amount:                          cctx.InboundTxParams.Amount,
				InboundTxObservedHash:           cctx.InboundTxParams.InboundTxObservedHash,
				InboundTxObservedExternalHeight: cctx.InboundTxParams.InboundTxObservedExternalHeight,
				InboundTxBallotIndex:            cctx.InboundTxParams.InboundTxBallotIndex,
				InboundTxFinalizedZetaHeight:    cctx.InboundTxParams.InboundTxFinalizedZetaHeight,
				TxFinalizationStatus:            cctx.InboundTxParams.TxFinalizationStatus,
			},
			Index:            cctx.Index,
			Creator:          cctx.Creator,
			OutboundTxParams: OutBoundParamsV15,
			CctxStatus:       cctx.CctxStatus,
			CoinType:         cctx.InboundTxParams.CoinType,
			ZetaFees:         cctx.ZetaFees,
			RelayedMessage:   cctx.RelayedMessage,
			EventIndex:       1, // We don't have this information in the old version
		}
		crosschainKeeper.SetCrossChainTx(tmpctx, cctxV15)
	}
	commit()
	return nil
}

func GetV14CCTX(ctx sdk.Context, crosschainKeeper crosschainKeeper) (list []types.CrossChainTxV14) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(crosschainKeeper.GetStoreKey()), p)

	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTxV14
		crosschainKeeper.GetCodec().MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}
