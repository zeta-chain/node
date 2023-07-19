package keeper

import (
	"github.com/zeta-chain/zetacore/common"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func EmitEventInboundFinalized(ctx sdk.Context, cctx *types.CrossChainTx) {
	currentOutParam := cctx.GetCurrentOutTxParam()
	err := ctx.EventManager().EmitTypedEvents(&types.EventInboundFinalized{
		MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgVoteOnObservedInboundTx{}),
		CctxIndex:      cctx.Index,
		Sender:         cctx.InboundTxParams.Sender,
		SenderChain:    common.GetChainFromChainID(cctx.InboundTxParams.SenderChainId).ChainName.String(),
		TxOrgin:        cctx.InboundTxParams.TxOrigin,
		Asset:          cctx.InboundTxParams.Asset,
		InTxHash:       cctx.InboundTxParams.InboundTxObservedHash,
		InBlockHeight:  strconv.FormatUint(cctx.InboundTxParams.InboundTxObservedExternalHeight, 10),
		Receiver:       currentOutParam.Receiver,
		ReceiverChain:  common.GetChainFromChainID(currentOutParam.ReceiverChainId).ChainName.String(),
		Amount:         cctx.InboundTxParams.Amount.String(),
		RelayedMessage: cctx.RelayedMessage,
		NewStatus:      cctx.CctxStatus.Status.String(),
		StatusMessage:  cctx.CctxStatus.StatusMessage,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventInboundFinalized :", err)
	}
}

func EmitZRCWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventZrcWithdrawCreated{
		MsgTypeUrl: "/zetachain.zetacore.crosschain.internal.ZRCWithdrawCreated",
		CctxIndex:  cctx.Index,
		Sender:     cctx.InboundTxParams.Sender,
		InTxHash:   cctx.InboundTxParams.InboundTxObservedHash,
		NewStatus:  cctx.CctxStatus.Status.String(),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting ZRCWithdrawCreated :", err)
	}
}

func EmitZetaWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	err := ctx.EventManager().EmitTypedEvent(&types.EventZetaWithdrawCreated{
		MsgTypeUrl: "/zetachain.zetacore.crosschain.internal.ZetaWithdrawCreated",
		CctxIndex:  cctx.Index,
		Sender:     cctx.InboundTxParams.Sender,
		InTxHash:   cctx.InboundTxParams.InboundTxObservedHash,
		NewStatus:  cctx.CctxStatus.Status.String(),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting ZetaWithdrawCreated :", err)
	}

}

func EmitOutboundSuccess(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, oldStatus string, newStatus string, cctx *types.CrossChainTx) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventOutboundSuccess{
		MsgTypeUrl: sdk.MsgTypeURL(&types.MsgVoteOnObservedOutboundTx{}),
		CctxIndex:  cctx.Index,
		ZetaMinted: msg.ZetaMinted.String(),
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting MsgVoteOnObservedOutboundTx :", err)
	}

}

func EmitOutboundFailure(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, oldStatus string, newStatus string, cctx *types.CrossChainTx) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventOutboundFailure{
		MsgTypeUrl: sdk.MsgTypeURL(&types.MsgVoteOnObservedOutboundTx{}),
		CctxIndex:  cctx.Index,
		ZetaMinted: msg.ZetaMinted.String(),
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting MsgVoteOnObservedOutboundTx :", err)
	}
}
