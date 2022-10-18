package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"strconv"
)

func (k msgServer) VoteOnObservedOutboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedOutboundTx) (*types.MsgVoteOnObservedOutboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_OutBoundTx
	observationChain := zetaObserverTypes.ConvertStringChaintoObservationChain(msg.OutTxChain)
	err := zetaObserverTypes.CheckReceiveStatus(msg.Status)
	if err != nil {
		return nil, err
	}
	//Check is msg.Creator is authorized to vote
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain, observationType.String())
	if !ok {
		return nil, err
	}

	ballotIndex := msg.Digest()
	// Add votes and Set Ballot
	ballot, err := k.GetBallot(ctx, ballotIndex, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.ConvertReceiveStatusToVoteType(msg.Status))
	if err != nil {
		return nil, err
	}
	// Check CCTX exists after confirmed vote
	cctx, err := k.CheckCCTXExists(ctx, ballot.Index, msg.CctxHash)
	if err != nil {
		return nil, err
	}
	ballot, isFinalized := k.CheckIfBallotIsFinalized(ctx, ballot)
	if !isFinalized {
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}

	if ballot.BallotStatus != zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation {
		if !msg.ZetaMinted.Equal(cctx.ZetaMint) {
			log.Error().Msgf("ReceiveConfirmation: Mint mismatch: %s vs %s", msg.ZetaMinted, cctx.ZetaMint)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ZetaMinted %s does not match send ZetaMint %s", msg.ZetaMinted, cctx.ZetaMint))
		}
	}
	cctx.OutBoundTxParams.OutBoundTxHash = msg.ObservedOutTxHash
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()

	oldStatus := cctx.CctxStatus.Status
	// FinalizeOutbound sets final status for a successful vote
	// FinalizeOutbound updates CCTX Prices and Nonce for a revert
	err = FinalizeOutbound(k, ctx, &cctx, msg, ballot.BallotStatus)
	if err != nil {
		return nil, err
	}
	// Remove OutTX tracker and change CCTX prefix store
	k.RemoveOutTxTracker(ctx, fmt.Sprintf("%s-%s", msg.OutTxChain, strconv.Itoa(int(msg.OutTxTssNonce))))
	k.CctxChangePrefixStore(ctx, cctx, oldStatus)

	return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
}

func HandleFeeBalances(k msgServer, ctx sdk.Context, balanceAmount sdk.Uint) error {
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(common.ZETADenom, sdk.NewIntFromBigInt(balanceAmount.BigInt()))))
	if err != nil {
		log.Error().Msgf("ReceiveConfirmation: failed to mint coins: %s", err.Error())
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to mint coins: %s", err.Error()))
	}
	return nil
}

func FinalizeOutbound(k msgServer, ctx sdk.Context, cctx *types.CrossChainTx, msg *types.MsgVoteOnObservedOutboundTx, status zetaObserverTypes.BallotStatus) error {
	cctx.OutBoundTxParams.OutBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	cctx.OutBoundTxParams.OutBoundTxObservedExternalHeight = msg.ObservedOutTxBlockHeight
	zetaBurnt := cctx.ZetaBurnt
	zetaMinted := cctx.ZetaMint
	oldStatus := cctx.CctxStatus.Status
	switch status {
	case zetaObserverTypes.BallotStatus_BallotFinalized_SuccessObservation:
		switch oldStatus {
		case types.CctxStatus_PendingRevert:
			cctx.CctxStatus.ChangeStatus(&ctx,
				types.CctxStatus_Reverted, "Set To Final status", cctx.LogIdentifierForCCTX())
		case types.CctxStatus_PendingOutbound:
			cctx.CctxStatus.ChangeStatus(&ctx,
				types.CctxStatus_OutboundMined, "Set To Final status", cctx.LogIdentifierForCCTX())
		}

		newStatus := cctx.CctxStatus.Status.String()
		if zetaBurnt.LT(zetaMinted) {
			// TODO :Handle Error ?
		}
		balanceAmount := zetaBurnt.Sub(zetaMinted)
		err := HandleFeeBalances(k, ctx, balanceAmount)
		if err != nil {
			return err
		}
		EmitOutboundSuccess(ctx, msg, oldStatus.String(), newStatus, cctx)
	case zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation:
		switch oldStatus {
		case types.CctxStatus_PendingOutbound:
			chain := cctx.InBoundTxParams.SenderChain
			err := k.UpdatePrices(ctx, chain, cctx)
			if err != nil {
				return err
			}
			err = k.UpdateNonce(ctx, chain, cctx)
			if err != nil {
				return err
			}
			cctx.CctxStatus.ChangeStatus(&ctx,
				types.CctxStatus_PendingRevert, "Outbound Failed , Starting Revert", cctx.LogIdentifierForCCTX())
		case types.CctxStatus_PendingRevert:
			cctx.CctxStatus.ChangeStatus(&ctx,
				types.CctxStatus_Aborted, "Outbound Failed & Revert Failed , Abort TX", cctx.LogIdentifierForCCTX())

		}
		newStatus := cctx.CctxStatus.Status.String()
		EmitOutboundFailure(ctx, msg, oldStatus.String(), newStatus, cctx)
	}
	return nil
}
