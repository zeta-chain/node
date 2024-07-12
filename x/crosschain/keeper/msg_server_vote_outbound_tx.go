package keeper

import (
	"context"
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerkeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
)

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
		msg.OutboundChain,
		msg.Status,
		msg.Creator)
	if err != nil {
		return nil, err
	}
	// if the ballot is new, set the index to the CCTX
	if isNew {
		observerkeeper.EmitEventBallotCreated(ctx, ballot, msg.ObservedOutboundHash, observationChain)
	}
	// if not finalized commit state here
	if !isFinalizingVote {
		return &types.MsgVoteOutboundResponse{}, nil
	}

	// if ballot successful, the value received should be the out tx amount
	err = cctx.AddOutbound(ctx, *msg, ballot.BallotStatus)
	if err != nil {
		return nil, err
	}
	// Fund the gas stability pool with the remaining funds
	k.FundStabilityPool(ctx, &cctx)

	err = k.ValidateOutboundObservers(ctx, &cctx, ballot.BallotStatus, msg.ValueReceived.String())
	if err != nil {
		k.SaveFailedOutbound(ctx, &cctx, err.Error(), ballotIndex)
		return &types.MsgVoteOutboundResponse{}, nil
	}
	k.SaveSuccessfulOutbound(ctx, &cctx, ballotIndex)
	return &types.MsgVoteOutboundResponse{}, nil
}

// FundStabilityPool funds the stability pool with the remaining fees of an outbound tx
// The funds are sent to the gas stability pool associated with the receiver chain
// This wraps the FundGasStabilityPoolFromRemainingFees function and logs an error if it fails.We do not return an error here.
// Event if the funding fails, the outbound tx is still processed.
func (k Keeper) FundStabilityPool(ctx sdk.Context, cctx *types.CrossChainTx) {
	// Fund the gas stability pool with the remaining funds
	if err := k.FundGasStabilityPoolFromRemainingFees(ctx, *cctx.GetCurrentOutboundParam(), cctx.GetCurrentOutboundParam().ReceiverChainId); err != nil {
		ctx.Logger().
			Error(fmt.Sprintf("VoteOutbound: CCTX: %s Can't fund the gas stability pool with remaining fees %s", cctx.Index, err.Error()))
	}
}

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
			remainingFees := math.NewUint(remainingGas).Mul(gasPrice).BigInt()

			// We fund the stability pool with a portion of the remaining fees
			remainingFees = percentOf(remainingFees, RemainingFeesToStabilityPoolPercent)
			// Fund the gas stability pool
			if err := k.fungibleKeeper.FundGasStabilityPool(ctx, chainID, remainingFees); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("VoteOutbound: The gas limit %d is less than the gas used %d", gasLimit, gasUsed)
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

/*
SaveFailedOutbound saves a failed outbound transaction.It does the following things in one function:

 1. Change the status of the CCTX to Aborted

 2. Save the outbound
*/

func (k Keeper) SaveFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, errMessage string, ballotIndex string) {
	cctx.SetAbort(errMessage)
	ctx.Logger().Error(errMessage)

	k.SaveOutbound(ctx, cctx, ballotIndex)
}

// SaveSuccessfulOutbound saves a successful outbound transaction.
// This function does not set the CCTX status, therefore all successful outbound transactions need
// to have their status set during processing
func (k Keeper) SaveSuccessfulOutbound(ctx sdk.Context, cctx *types.CrossChainTx, ballotIndex string) {
	k.SaveOutbound(ctx, cctx, ballotIndex)
}

/*
SaveOutbound saves the outbound transaction.It does the following things in one function:

 1. Set the ballot index for the outbound vote to the cctx

 2. Remove the nonce from the pending nonces

 3. Remove the outbound tx tracker

 4. Set the cctx and nonce to cctx and inboundHash to cctx
*/
func (k Keeper) SaveOutbound(ctx sdk.Context, cctx *types.CrossChainTx, ballotIndex string) {
	receiverChain := cctx.GetCurrentOutboundParam().ReceiverChainId
	tssPubkey := cctx.GetCurrentOutboundParam().TssPubkey
	outTxTssNonce := cctx.GetCurrentOutboundParam().TssNonce

	cctx.GetCurrentOutboundParam().BallotIndex = ballotIndex
	// #nosec G115 always in range
	k.GetObserverKeeper().RemoveFromPendingNonces(ctx, tssPubkey, receiverChain, int64(outTxTssNonce))
	k.RemoveOutboundTrackerFromStore(ctx, receiverChain, outTxTssNonce)
	ctx.Logger().
		Info(fmt.Sprintf("Remove tracker %s: , Block Height : %d ", getOutboundTrackerIndex(receiverChain, outTxTssNonce), ctx.BlockHeight()))
	// This should set nonce to cctx only if a new revert is created.
	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, *cctx)
}

func (k Keeper) ValidateOutboundMessage(ctx sdk.Context, msg types.MsgVoteOutbound) (types.CrossChainTx, error) {
	// check if CCTX exists and if the nonce matches
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxHash)
	if !found {
		return types.CrossChainTx{}, cosmoserrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("CCTX %s does not exist", msg.CctxHash),
		)
	}
	if cctx.GetCurrentOutboundParam().TssNonce != msg.OutboundTssNonce {
		return types.CrossChainTx{}, cosmoserrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf(
				"OutboundTssNonce %d does not match CCTX OutboundTssNonce %d",
				msg.OutboundTssNonce,
				cctx.GetCurrentOutboundParam().TssNonce,
			),
		)
	}
	// do not process an outbound vote if TSS is not found
	_, found = k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return types.CrossChainTx{}, types.ErrCannotFindTSSKeys
	}
	if cctx.GetCurrentOutboundParam().ReceiverChainId != msg.OutboundChain {
		return types.CrossChainTx{}, cosmoserrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf(
				"OutboundChain %d does not match CCTX OutboundChain %d",
				msg.OutboundChain,
				cctx.GetCurrentOutboundParam().ReceiverChainId,
			),
		)
	}
	return cctx, nil
}
