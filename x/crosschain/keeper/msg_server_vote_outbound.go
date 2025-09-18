package keeper

import (
	"context"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	cctxerror "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/x/crosschain/types"
	observerkeeper "github.com/zeta-chain/node/x/observer/keeper"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const voteOutboundID = "Vote Outbound"

// VoteOutbound casts a vote on an outbound transaction observed on a connected chain (after
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
func (k msgServer) VoteOutbound(
	goCtx context.Context,
	msg *types.MsgVoteOutbound,
) (*types.MsgVoteOutboundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if TSS exists
	// It is almost impossible to reach this point without a TSS,
	// as the check for TSS was already done when creating the inbound, but we check anyway.
	// We also expect the tss.Pubkey to be the same as the one in the outbound params,
	// As we would have processed all CCTXs before migrating
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, cosmoserrors.Wrap(types.ErrCannotFindTSSKeys, voteOutboundID)
	}

	// Validate the message params to verify it against an existing CCTX.
	cctx, err := k.ValidateOutboundMessage(ctx, *msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(err, voteOutboundID)
	}

	ballotIndex := msg.Digest()
	isFinalizingVote, isNew, ballot, observationChain, err := k.zetaObserverKeeper.VoteOnOutboundBallot(
		ctx,
		ballotIndex,
		msg.OutboundChain,
		msg.Status,
		msg.Creator)
	if err != nil {
		return nil, cosmoserrors.Wrap(err, voteOutboundID)
	}

	// If the ballot is new, set the index to the CCTX.
	if isNew {
		observerkeeper.EmitEventBallotCreated(ctx, ballot, msg.ObservedOutboundHash, observationChain)
	}

	// If not finalized commit state here.
	if !isFinalizingVote {
		return &types.MsgVoteOutboundResponse{}, nil
	}

	// If the CCTX is in a terminal state, we do not need to process it.
	if cctx.CctxStatus.Status.IsTerminal() {
		return &types.MsgVoteOutboundResponse{}, cosmoserrors.Wrap(
			types.ErrCCTXAlreadyFinalized,
			fmt.Sprintf("CCTX status %s", cctx.CctxStatus.Status),
		)
	}

	// Set the finalized ballot to the current outbound params.
	cctx.SetOutboundBallotIndex(ballotIndex)

	err = cctx.UpdateCurrentOutbound(ctx, *msg, ballot.BallotStatus)
	if err != nil {
		return nil, cosmoserrors.Wrap(err, voteOutboundID)
	}

	// Fund the gas stability pool with the remaining funds.
	k.ManageUnusedGasFee(ctx, &cctx)

	// Validate and process the observed outbound
	err = k.ValidateOutboundObservers(ctx, &cctx, ballot.BallotStatus, msg.ValueReceived.String())
	if err != nil {
		// Should not happen
		cctx.SetAbort(types.StatusMessages{
			StatusMessage:        "outbound failed unable to process",
			ErrorMessageOutbound: cctxerror.NewCCTXErrorJSONMessage("", err),
		})
		ctx.Logger().Error(err.Error())
	}
	k.SaveOutbound(ctx, &cctx, tss.TssPubkey)

	return &types.MsgVoteOutboundResponse{}, nil
}

// ManageUnusedGasFee funds the stability pool with the remaining fees of an outbound tx
// The funds are sent to the gas stability pool associated with the receiver chain and the user is refunded if possible.
// This wraps the UseRemainingGasFee function and logs an error if it fails.We do not return an error here.
// Event if the funding fails, the outbound tx is still processed.
func (k Keeper) ManageUnusedGasFee(ctx sdk.Context, cctx *types.CrossChainTx) {
	outboundParams := cctx.GetCurrentOutboundParam()
	// We skip funding the gas stability pool if the userGasFeePaid is nil or zero.
	// This is a legacy outbound where the user fee was not recorded as part of the cctx struct. This is only to handle cctxs which might be in pending outbound state when the upgrade happens and get finalized after the upgrade
	if outboundParams.UserGasFeePaid.IsNil() || outboundParams.UserGasFeePaid.IsZero() {
		err := k.FundGasStabilityPoolFromRemainingFees(
			ctx,
			*outboundParams,
			outboundParams.ReceiverChainId,
		)
		if err != nil {
			ctx.Logger().Error("Failed to fund gas stability pool (Legacy) with remaining fees",
				"voteOutboundID", voteOutboundID,
				"cctxIndex", cctx.Index,
				"error", err,
			)

		}
		return
	}

	if err := k.useRemainingGasFee(ctx, cctx); err != nil {
		ctx.Logger().Error("Failed to fund gas stability pool with remaining fees",
			"voteOutboundID", voteOutboundID,
			"cctxIndex", cctx.Index,
			"error", err,
		)
	}
}

