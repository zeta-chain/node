package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func EmitEventInboundFinalized(ctx sdk.Context, cctx *types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.InboundFinalized,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InBoundTxParams.Sender),
			sdk.NewAttribute(types.SenderChain, cctx.InBoundTxParams.SenderChain),
			sdk.NewAttribute(types.InTxHash, cctx.InBoundTxParams.InBoundTxObservedHash),
			sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", cctx.InBoundTxParams.InBoundTxObservedExternalHeight)),
			sdk.NewAttribute(types.Receiver, cctx.OutBoundTxParams.Receiver),
			sdk.NewAttribute(types.ReceiverChain, cctx.OutBoundTxParams.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, cctx.ZetaBurnt.String()),
			sdk.NewAttribute(types.ZetaMint, cctx.ZetaMint.String()),
			sdk.NewAttribute(types.RelayedMessage, cctx.RelayedMessage),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.Status.String()),
			sdk.NewAttribute(types.StatusMessage, cctx.CctxStatus.StatusMessage),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

func EmitEventCCTXCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.CctxCreated,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InBoundTxParams.Sender),
			sdk.NewAttribute(types.SenderChain, cctx.InBoundTxParams.SenderChain),
			sdk.NewAttribute(types.InTxHash, cctx.InBoundTxParams.InBoundTxObservedHash),
			sdk.NewAttribute(types.Receiver, cctx.OutBoundTxParams.Receiver),
			sdk.NewAttribute(types.ReceiverChain, cctx.OutBoundTxParams.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, cctx.ZetaBurnt.String()),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.String()),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

func EmitOutboundSuccess(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, oldStatus string, newStatus string, cctx *types.CrossChainTx) {
	event := sdk.NewEvent(types.OutboundTxSuccessful,
		sdk.NewAttribute(types.CctxIndex, cctx.Index),
		sdk.NewAttribute(types.OutTxHash, cctx.OutBoundTxParams.OutBoundTxHash),
		sdk.NewAttribute(types.ZetaMint, msg.ZetaMinted.String()),
		sdk.NewAttribute(types.OutTXVotingChain, msg.OutTxChain),
		sdk.NewAttribute(types.OldStatus, oldStatus),
		sdk.NewAttribute(types.NewStatus, newStatus),
		sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
	)
	ctx.EventManager().EmitEvent(event)
}

func EmitOutboundFailure(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, oldStatus string, newStatus string, cctx *types.CrossChainTx) {
	event := sdk.NewEvent(types.OutboundTxFailed,
		sdk.NewAttribute(types.CctxIndex, cctx.Index),
		sdk.NewAttribute(types.OutTxHash, cctx.OutBoundTxParams.OutBoundTxHash),
		sdk.NewAttribute(types.ZetaMint, msg.ZetaMinted.String()),
		sdk.NewAttribute(types.OutTXVotingChain, cctx.OutBoundTxParams.ReceiverChain),
		sdk.NewAttribute(types.OldStatus, oldStatus),
		sdk.NewAttribute(types.NewStatus, newStatus),
		sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
	)
	ctx.EventManager().EmitEvent(event)
}

func EmitCCTXScrubbed(ctx sdk.Context, cctx types.CrossChainTx, oldGasPrice, newGasPrice, chain string) {
	event := sdk.NewEvent(types.CctxScrubbed,
		sdk.NewAttribute(types.CctxIndex, cctx.Index),
		sdk.NewAttribute("OldGasPrice", oldGasPrice),
		sdk.NewAttribute("NewGasPrice", newGasPrice),
		sdk.NewAttribute("Chain", chain),
		sdk.NewAttribute("Nonce", fmt.Sprintf("%d", cctx.OutBoundTxParams.OutBoundTxTSSNonce)),
	)
	ctx.EventManager().EmitEvent(event)
}
