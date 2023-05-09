package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) VoteOnObservedOutboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedOutboundTx) (*types.MsgVoteOnObservedOutboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_OutBoundTx
	// Observer Chain already checked then inbound is created
	/* EDGE CASE : Params updated in during the finalization process
	   i.e Inbound has been finalized but outbound is still pending
	*/
	observationChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.OutTxChain)
	err := zetaObserverTypes.CheckReceiveStatus(msg.Status)
	if err != nil {
		return nil, err
	}
	//Check is msg.Creator is authorized to vote
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain)
	if !ok {
		return nil, err
	}

	ballotIndex := msg.Digest()
	// Add votes and Set Ballot
	ballot, isNew, err := k.GetBallot(ctx, ballotIndex, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	if isNew {
		EmitEventBallotCreated(ctx, ballot, msg.ObservedOutTxHash, observationChain.String())
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.ConvertReceiveStatusToVoteType(msg.Status))
	if err != nil {
		return nil, err
	}
	// Check CCTX exists after confirmed vote
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxHash)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("CCTX %s does not exist", msg.CctxHash))
	}
	ballot, isFinalized := k.CheckIfBallotIsFinalized(ctx, ballot)
	if !isFinalized {
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}
	if ballot.BallotStatus != zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation {
		if !msg.ZetaMinted.Equal(cctx.GetCurrentOutTxParam().Amount) {
			log.Error().Msgf("ReceiveConfirmation: Mint mismatch: %s vs %s", msg.ZetaMinted, cctx.GetCurrentOutTxParam().Amount)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ZetaMinted %s does not match send ZetaMint %s", msg.ZetaMinted, cctx.GetCurrentOutTxParam().Amount))
		}
	}

	cctx.GetCurrentOutTxParam().OutboundTxHash = msg.ObservedOutTxHash
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()

	tss, _ := k.GetTSS(ctx)
	// FinalizeOutbound sets final status for a successful vote
	// FinalizeOutbound updates CCTX Prices and Nonce for a revert
	err = FinalizeOutbound(k, ctx, &cctx, msg, ballot.BallotStatus)
	if err != nil {
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
		ctx.Logger().Error(err.Error())
		k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
		k.RemoveFromPendingNonces(ctx, tss.TssPubkey, msg.OutTxChain, int64(msg.OutTxTssNonce))
		k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
		k.SetCrossChainTx(ctx, cctx)
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}

	k.RemoveFromPendingNonces(ctx, tss.TssPubkey, msg.OutTxChain, int64(msg.OutTxTssNonce))
	k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
}

func HandleFeeBalances(k msgServer, ctx sdk.Context, balanceAmount math.Uint) error {
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(common.ZETADenom, sdk.NewIntFromBigInt(balanceAmount.BigInt()))))
	if err != nil {
		log.Error().Msgf("ReceiveConfirmation: failed to mint coins: %s", err.Error())
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to mint coins: %s", err.Error()))
	}
	return nil
}

func FinalizeOutbound(k msgServer, ctx sdk.Context, cctx *types.CrossChainTx, msg *types.MsgVoteOnObservedOutboundTx, status zetaObserverTypes.BallotStatus) error {
	//cctx.GetCurrentOutTxParam().OutboundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = msg.ObservedOutTxBlockHeight
	zetaBurnt := cctx.InboundTxParams.Amount
	zetaMinted := cctx.GetCurrentOutTxParam().Amount
	oldStatus := cctx.CctxStatus.Status
	switch status {
	case zetaObserverTypes.BallotStatus_BallotFinalized_SuccessObservation:
		switch oldStatus {
		case types.CctxStatus_PendingRevert:
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Reverted, "Set To Final status", cctx.LogIdentifierForCCTX())
		case types.CctxStatus_PendingOutbound:
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_OutboundMined, "Set To Final status", cctx.LogIdentifierForCCTX())
		}

		newStatus := cctx.CctxStatus.Status.String()
		if zetaBurnt.LT(zetaMinted) {
			// TODO :Handle Error ?
		}
		balanceAmount := zetaBurnt.Sub(zetaMinted)
		if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta { // TODO : Handle Fee for other coins
			err := HandleFeeBalances(k, ctx, balanceAmount)
			if err != nil {
				return err
			}
		}
		EmitOutboundSuccess(ctx, msg, oldStatus.String(), newStatus, cctx)
	case zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation:
		switch oldStatus {
		case types.CctxStatus_PendingOutbound:
			// create new OutboundTxParams for the revert
			cctx.OutboundTxParams = append(cctx.OutboundTxParams, &types.OutboundTxParams{
				Receiver:           cctx.InboundTxParams.Sender,
				ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
				Amount:             cctx.InboundTxParams.Amount,
				CoinType:           cctx.InboundTxParams.CoinType,
				OutboundTxGasLimit: cctx.OutboundTxParams[0].OutboundTxGasLimit, // NOTE(pwu): revert gas limit = initial outbound gas limit set by user;
			})
			err := k.UpdatePrices(ctx, cctx.InboundTxParams.SenderChainId, cctx)
			if err != nil {
				return err
			}
			err = k.UpdateNonce(ctx, cctx.InboundTxParams.SenderChainId, cctx)
			if err != nil {
				return err
			}
			cctx.CctxStatus.ChangeStatus(&ctx,
				types.CctxStatus_PendingRevert, "Outbound failed, start revert", cctx.LogIdentifierForCCTX())
		case types.CctxStatus_PendingRevert:
			cctx.CctxStatus.ChangeStatus(&ctx,
				types.CctxStatus_Aborted, "Outbound failed: revert failed; abort TX", cctx.LogIdentifierForCCTX())

		}
		newStatus := cctx.CctxStatus.Status.String()
		EmitOutboundFailure(ctx, msg, oldStatus.String(), newStatus, cctx)
	}
	return nil
}
