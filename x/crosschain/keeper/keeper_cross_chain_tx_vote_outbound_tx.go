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
	observerKeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Casts a vote on an outbound transaction observed on a connected chain (after
// it has been broadcasted to and finalized on a connected chain). If this is
// the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized. When a ballot is finalized, the outbound
// transaction is processed.
//
// If the observation was successful, the difference between zeta burned
// and minted is minted by the bank module and deposited into the module
// account.
//
// If the observation was unsuccessful, the logic depends on the previous
// status.
//
// If the previous status was `PendingOutbound`, a new revert transaction is
// created. To cover the revert transaction fee, the required amount of tokens
// submitted with the CCTX are swapped using a Uniswap V2 contract instance on
// ZetaChain for the ZRC20 of the gas token of the receiver chain. The ZRC20
// tokens are then
// burned. The nonce is updated. If everything is successful, the CCTX status is
// changed to `PendingRevert`.
//
// If the previous status was `PendingRevert`, the CCTX is aborted.
//
// ```mermaid
// stateDiagram-v2
//
//	state finalize_outbound <<choice>>
//	state observation <<choice>>
//	state success_old_status <<choice>>
//	state fail_old_status <<choice>>
//	[*] --> finalize_outbound
//	finalize_outbound --> observation: Finalize outbound
//	observation --> success_old_status: Observation succeeded
//	success_old_status --> Reverted: Old status is PendingRevert
//	success_old_status --> OutboundMined: Old status is PendingOutbound
//	observation --> fail_old_status: Observation failed
//	fail_old_status --> PendingRevert: Old status is PendingOutbound
//	fail_old_status --> Aborted: Old status is PendingRevert
//	finalize_outbound --> Aborted: Finalize outbound error
//	Aborted --> [*]
//	Reverted --> [*]
//	OutboundMined --> [*]
//	PendingRevert --> [*]
//
// ```
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteOnObservedOutboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedOutboundTx) (*types.MsgVoteOnObservedOutboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := observerTypes.ObservationType_OutBoundTx
	// Observer Chain already checked then inbound is created
	/* EDGE CASE : Params updated in during the finalization process
	   i.e Inbound has been finalized but outbound is still pending
	*/
	observationChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.OutTxChain)
	if observationChain == nil {
		return nil, observerTypes.ErrSupportedChains
	}
	err := observerTypes.CheckReceiveStatus(msg.Status)
	if err != nil {
		return nil, err
	}
	//Check is msg.Creator is authorized to vote
	ok, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, observationChain)
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
	ballot, isNew, err := k.zetaObserverKeeper.FindBallot(ctx, ballotIndex, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	if isNew {
		observerKeeper.EmitEventBallotCreated(ctx, ballot, msg.ObservedOutTxHash, observationChain.String())
		// Set this the first time when the ballot is created
		// The ballot might change if there are more votes in a different outbound ballot for this cctx hash
		cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
		//k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.zetaObserverKeeper.AddVoteToBallot(ctx, ballot, msg.Creator, observerTypes.ConvertReceiveStatusToVoteType(msg.Status))
	if err != nil {
		return nil, err
	}

	ballot, isFinalized := k.zetaObserverKeeper.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}
	if ballot.BallotStatus != observerTypes.BallotStatus_BallotFinalized_FailureObservation {
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

	err = func() error { //err = FinalizeOutbound(k, ctx, &cctx, msg, ballot.BallotStatus) // remove this line
		cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = msg.ObservedOutTxBlockHeight
		zetaBurnt := cctx.InboundTxParams.Amount
		zetaMinted := cctx.GetCurrentOutTxParam().Amount
		oldStatus := cctx.CctxStatus.Status
		switch ballot.BallotStatus {
		case observerTypes.BallotStatus_BallotFinalized_SuccessObservation:
			switch oldStatus {
			case types.CctxStatus_PendingRevert:
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Reverted, "", cctx.LogIdentifierForCCTX())
			case types.CctxStatus_PendingOutbound:
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_OutboundMined, "", cctx.LogIdentifierForCCTX())
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
		case observerTypes.BallotStatus_BallotFinalized_FailureObservation:
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
				err := k.UpdatePrices(ctx, cctx.InboundTxParams.SenderChainId, &cctx)
				if err != nil {
					return err
				}
				err = k.UpdateNonce(ctx, cctx.InboundTxParams.SenderChainId, &cctx)
				if err != nil {
					return err
				}
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingRevert, "Outbound failed, start revert", cctx.LogIdentifierForCCTX())
			case types.CctxStatus_PendingRevert:
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, "Outbound failed: revert failed; abort TX", cctx.LogIdentifierForCCTX())
			}
			newStatus := cctx.CctxStatus.Status.String()
			EmitOutboundFailure(ctx, msg, oldStatus.String(), newStatus, cctx)
		}
		return nil
	}()
	if err != nil {
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
		ctx.Logger().Error(err.Error())
		k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
		k.RemoveFromPendingNonces(ctx, tss.TssPubkey, msg.OutTxChain, int64(msg.OutTxTssNonce))
		k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
		k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}
	// Set the ballot index to the finalized ballot
	cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
	k.RemoveFromPendingNonces(ctx, tss.TssPubkey, msg.OutTxChain, int64(msg.OutTxTssNonce))
	k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
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
