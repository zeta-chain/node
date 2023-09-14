package keeper

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerKeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// VoteOnObservedOutboundTx casts a vote on an outbound transaction observed on a connected chain (after
// it has been broadcasted to and finalized on a connected chain). If this is
// the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized. When a ballot is finalized, the outbound
// transaction is processed.
//
// If the observation is successful, the difference between zeta burned
// and minted is minted by the bank module and deposited into the module
// account.
//
// If the observation is unsuccessful, the logic depends on the previous
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
//	state observation <<choice>>
//	state success_old_status <<choice>>
//	state fail_old_status <<choice>>
//	PendingOutbound --> observation: Finalize outbound
//	observation --> success_old_status: Observation succeeded
//	success_old_status --> Reverted: Old status is PendingRevert
//	success_old_status --> OutboundMined: Old status is PendingOutbound
//	observation --> fail_old_status: Observation failed
//	fail_old_status --> PendingRevert: Old status is PendingOutbound
//	fail_old_status --> Aborted: Old status is PendingRevert
//	PendingOutbound --> Aborted: Finalize outbound error
//
// ```
//
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
			log.Error().Msgf("VoteOnObservedOutboundTx: Mint mismatch: %s zeta minted vs %s cctx amount",
				msg.ZetaMinted,
				cctx.GetCurrentOutTxParam().Amount)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ZetaMinted %s does not match send ZetaMint %s", msg.ZetaMinted, cctx.GetCurrentOutTxParam().Amount))
		}
	}

	// Update CCTX values
	cctx.GetCurrentOutTxParam().OutboundTxHash = msg.ObservedOutTxHash
	cctx.GetCurrentOutTxParam().OutboundTxGasUsed = msg.ObservedOutTxGasUsed
	cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasPrice = msg.ObservedOutTxEffectiveGasPrice
	cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasLimit = msg.ObservedOutTxEffectiveGasLimit
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()

	// Fund the gas stability pool with the remaining funds
	if err := k.FundGasStabilityPoolFromRemainingFees(ctx, *cctx.GetCurrentOutTxParam(), msg.OutTxChain); err != nil {
		log.Error().Msgf(
			"VoteOnObservedOutboundTx: CCTX: %s Can't fund the gas stability pool with remaining fees %s", cctx.Index, err.Error(),
		)
	}

	tss, _ := k.GetTSS(ctx)

	// FinalizeOutbound sets final status for a successful vote
	// FinalizeOutbound updates CCTX Prices and Nonce for a revert

	tmpCtx, commit := ctx.CacheContext()
	err = func() error { //err = FinalizeOutbound(k, ctx, &cctx, msg, ballot.BallotStatus)
		cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = msg.ObservedOutTxBlockHeight
		oldStatus := cctx.CctxStatus.Status
		switch ballot.BallotStatus {
		case observerTypes.BallotStatus_BallotFinalized_SuccessObservation:
			switch oldStatus {
			case types.CctxStatus_PendingRevert:
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Reverted, "")
			case types.CctxStatus_PendingOutbound:
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "")
			}
			newStatus := cctx.CctxStatus.Status.String()
			EmitOutboundSuccess(tmpCtx, msg, oldStatus.String(), newStatus, cctx)
		case observerTypes.BallotStatus_BallotFinalized_FailureObservation:
			if msg.CoinType == common.CoinType_Cmd || cctx.InboundTxParams.SenderChainId == common.ZetaChain().ChainId {
				// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "")
			} else {
				switch oldStatus {
				case types.CctxStatus_PendingOutbound:
					// create new OutboundTxParams for the revert
					cctx.OutboundTxParams = append(cctx.OutboundTxParams, &types.OutboundTxParams{
						Receiver:        cctx.InboundTxParams.Sender,
						ReceiverChainId: cctx.InboundTxParams.SenderChainId,
						Amount:          cctx.InboundTxParams.Amount,
						CoinType:        cctx.InboundTxParams.CoinType,
						// NOTE(pwu): revert gas limit = initial outbound gas limit set by user
						//TODO: determine a specific revert gas limit https://github.com/zeta-chain/node/issues/1065
						OutboundTxGasLimit: cctx.OutboundTxParams[0].OutboundTxGasLimit,
					})
					err := k.PayGasAndUpdateCctx(
						tmpCtx,
						cctx.InboundTxParams.SenderChainId,
						&cctx,
						cctx.OutboundTxParams[0].Amount,
						false,
					)
					if err != nil {
						return err
					}
					err = k.UpdateNonce(tmpCtx, cctx.InboundTxParams.SenderChainId, &cctx)
					if err != nil {
						return err
					}
					cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, "Outbound failed, start revert")
				case types.CctxStatus_PendingRevert:
					cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "Outbound failed: revert failed; abort TX")
				}
			}
			newStatus := cctx.CctxStatus.Status.String()
			EmitOutboundFailure(ctx, msg, oldStatus.String(), newStatus, cctx)
		}
		return nil
	}()
	if err != nil {
		// do not commit tmpCtx
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		ctx.Logger().Error(err.Error())
		k.RemoveFromPendingNonces(ctx, tss.TssPubkey, msg.OutTxChain, int64(msg.OutTxTssNonce))
		k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
		k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}
	commit()
	// Set the ballot index to the finalized ballot
	cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
	k.RemoveFromPendingNonces(ctx, tss.TssPubkey, msg.OutTxChain, int64(msg.OutTxTssNonce))
	k.RemoveOutTxTracker(ctx, msg.OutTxChain, msg.OutTxTssNonce)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
}

func percentOf(n *big.Int, percent int64) *big.Int {
	n = n.Mul(n, big.NewInt(percent))
	n = n.Div(n, big.NewInt(100))
	return n
}

// FundGasStabilityPoolFromRemainingFees funds the gas stability pool with the remaining fees of an outbound tx
func (k Keeper) FundGasStabilityPoolFromRemainingFees(ctx sdk.Context, outboundTxParams types.OutboundTxParams, chainID int64) error {
	gasUsed := outboundTxParams.OutboundTxGasUsed
	gasLimit := outboundTxParams.OutboundTxEffectiveGasLimit
	gasPrice := math.NewUintFromBigInt(outboundTxParams.OutboundTxEffectiveGasPrice.BigInt())

	// We skip gas stability pool funding if one of the params is zero
	if gasLimit > 0 && gasUsed > 0 && !gasPrice.IsZero() {
		if gasLimit >= gasUsed {
			remainingGas := gasLimit - gasUsed
			remainingFees := math.NewUint(remainingGas).Mul(gasPrice).BigInt()

			// We fund the stability pool with a portion of the remaining fees
			remainingFees = percentOf(remainingFees, RemainingFeesToStabilityPoolPercent)

			// Fund the gas stability pool
			if err := k.fungibleKeeper.FundGasStabilityPool(ctx, chainID, remainingFees); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("VoteOnObservedOutboundTx: The gas limit %d is less than the gas used %d", gasLimit, gasUsed)
		}
	}
	return nil
}
