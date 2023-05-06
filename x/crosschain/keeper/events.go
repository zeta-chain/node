package keeper

import (
	"fmt"
	types2 "github.com/coinbase/rosetta-sdk-go/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func EmitEventInboundFinalized(ctx sdk.Context, cctx *types.CrossChainTx) {
	currentOutParam := cctx.GetCurrentOutTxParam()
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.InboundFinalized,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InboundTxParams.Sender),
			sdk.NewAttribute(types.SenderChain, fmt.Sprintf("%d", cctx.InboundTxParams.SenderChainId)),
			sdk.NewAttribute(types.TxOrigin, cctx.InboundTxParams.TxOrigin),
			sdk.NewAttribute(types.Asset, cctx.InboundTxParams.Asset),
			sdk.NewAttribute(types.InTxHash, cctx.InboundTxParams.InboundTxObservedHash),
			sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", cctx.InboundTxParams.InboundTxObservedExternalHeight)),
			sdk.NewAttribute(types.Receiver, currentOutParam.Receiver),
			sdk.NewAttribute(types.ReceiverChain, fmt.Sprintf("%d", currentOutParam.ReceiverChainId)),
			sdk.NewAttribute(types.Amount, cctx.InboundTxParams.Amount.String()),
			//sdk.NewAttribute(types.ZetaMint, cctx.ZetaMint.String()),
			sdk.NewAttribute(types.RelayedMessage, cctx.RelayedMessage),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.Status.String()),
			sdk.NewAttribute(types.StatusMessage, cctx.CctxStatus.StatusMessage),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

func EmitEventKeyGenBlockUpdated(ctx sdk.Context, keygen *types.Keygen) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.CctxNewKeygenBlock,
			sdk.NewAttribute(types.KeyGenBlock, strconv.Itoa(int(keygen.BlockNumber))),
			sdk.NewAttribute(types.KeyGenPubKeys, types2.PrettyPrintStruct(keygen.GranteePubkeys)),
		),
	)
}

func EmitZRCWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.ZrcWithdrawCreated,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InboundTxParams.Sender),
			//sdk.NewAttribute(types.SenderChain, cctx.InboundTxParams.SenderChain),
			sdk.NewAttribute(types.InTxHash, cctx.InboundTxParams.InboundTxObservedHash),
			//sdk.NewAttribute(types.Receiver, cctx.OutboundTxParams.Receiver),
			//sdk.NewAttribute(types.ReceiverChain, cctx.OutboundTxParams.ReceiverChain),
			//sdk.NewAttribute(types.Amount, cctx.ZetaBurnt.String()),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.Status.String()),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

func EmitEventBallotCreated(ctx sdk.Context, ballot zetaObserverTypes.Ballot, observationHash string, obserVationChain string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.BallotCreated,
			sdk.NewAttribute(types.BallotIdentifier, ballot.BallotIdentifier),
			sdk.NewAttribute(types.CCTXIndex, ballot.BallotIdentifier),
			sdk.NewAttribute(types.BallotObservationHash, observationHash),
			sdk.NewAttribute(types.BallotObservationChain, obserVationChain),
		),
	)
}

func EmitZetaWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.ZetaWithdrawCreated,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InboundTxParams.Sender),
			//sdk.NewAttribute(types.SenderChain, cctx.InboundTxParams.SenderChain),
			sdk.NewAttribute(types.InTxHash, cctx.InboundTxParams.InboundTxObservedHash),
			//sdk.NewAttribute(types.Receiver, cctx.OutboundTxParams.Receiver),
			//sdk.NewAttribute(types.ReceiverChain, cctx.OutboundTxParams.ReceiverChain),
			//sdk.NewAttribute(types.Amount, cctx.ZetaBurnt.String()),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.Status.String()),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

func EmitOutboundSuccess(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, oldStatus string, newStatus string, cctx *types.CrossChainTx) {
	event := sdk.NewEvent(types.OutboundTxSuccessful,
		sdk.NewAttribute(types.CctxIndex, cctx.Index),
		//sdk.NewAttribute(types.OutTxHash, cctx.OutboundTxParams.OutboundTxHash),
		sdk.NewAttribute(types.ZetaMint, msg.ZetaMinted.String()),
		//sdk.NewAttribute(types.OutTXVotingChain, cctx.OutboundTxParams.ReceiverChain),
		sdk.NewAttribute(types.OldStatus, oldStatus),
		sdk.NewAttribute(types.NewStatus, newStatus),
		sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
	)
	ctx.EventManager().EmitEvent(event)
}

func EmitOutboundFailure(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, oldStatus string, newStatus string, cctx *types.CrossChainTx) {
	event := sdk.NewEvent(types.OutboundTxFailed,
		sdk.NewAttribute(types.CctxIndex, cctx.Index),
		//sdk.NewAttribute(types.OutTxHash, cctx.OutboundTxParams.OutboundTxHash),
		sdk.NewAttribute(types.ZetaMint, msg.ZetaMinted.String()),
		//sdk.NewAttribute(types.OutTXVotingChain, cctx.OutboundTxParams.ReceiverChain),
		sdk.NewAttribute(types.OldStatus, oldStatus),
		sdk.NewAttribute(types.NewStatus, newStatus),
		sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
	)
	ctx.EventManager().EmitEvent(event)
}

func EmitCCTXScrubbed(ctx sdk.Context, cctx types.CrossChainTx, chainID int64, oldGasPrice, newGasPrice string) {
	event := sdk.NewEvent(types.CctxScrubbed,
		sdk.NewAttribute(types.CctxIndex, cctx.Index),
		sdk.NewAttribute("OldGasPrice", oldGasPrice),
		sdk.NewAttribute("NewGasPrice", newGasPrice),
		sdk.NewAttribute("Chain ID", strconv.FormatInt(chainID, 10)),
		//sdk.NewAttribute("Nonce", fmt.Sprintf("%d", cctx.OutboundTxParams.OutboundTxTssNonce)),
	)
	ctx.EventManager().EmitEvent(event)
}
