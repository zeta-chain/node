package keeper

import (
	"context"
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerkeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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

	// Validate the message params to verify it against an existing cctx
	cctx, err := k.ValidateOutboundMessage(ctx, *msg)
	if err != nil {
		return nil, err
	}
	// get ballot index
	ballotIndex := msg.Digest()
	// vote on outbound ballot
	isFinalizingVote, isNew, ballot, observationChain, err := k.zetaObserverKeeper.VoteOnOutboundBallot(
		ctx,
		ballotIndex,
		msg.OutTxChain,
		msg.Status,
		msg.Creator)
	if err != nil {
		return nil, err
	}
	// if the ballot is new, set the index to the CCTX
	if isNew {
		observerkeeper.EmitEventBallotCreated(ctx, ballot, msg.ObservedOutTxHash, observationChain)
	}
	// if not finalized commit state here
	if !isFinalizingVote {
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}

	// if ballot successful, the value received should be the out tx amount
	err = cctx.AddOutbound(ctx, *msg, ballot.BallotStatus)
	if err != nil {
		return nil, err
	}
	// Fund the gas stability pool with the remaining funds
	k.FundStabilityPool(ctx, &cctx)

	err = k.ProcessOutbound(ctx, &cctx, ballot.BallotStatus, msg.ValueReceived.String())
	if err != nil {
		k.SaveFailedOutbound(ctx, &cctx, err.Error(), ballotIndex)
		return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
	}
	k.SaveSuccessfulOutbound(ctx, &cctx, ballotIndex)
	return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
}

// FundStabilityPool funds the stability pool with the remaining fees of an outbound tx
// The funds are sent to the gas stability pool associated with the receiver chain
func (k Keeper) FundStabilityPool(ctx sdk.Context, cctx *types.CrossChainTx) {
	// Fund the gas stability pool with the remaining funds
	if err := k.FundGasStabilityPoolFromRemainingFees(ctx, *cctx.GetCurrentOutTxParam(), cctx.GetCurrentOutTxParam().ReceiverChainId); err != nil {
		ctx.Logger().Error(fmt.Sprintf("VoteOnObservedOutboundTx: CCTX: %s Can't fund the gas stability pool with remaining fees %s", cctx.Index, err.Error()))
	}
}

