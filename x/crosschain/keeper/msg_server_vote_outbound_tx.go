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

	// If ballot is successful, the value received should be the out tx amount.
	err = cctx.UpdateCurrentOutbound(ctx, *msg, ballot.BallotStatus)
	if err != nil {
		return nil, cosmoserrors.Wrap(err, voteOutboundID)
	}

	// Fund the gas stability pool with the remaining funds.
	k.FundStabilityPool(ctx, &cctx)

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

// FundStabilityPool funds the stability pool with the remaining fees of an outbound tx
// The funds are sent to the gas stability pool associated with the receiver chain
// This wraps the FundGasStabilityPoolFromRemainingFees function and logs an error if it fails.We do not return an error here.
// Event if the funding fails, the outbound tx is still processed.

func (k Keeper) FundStabilityPool(ctx sdk.Context, cctx *types.CrossChainTx) {
	// Fund the gas stability pool with the remaining funds
	if err := k.FundGasStabilityPoolFromRemainingFees(ctx,
		*cctx.GetCurrentOutboundParam(),
		cctx.GetCurrentOutboundParam().ReceiverChainId,
		cctx.InboundParams.SenderChainId,
		cctx.InboundParams.Sender); err != nil {
		ctx.Logger().
			Error("%s: CCTX: %s Can't fund the gas stability pool with remaining fees %s", voteOutboundID, cctx.Index, err.Error())
	}
}

// FundGasStabilityPoolFromRemainingFees funds the gas stability pool with the remaining fees of an outbound tx
func (k Keeper) FundGasStabilityPoolFromRemainingFees(
	ctx sdk.Context,
	OutboundParams types.OutboundParams,
	receiverChainID int64,
	senderChainID int64,
	sender string,
) error {
	outboundTxGasUsed := math.NewUint(OutboundParams.GasUsed)
	outboundTxFinalGasPrice := math.NewUintFromBigInt(OutboundParams.EffectiveGasPrice.BigInt())
	outboundTxFeePaid := outboundTxGasUsed.Mul(outboundTxFinalGasPrice)

	userGasFeePaid := OutboundParams.UserGasFeePaid

	// The final fee paid is greater than what the user paid originally.The stability pool would cover the extra cost in this case.
	if outboundTxFeePaid.GTE(userGasFeePaid) {
		return nil
	}

	remainingFees := userGasFeePaid.Sub(outboundTxFeePaid)

	chainParams, found := k.GetObserverKeeper().GetChainParamsByChainID(ctx, receiverChainID)
	if !found {
		return errors.Wrap(observertypes.ErrChainParamsNotFound, fmt.Sprintf("chainID: %d", receiverChainID))
	}

	// Send all tokens to stability pool by default
	stabilityPoolPercentage := uint64(100)
	refundToUser := false
	// Refund tokens to user if it's a withdrawal originating from zEVM.
	// Refund to the sender irrespective of weather its EOA or contract address
	// For v1 msg passing, we cannot refund the user on zEVM
	if chains.IsZetaChain(
		senderChainID,
		k.GetAuthorityKeeper().GetAdditionalChainList(ctx),
	) && ethcommon.IsHexAddress(sender) {
		refundToUser = true
		stabilityPoolPercentage = chainParams.StabilityPoolPercentage
	}

	stabilityPoolAmount := PercentOf(remainingFees, stabilityPoolPercentage)
	// Refund the remaining fees to the user
	// For v2 withdraw: The fees are paid by burning GASZRC20 tokens of a receiver chain. So we can directly refund the calculated amount to user in the same tokens.
	// For v1 withdraw
	//  - ZRC20 withdrawal: The fees are paid by burning GASZRC20 tokens of a receiver chain. So we can directly refund the calculated amount to user in the same tokens.
	//  - Zeta withdrawal: We use a portion of the amount to buy gas tokens and burn it. We can still refund the user the remaining amount in GAS ZRC20 instead of ZETA.
	// For v1 Msg Passing:
	// - Zeta : We use a portion of the amount to buy gas tokens and burn it. We can still refund the user the remaining amount in GAS ZRC20 instead of ZETA.
	// - GAS : The fees are paid by burning GASZRC20 tokens of a receiver chain. So we can directly refund the calculated amount to user in the same tokens.
	// - ERC20 : We use a portion of the amount to buy gas tokens and burn it. We can still refund the user the remaining amount in GAS ZRC20 instead of ZETA.
	refundAmount := remainingFees.Sub(stabilityPoolAmount)
	refundAddress := ethcommon.HexToAddress(sender)

	if stabilityPoolAmount.GT(math.ZeroUint()) {
		if err := k.fungibleKeeper.FundGasStabilityPool(ctx, receiverChainID, stabilityPoolAmount.BigInt()); err != nil {
			return err
		}
	}

	if refundAmount.GT(math.ZeroUint()) && refundToUser {
		if err := k.fungibleKeeper.RefundRemainGasFess(ctx, receiverChainID, refundAmount.BigInt(), refundAddress); err != nil {
			return err
		}
	}
	return nil
}

// PercentOf returns the percentage of a number
func PercentOf(n math.Uint, percent uint64) math.Uint {
	// Convert percent to math.Uint
	percentUint := math.NewUint(percent)

	// Calculate n * percent
	result := n.Mul(percentUint)

	// Divide by 100
	hundred := math.NewUint(100)
	result = result.Quo(hundred)

	return result
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