// UseRemainingGasFee uses the remaining fees of an outbound tx to fund the gas stability pool and refund the user if possible
func (k Keeper) useRemainingGasFee(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
) error {
	outboundParams := cctx.GetCurrentOutboundParam()
	outboundTxFeePaid := math.NewUint(outboundParams.GasUsed).
		Mul(math.NewUintFromBigInt(outboundParams.EffectiveGasPrice.BigInt()))
	// We already know that the userGasFeePaid is not nil
	userGasFeePaid := outboundParams.UserGasFeePaid

	// The final fee paid is greater than what the user paid originally.The stability pool would cover the extra cost in this case.
	if outboundTxFeePaid.GTE(userGasFeePaid) {
		return nil
	}

	// We use a maximum for 95 percent of the remaining fees to fund the stability pool and refund the user.
	usableRemainingFees := PercentOf(
		outboundParams.UserGasFeePaid.Sub(outboundTxFeePaid),
		types.UsableRemainingFeesPercentage,
	)
	if !usableRemainingFees.GT(math.ZeroUint()) {
		return nil
	}

	// Send all tokens to stability pool by default
	stabilityPoolPercentage := uint64(100)
	refundToUser := false

	// Refund to the sender irrespective of whether its EOA or contract address if it's a withdrawal originating from zEVM.
	if chains.IsZetaChain(cctx.InboundParams.SenderChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx)) &&
		ethcommon.IsHexAddress(cctx.InboundParams.Sender) {
		chainParams, found := k.GetObserverKeeper().GetChainParamsByChainID(ctx, outboundParams.ReceiverChainId)
		if !found {
			return errors.Wrap(
				observertypes.ErrChainParamsNotFound,
				fmt.Sprintf("chainID: %d", outboundParams.ReceiverChainId),
			)
		}
		refundToUser = true
		stabilityPoolPercentage = chainParams.StabilityPoolPercentage
	}
	stabilityPoolAmount := PercentOf(usableRemainingFees, stabilityPoolPercentage)
	refundAmount := usableRemainingFees.Sub(stabilityPoolAmount)

	if stabilityPoolAmount.GT(math.ZeroUint()) {
		if err := k.fungibleKeeper.FundGasStabilityPool(ctx, outboundParams.ReceiverChainId, stabilityPoolAmount.BigInt()); err != nil {
			return err
		}
	}

	if refundAmount.GT(math.ZeroUint()) && refundToUser {
		if err := k.fungibleKeeper.RefundRemainingGasFees(ctx, outboundParams.ReceiverChainId, refundAmount.BigInt(), ethcommon.HexToAddress(cctx.InboundParams.Sender)); err != nil {
			return err
		}
	}

	return nil
}

// TODO : Remove this function
// https://github.com/zeta-chain/node/issues/3768
// FundGasStabilityPoolFromRemainingFees funds the gas stability pool with the remaining fees of an outbound tx
func (k Keeper) FundGasStabilityPoolFromRemainingFees(
	ctx sdk.Context,
	OutboundParams types.OutboundParams,
	chainID int64,
) error {
	gasUsed := OutboundParams.GasUsed
	gasLimit := OutboundParams.EffectiveGasLimit
	gasPrice := math.NewUintFromBigInt(OutboundParams.EffectiveGasPrice.BigInt())

	if gasLimit == gasUsed {
		return nil
	}

	// We skip gas stability pool funding if one of the params is zero
	if gasLimit > 0 && gasUsed > 0 && !gasPrice.IsZero() {
		if gasLimit > gasUsed {
			remainingGas := gasLimit - gasUsed
			remainingFees := math.NewUint(remainingGas).Mul(gasPrice)

			// We fund the stability pool with a portion of the remaining fees.
			remainingFees = PercentOf(remainingFees, RemainingFeesToStabilityPoolPercent)

			// Fund the gas stability pool.
			if err := k.fungibleKeeper.FundGasStabilityPool(ctx, chainID, remainingFees.BigInt()); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s: The gas limit %d is less than the gas used %d", voteOutboundID, gasLimit, gasUsed)
		}
	}

	return nil
}

// PercentOf returns the percentage of a number
func PercentOf(n math.Uint, percent uint64) math.Uint {
	percentUint := math.NewUint(percent)
	result := n.Mul(percentUint)
	return result.Quo(math.NewUint(100))
}

// SaveOutbound saves the outbound transaction.It does the following things in one function:
// 1. Set the ballot index for the outbound vote to the cctx
// 2. Remove the nonce from the pending nonces
// 3. Remove the outbound tx tracker
// 4. Set the cctx and nonce to cctx and inboundHash to cctx
func (k Keeper) SaveOutbound(ctx sdk.Context, cctx *types.CrossChainTx, tssPubkey string) {
	// #nosec G115 always in range
	for _, outboundParams := range cctx.OutboundParams {
		k.GetObserverKeeper().
			RemoveFromPendingNonces(ctx, outboundParams.TssPubkey, outboundParams.ReceiverChainId, int64(outboundParams.TssNonce))
		k.RemoveOutboundTrackerFromStore(ctx, outboundParams.ReceiverChainId, outboundParams.TssNonce)
	}
	// This should set nonce to cctx only if a new revert is created.
	k.SaveCCTXUpdate(ctx, *cctx, tssPubkey)
}

func (k Keeper) ValidateOutboundMessage(ctx sdk.Context, msg types.MsgVoteOutbound) (types.CrossChainTx, error) {
	// Check if CCTX exists and if the nonce matches.
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxHash)
	if !found {
		return types.CrossChainTx{}, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"%s, CCTX %s does not exist", voteOutboundID, msg.CctxHash)
	}

	if cctx.GetCurrentOutboundParam().TssNonce != msg.OutboundTssNonce {
		return types.CrossChainTx{}, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"%s, OutboundTssNonce %d does not match CCTX OutboundTssNonce %d",
			voteOutboundID,
			msg.OutboundTssNonce,
			cctx.GetCurrentOutboundParam().TssNonce,
		)
	}

	if cctx.GetCurrentOutboundParam().ReceiverChainId != msg.OutboundChain {
		return types.CrossChainTx{}, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"%s, OutboundChain %d does not match CCTX OutboundChain %d",
			voteOutboundID,
			msg.OutboundChain,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
		)
	}

	return cctx, nil
}