// FundGasStabilityPoolFromRemainingFees funds the gas stability pool with the remaining fees of an outbound tx
func (k Keeper) FundGasStabilityPoolFromRemainingFees(ctx sdk.Context, outboundTxParams types.OutboundTxParams, chainID int64) error {
	gasUsed := outboundTxParams.OutboundTxGasUsed
	gasLimit := outboundTxParams.OutboundTxEffectiveGasLimit
	gasPrice := math.NewUintFromBigInt(outboundTxParams.OutboundTxEffectiveGasPrice.BigInt())

	if gasLimit == gasUsed {
		return nil
	}

	// We skip gas stability pool funding if one of the params is zero
	if gasLimit > 0 && gasUsed > 0 && !gasPrice.IsZero() {
		if gasLimit > gasUsed {
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

// percentOf returns the percentage of a number
func percentOf(n *big.Int, percent int64) *big.Int {
	n = n.Mul(n, big.NewInt(percent))
	n = n.Div(n, big.NewInt(100))
	return n
}

/* ProcessSuccessfulOutbound processes a successful outbound transaction. It does the following things in one function:

	1. Change the status of the CCTX from
	 - PendingRevert to Reverted
     - PendingOutbound to OutboundMined

	2. Set the finalization status of the current outbound tx to executed

	3. Emit an event for the successful outbound transaction
*/

func (k Keeper) ProcessSuccessfulOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) {
	oldStatus := cctx.CctxStatus.Status
	switch oldStatus {
	case types.CctxStatus_PendingRevert:
		cctx.SetReverted("Outbound succeeded, revert executed")
	case types.CctxStatus_PendingOutbound:
		cctx.SetOutBoundMined("Outbound succeeded, mined")
	default:
		return
	}
	cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundSuccess(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
}

/*
ProcessFailedOutbound processes a failed outbound transaction. It does the following things in one function:

 1. For Admin Tx or a withdrawal from Zeta chain, it aborts the CCTX

 2. For other CCTX
    - If the CCTX is in PendingOutbound, it creates a revert tx and sets the finalization status of the current outbound tx to executed
    - If the CCTX is in PendingRevert, it sets the Status to Aborted

 3. Emit an event for the failed outbound transaction

 4. Set the finalization status of the current outbound tx to executed. If a revert tx is is created, the finalization status is not set, it would get set when the revert is processed via a subsequent transaction
*/
func (k Keeper) ProcessFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) error {
	oldStatus := cctx.CctxStatus.Status
	if cctx.InboundTxParams.CoinType == common.CoinType_Cmd || common.IsZetaChain(cctx.InboundTxParams.SenderChainId) {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "")
	} else {
		switch oldStatus {
		case types.CctxStatus_PendingOutbound:

			gasLimit, err := k.GetRevertGasLimit(ctx, *cctx)
			if err != nil {
				return cosmoserrors.Wrap(err, "GetRevertGasLimit")
			}
			if gasLimit == 0 {
				// use same gas limit of outbound as a fallback -- should not happen
				gasLimit = cctx.OutboundTxParams[0].OutboundTxGasLimit
			}

			// create new OutboundTxParams for the revert
			err = cctx.AddRevertOutbound(gasLimit)
			if err != nil {
				return cosmoserrors.Wrap(err, "AddRevertOutbound")
			}

			err = k.PayGasAndUpdateCctx(
				ctx,
				cctx.InboundTxParams.SenderChainId,
				cctx,
				cctx.OutboundTxParams[0].Amount,
				false,
			)
			if err != nil {
				return err
			}
			err = k.UpdateNonce(ctx, cctx.InboundTxParams.SenderChainId, cctx)
			if err != nil {
				return err
			}
			// Not setting the finalization status here, the required changes have been made while creating the revert tx
			cctx.SetPendingRevert("Outbound failed, start revert")
		case types.CctxStatus_PendingRevert:
			cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
			cctx.SetAbort("Outbound failed: revert failed; abort TX")
		}
	}
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundFailure(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
	return nil
}

// ProcessOutbound processes the finalization of an outbound transaction based on the ballot status
// The state is committed only if the individual steps are successful
func (k Keeper) ProcessOutbound(ctx sdk.Context, cctx *types.CrossChainTx, ballotStatus observertypes.BallotStatus, valueReceived string) error {
	tmpCtx, commit := ctx.CacheContext()
	err := func() error {
		switch ballotStatus {
		case observertypes.BallotStatus_BallotFinalized_SuccessObservation:
			k.ProcessSuccessfulOutbound(tmpCtx, cctx, valueReceived)
		case observertypes.BallotStatus_BallotFinalized_FailureObservation:
			err := k.ProcessFailedOutbound(tmpCtx, cctx, valueReceived)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}
	err = cctx.Validate()
	if err != nil {
		return err
	}
	commit()
	return nil
}

/*
SaveFailedOutbound saves a failed outbound transaction.It does the following things in one function:

 1. Change the status of the CCTX to Aborted

 2. Save the outbound
*/
func (k Keeper) SaveFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, errMessage string, ballotIndex string) {
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, errMessage)
	ctx.Logger().Error(errMessage)

	k.SaveOutbound(ctx, cctx, ballotIndex)
}

// SaveSuccessfulOutbound saves a successful outbound transaction.
func (k Keeper) SaveSuccessfulOutbound(ctx sdk.Context, cctx *types.CrossChainTx, ballotIndex string) {
	k.SaveOutbound(ctx, cctx, ballotIndex)
}

/*
SaveOutbound saves the outbound transaction.It does the following things in one function:

 1. Set the ballot index for the outbound vote to the cctx

 2. Remove the nonce from the pending nonces

 3. Remove the outbound tx tracker

 4. Set the cctx and nonce to cctx and inTxHash to cctx
*/
func (k Keeper) SaveOutbound(ctx sdk.Context, cctx *types.CrossChainTx, ballotIndex string) {
	receiverChain := cctx.GetCurrentOutTxParam().ReceiverChainId
	tssPubkey := cctx.GetCurrentOutTxParam().TssPubkey
	outTxTssNonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce

	cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
	// #nosec G701 always in range
	k.GetObserverKeeper().RemoveFromPendingNonces(ctx, tssPubkey, receiverChain, int64(outTxTssNonce))
	k.RemoveOutTxTracker(ctx, receiverChain, outTxTssNonce)
	ctx.Logger().Info(fmt.Sprintf("Remove tracker %s: , Block Height : %d ", getOutTrackerIndex(receiverChain, outTxTssNonce), ctx.BlockHeight()))
	// This should set nonce to cctx only if a new revert is created.
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *cctx)
}

func (k Keeper) ValidateOutboundMessage(ctx sdk.Context, msg types.MsgVoteOnObservedOutboundTx) (types.CrossChainTx, error) {
	// check if CCTX exists and if the nonce matches
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxHash)
	if !found {
		return types.CrossChainTx{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("CCTX %s does not exist", msg.CctxHash))
	}
	if cctx.GetCurrentOutTxParam().OutboundTxTssNonce != msg.OutTxTssNonce {
		return types.CrossChainTx{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("OutTxTssNonce %d does not match CCTX OutTxTssNonce %d", msg.OutTxTssNonce, cctx.GetCurrentOutTxParam().OutboundTxTssNonce))
	}
	// do not process an outbound vote if TSS is not found
	_, found = k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return types.CrossChainTx{}, types.ErrCannotFindTSSKeys
	}
	if cctx.GetCurrentOutTxParam().ReceiverChainId != msg.OutTxChain {
		return types.CrossChainTx{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("OutTxChain %d does not match CCTX OutTxChain %d", msg.OutTxChain, cctx.GetCurrentOutTxParam().ReceiverChainId))
	}
	return cctx, nil
}
