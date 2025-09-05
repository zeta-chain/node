package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func EmitEventInboundFinalized(ctx sdk.Context, cctx *types.CrossChainTx) {
	currentOutParam := cctx.GetCurrentOutboundParam()
	err := ctx.EventManager().EmitTypedEvents(&types.EventInboundFinalized{
		MsgTypeUrl:         sdk.MsgTypeURL(&types.MsgVoteInbound{}),
		CctxIndex:          cctx.Index,
		Sender:             cctx.InboundParams.Sender,
		TxOrgin:            cctx.InboundParams.TxOrigin,
		Asset:              cctx.InboundParams.Asset,
		InboundHash:        cctx.InboundParams.ObservedHash,
		InboundBlockHeight: strconv.FormatUint(cctx.InboundParams.ObservedExternalHeight, 10),
		Receiver:           currentOutParam.Receiver,
		Amount:             cctx.InboundParams.Amount.String(),
		RelayedMessage:     cctx.RelayedMessage,
		NewStatus:          cctx.CctxStatus.Status.String(),
		StatusMessage:      cctx.CctxStatus.StatusMessage,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventInboundFinalized :", err)
	}
}

func EmitZRCWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventZrcWithdrawCreated{
		MsgTypeUrl:  "/zetachain.zetacore.crosschain.internal.ZRCWithdrawCreated",
		CctxIndex:   cctx.Index,
		Sender:      cctx.InboundParams.Sender,
		InboundHash: cctx.InboundParams.ObservedHash,
		NewStatus:   cctx.CctxStatus.Status.String(),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting ZRCWithdrawCreated :", err)
	}
}

func EmitZetaWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	err := ctx.EventManager().EmitTypedEvent(&types.EventZetaWithdrawCreated{
		MsgTypeUrl:  "/zetachain.zetacore.crosschain.internal.ZetaWithdrawCreated",
		CctxIndex:   cctx.Index,
		Sender:      cctx.InboundParams.Sender,
		InboundHash: cctx.InboundParams.ObservedHash,
		NewStatus:   cctx.CctxStatus.Status.String(),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting ZetaWithdrawCreated :", err)
	}
}

func EmitOutboundSuccess(ctx sdk.Context, valueReceived string, oldStatus string, newStatus string, cctxIndex string) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventOutboundSuccess{
		MsgTypeUrl:    sdk.MsgTypeURL(&types.MsgVoteOutbound{}),
		CctxIndex:     cctxIndex,
		ValueReceived: valueReceived,
		OldStatus:     oldStatus,
		NewStatus:     newStatus,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting MsgVoteOutbound :", err)
	}
}

func EmitOutboundFailure(ctx sdk.Context, valueReceived string, oldStatus string, newStatus string, cctxIndex string) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventOutboundFailure{
		MsgTypeUrl:    sdk.MsgTypeURL(&types.MsgVoteOutbound{}),
		CctxIndex:     cctxIndex,
		ValueReceived: valueReceived,
		OldStatus:     oldStatus,
		NewStatus:     newStatus,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting MsgVoteOutbound :", err)
	}
}

func EmitInboundProcessingFailure(ctx sdk.Context, inboundHash string, errorMessage string) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventInboundProcessingFailure{
		InboundHash:  inboundHash,
		ErrorMessage: errorMessage,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventInboundProcessingFailure :", err)
	}
}
