package keeper

import (
	"errors"
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

func (k Keeper) GetOutbound(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedOutboundTx, ballotStatus observertypes.BallotStatus) error {
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

	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return types.ErrCannotFindTSSKeys
	}
	if tss.TssPubkey != cctx.GetCurrentOutTxParam().TssPubkey {
		return types.ErrTssMismatch
	}
	return nil
}

func (k Keeper) FundStabiltityPool(ctx sdk.Context, cctx *types.CrossChainTx) {
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

func percentOf(n *big.Int, percent int64) *big.Int {
	n = n.Mul(n, big.NewInt(percent))
	n = n.Div(n, big.NewInt(100))
	return n
}

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
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundSuccess(ctx, valueReceived, oldStatus.String(), newStatus, *cctx)
}

func (k Keeper) ProcessFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) error {
	oldStatus := cctx.CctxStatus.Status
	if cctx.CoinType == common.CoinType_Cmd || common.IsZetaChain(cctx.InboundTxParams.SenderChainId) {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "")
	} else {
		switch oldStatus {
		case types.CctxStatus_PendingOutbound:

			gasLimit, err := k.GetRevertGasLimit(ctx, cctx)
			if err != nil {
				return errors.New("can't get revert tx gas limit" + err.Error())
			}
			if gasLimit == 0 {
				// use same gas limit of outbound as a fallback -- should not happen
				gasLimit = cctx.OutboundTxParams[0].OutboundTxGasLimit
			}

			// create new OutboundTxParams for the revert
			revertTxParams := &types.OutboundTxParams{
				Receiver:           cctx.InboundTxParams.Sender,
				ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
				Amount:             cctx.InboundTxParams.Amount,
				OutboundTxGasLimit: gasLimit,
			}
			cctx.OutboundTxParams = append(cctx.OutboundTxParams, revertTxParams)

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
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, "Outbound failed, start revert")
		case types.CctxStatus_PendingRevert:
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "Outbound failed: revert failed; abort TX")
		}
	}
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundFailure(ctx, valueReceived, oldStatus.String(), newStatus, *cctx)
	return nil
}
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

func (k Keeper) SaveFailedOutBound(ctx sdk.Context, cctx *types.CrossChainTx, errMessage string) {
	receiverChain := cctx.GetCurrentOutTxParam().ReceiverChainId
	tssPubkey := cctx.GetCurrentOutTxParam().TssPubkey
	outTxTssNonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce

	cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, errMessage)
	cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	ctx.Logger().Error(errMessage)
	// #nosec G701 always in range
	k.GetObserverKeeper().RemoveFromPendingNonces(ctx, tssPubkey, receiverChain, int64(outTxTssNonce))
	k.RemoveOutTxTracker(ctx, receiverChain, outTxTssNonce)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *cctx)
}

func (k Keeper) SaveSucessfullOutBound(ctx sdk.Context, cctx *types.CrossChainTx, ballotIndex string) {
	receiverChain := cctx.GetCurrentOutTxParam().ReceiverChainId
	tssPubkey := cctx.GetCurrentOutTxParam().TssPubkey
	outTxTssNonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce

	cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex
	cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	// #nosec G701 always in range
	k.GetObserverKeeper().RemoveFromPendingNonces(ctx, tssPubkey, receiverChain, int64(outTxTssNonce))
	k.RemoveOutTxTracker(ctx, receiverChain, outTxTssNonce)
	ctx.Logger().Info(fmt.Sprintf("Remove tracker %s: , Block Height : %d ", getOutTrackerIndex(receiverChain, outTxTssNonce), ctx.BlockHeight()))
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *cctx)
}
