package keeper

import (
	types2 "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/zeta-chain/zetacore/common"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func EmitEventInboundFinalized(ctx sdk.Context, cctx *types.CrossChainTx) error {
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
		return err
	}
	return nil
}

func EmitEventKeyGenBlockUpdated(ctx sdk.Context, keygen *types.Keygen) error {
	err := ctx.EventManager().EmitTypedEvents(&types.EventKeygenBlockUpdated{
		MsgTypeUrl:    sdk.MsgTypeURL(&types.MsgUpdateKeygen{}),
		KeygenBlock:   strconv.Itoa(int(keygen.BlockNumber)),
		KeygenPubkeys: types2.PrettyPrintStruct(keygen.GranteePubkeys),
	})
	if err != nil {
		return err
	}
	return nil
}

func EmitZRCWithdrawCreated(ctx sdk.Context, cctx types.CrossChainTx) error {
	err := ctx.EventManager().EmitTypedEvents(&types.EventZrcWithdrawCreated{
		MsgTypeUrl: " /zetachain.zetacore.crosschain. ZRCWithdrawCreated",
		CctxIndex:  cctx.Index,
		Sender:     cctx.InboundTxParams.Sender,
		InTxHash:   cctx.InboundTxParams.InboundTxObservedHash,
		NewStatus:  cctx.CctxStatus.Status.String(),
	})
	if err != nil {
		return err
	}
	return nil
}

func EmitEventBallotCreated(ctx sdk.Context, ballot zetaObserverTypes.Ballot, observationHash, obserVationChain string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.BallotCreated,
			sdk.NewAttribute(types.BallotIdentifier, ballot.BallotIdentifier),
			sdk.NewAttribute(types.CCTXIndex, ballot.BallotIdentifier),
			sdk.NewAttribute(types.BallotObservationHash, observationHash),
			sdk.NewAttribute(types.BallotObservationChain, obserVationChain),
			sdk.NewAttribute(types.BallotType, ballot.ObservationType.String()),
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
