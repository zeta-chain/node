package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Casts a vote on an outbound transaction observed on a connected chain (after
// it has been broadcasted to and finalized on a connected chain). If this is
// the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized. When a ballot is finalized, if the amount
// of zeta minted does not match the outbound transaction amount an error is
// thrown. If the amounts match, the outbound transaction hash and the "last
// updated" timestamp are updated.
//
// The transaction is proceeded to be finalized:
//
// If the observation was successful, the status is changed from "pending
// revert/outbound" to "reverted/mined". The difference between zeta burned
// and minted is minted by the bank module and deposited into the module
// account.
//
// If the observation was unsuccessful, and if the status is "pending outbound",
// prices and nonce are updated and the status is changed to "pending revert".
// If the status was "pending revert", the status is changed to "aborted".
//
// If there's an error in the finalization process, the CCTX status is set to
// 'aborted'.
//
// After finalization the outbound transaction tracker and pending nonces are
// removed, and the CCTX is updated in the store.
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteOnObservedOutboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedOutboundTx) (*types.MsgVoteOnObservedOutboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_OutBoundTx
	// Observer Chain already checked then inbound is created
	/* EDGE CASE : Params updated in during the finalization process
	   i.e Inbound has been finalized but outbound is still pending
	*/
	observationChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.OutTxChain)
	if observationChain == nil {
		return nil, zetaObserverTypes.ErrSupportedChains
	}
	err := zetaObserverTypes.CheckReceiveStatus(msg.Status)
	if err != nil {
		return nil, err
	}
	//Check is msg.Creator is authorized to vote
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain)
	if !ok {
		return nil, err
	}

	// Check if CCTX exists
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxHash)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("CCTX %s does not exist", msg.CctxHash))
	}

	ballotIndex := msg.Digest()
	// Add votes and Set Ballot
	ballot, isNew, err := k.GetBallot(ctx, ballotIndex, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	if isNew {
		EmitEventBallotCreated(ctx, ballot, msg.ObservedOutTxHash, observationChain.String())
		// Set this the first time when the ballot is created
		// The ballot might change if there are more votes in a different outbound ballot for this cctx hash
		// A new ballot is created everytime for every outtxhash for a cctx , but only the ballot with max votes gets finalized
		// If an observer adds their vote to a different ballot , it would be treated as NotYetVoted on the finalized one
		cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
		k.SetCrossChainTx(ctx, cctx)
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.ConvertReceiveStatusToVoteType(msg.Status))
	if err != nil {
		return nil, err
	}

	ballot, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
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
	// Set the ballot index to the finalized ballot
	cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
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
