package keeper

import (
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SetRevertOutboundValues does the following things in one function:
// 1. create a new OutboundTxParams for the revert
// 2. append the new OutboundTxParams to the current OutboundTxParams
// 3. update the TxFinalizationStatus of the current OutboundTxParams to Executed.
func SetRevertOutboundValues(cctx *types.CrossChainTx, gasLimit uint64) {
	revertTxParams := &types.OutboundTxParams{
		Receiver:           cctx.InboundTxParams.Sender,
		ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
		Amount:             cctx.InboundTxParams.Amount,
		OutboundTxGasLimit: gasLimit,
		TssPubkey:          cctx.GetCurrentOutTxParam().TssPubkey,
	}
	// The original outbound has been finalized, the new outbound is pending
	cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	cctx.OutboundTxParams = append(cctx.OutboundTxParams, revertTxParams)
}

// SetOutboundValues sets the required values for the outbound transaction
// Note: It expects the cctx to already have been created,
// it updates the cctx based on the MsgVoteOnObservedOutboundTx message which is signed and broadcasted by the observer
func SetOutboundValues(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedOutboundTx, ballotStatus observertypes.BallotStatus) error {
	if ballotStatus != observertypes.BallotStatus_BallotFinalized_FailureObservation {
		if !msg.ValueReceived.Equal(cctx.GetCurrentOutTxParam().Amount) {
			ctx.Logger().Error(fmt.Sprintf("VoteOnObservedOutboundTx: Mint mismatch: %s value received vs %s cctx amount",
				msg.ValueReceived,
				cctx.GetCurrentOutTxParam().Amount))
			return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ValueReceived %s does not match sent value %s", msg.ValueReceived, cctx.GetCurrentOutTxParam().Amount))
		}
	}
	// Update CCTX values
	cctx.GetCurrentOutTxParam().OutboundTxHash = msg.ObservedOutTxHash
	cctx.GetCurrentOutTxParam().OutboundTxGasUsed = msg.ObservedOutTxGasUsed
	cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasPrice = msg.ObservedOutTxEffectiveGasPrice
	cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasLimit = msg.ObservedOutTxEffectiveGasLimit
	cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = msg.ObservedOutTxBlockHeight
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()

	return nil
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

// ProcessSuccessfulOutbound processes a successful outbound transaction. It does the following things in one function:
// 1. Change the status of the CCTX from PendingRevert to Reverted or from PendingOutbound to OutboundMined
// 2. Set the finalization status of the current outbound tx to executed
// 3. Emit an event for the successful outbound transaction
func (k Keeper) ProcessSuccessfulOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) {
	oldStatus := cctx.CctxStatus.Status
	switch oldStatus {
	case types.CctxStatus_PendingRevert:
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Reverted, "")
	case types.CctxStatus_PendingOutbound:
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "")
	default:
		return
	}
	cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundSuccess(ctx, valueReceived, oldStatus.String(), newStatus, *cctx)
}

// ProcessFailedOutbound processes a failed outbound transaction. It does the following things in one function:
// 1. For Admin Tx or a withdrawal from Zeta chain, it aborts the CCTX
// 2. For other CCTX, it creates a revert tx if the outbound tx is pending. If the status is pending revert, it aborts the CCTX
// 3. Emit an event for the failed outbound transaction
// 4. Set the finalization status of the current outbound tx to executed. If a revert tx is is created, the finalization status is not set, it would get set when the revert is processed via a subsequent transaction
func (k Keeper) ProcessFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) error {
	oldStatus := cctx.CctxStatus.Status
	if cctx.CoinType == common.CoinType_Cmd || common.IsZetaChain(cctx.InboundTxParams.SenderChainId) {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "")
	} else {
		switch oldStatus {
		case types.CctxStatus_PendingOutbound:

			gasLimit, err := k.GetRevertGasLimit(ctx, cctx)
			if err != nil {
				return cosmoserrors.Wrap(err, "GetRevertGasLimit")
			}
			if gasLimit == 0 {
				// use same gas limit of outbound as a fallback -- should not happen
				gasLimit = cctx.OutboundTxParams[0].OutboundTxGasLimit
			}

			// create new OutboundTxParams for the revert
			SetRevertOutboundValues(cctx, gasLimit)

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
			// Not setting the finalization status here, the required changes have been mad while creating the revert tx
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, "Outbound failed, start revert")
		case types.CctxStatus_PendingRevert:
			cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "Outbound failed: revert failed; abort TX")
		}
	}
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundFailure(ctx, valueReceived, oldStatus.String(), newStatus, *cctx)
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
	commit()
	return nil
}

// SaveFailedOutBound saves a failed outbound transaction.
// It does the following things in one function:
// 1. Change the status of the CCTX to Aborted
// 2. Save the outbound
func (k Keeper) SaveFailedOutBound(ctx sdk.Context, cctx *types.CrossChainTx, errMessage string, ballotIndex string) {
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, errMessage)
	ctx.Logger().Error(errMessage)

	k.SaveOutbound(ctx, cctx, ballotIndex)
}

// SaveSuccessfulOutBound saves a successful outbound transaction.
func (k Keeper) SaveSuccessfulOutBound(ctx sdk.Context, cctx *types.CrossChainTx, ballotIndex string) {
	k.SaveOutbound(ctx, cctx, ballotIndex)
}

// SaveOutbound saves the outbound transaction.It does the following things in one function:
// 1. Set the ballot index for the outbound vote to the cctx
// 2. Remove the nonce from the pending nonces
// 3. Remove the outbound tx tracker
// 4. Set the cctx and nonce to cctx and inTxHash to cctx
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
