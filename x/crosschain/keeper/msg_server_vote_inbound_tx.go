package keeper

import (
	"context"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// FIXME: use more specific error types & codes

// VoteInbound casts a vote on an inbound transaction observed on a connected chain. If this
// is the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized. When a ballot is finalized, a new CCTX is
// created.
//
// If the receiver chain is ZetaChain, `HandleEVMDeposit` is called. If the
// tokens being deposited are ZETA, `MintZetaToEVMAccount` is called and the
// tokens are minted to the receiver account on ZetaChain. If the tokens being
// deposited are gas tokens or ERC20 of a connected chain, ZRC20's `deposit`
// method is called and the tokens are deposited to the receiver account on
// ZetaChain. If the message is not empty, system contract's `depositAndCall`
// method is also called and an omnichain contract on ZetaChain is executed.
// Omnichain contract address and arguments are passed as part of the message.
// If everything is successful, the CCTX status is changed to `OutboundMined`.
//
// If the receiver chain is a connected chain, the `FinalizeInbound` method is
// called to prepare the CCTX to be processed as an outbound transaction. To
// cover the outbound transaction fee, the required amount of tokens submitted
// with the CCTX are swapped using a Uniswap V2 contract instance on ZetaChain
// for the ZRC20 of the gas token of the receiver chain. The ZRC20 tokens are
// then burned. The nonce is updated. If everything is successful, the CCTX
// status is changed to `PendingOutbound`.
//
// ```mermaid
// stateDiagram-v2
//
//	state evm_deposit_success <<choice>>
//	state finalize_inbound <<choice>>
//	state evm_deposit_error <<choice>>
//	PendingInbound --> evm_deposit_success: Receiver is ZetaChain
//	evm_deposit_success --> OutboundMined: EVM deposit success
//	evm_deposit_success --> evm_deposit_error: EVM deposit error
//	evm_deposit_error --> PendingRevert: Contract error
//	evm_deposit_error --> Aborted: Internal error, invalid chain, gas, nonce
//	PendingInbound --> finalize_inbound: Receiver is connected chain
//	finalize_inbound --> Aborted: Finalize inbound error
//	finalize_inbound --> PendingOutbound: Finalize inbound success
//
// ```
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteInbound(
	goCtx context.Context,
	msg *types.MsgVoteInbound,
) (*types.MsgVoteInboundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	index := msg.Digest()

	// vote on inbound ballot
	// use a temporary context to not commit any ballot state change in case of error
	tmpCtx, commit := ctx.CacheContext()
	finalized, isNew, err := k.zetaObserverKeeper.VoteOnInboundBallot(
		tmpCtx,
		msg.SenderChainId,
		msg.ReceiverChain,
		msg.CoinType,
		msg.Creator,
		index,
		msg.InboundHash,
	)
	if err != nil {
		return nil, err
	}

	// If it is a new ballot, check if an inbound with the same hash, sender chain and event index has already been finalized
	// This may happen if the same inbound is observed twice where msg.Digest gives a different index
	// This check prevents double spending
	if isNew {
		if k.IsFinalizedInbound(tmpCtx, msg.InboundHash, msg.SenderChainId, msg.EventIndex) {
			return nil, cosmoserrors.Wrap(
				types.ErrObservedTxAlreadyFinalized,
				fmt.Sprintf(
					"inboundHash:%s, SenderChainID:%d, EventIndex:%d",
					msg.InboundHash,
					msg.SenderChainId,
					msg.EventIndex,
				),
			)
		}
	}
	commit()
	// If the ballot is not finalized return nil here to add vote to commit state
	if !finalized {
		return &types.MsgVoteInboundResponse{}, nil
	}

	cctx, err := k.ValidateInbound(ctx, msg, true)
	if err != nil {
		return nil, err
	}

	// Save the inbound CCTX to the store. This is called irrespective of the status of the CCTX or the outcome of the process function.
	k.SaveObservedInboundInformation(ctx, cctx, msg.EventIndex)
	return &types.MsgVoteInboundResponse{}, nil
}

/* SaveObservedInboundInformation saves the inbound CCTX to the store.It does the following:
    - Emits an event for the finalized inbound CCTX.
	- Adds the inbound CCTX to the finalized inbound CCTX store.This is done to prevent double spending, using the same inbound tx hash and event index.
	- Updates the CCTX with the finalized height and finalization status.
	- Removes the inbound CCTX from the inbound transaction tracker store.This is only for inbounds created via Inbound tracker suggestions
	- Sets the CCTX and nonce to the CCTX and inbound transaction hash to CCTX store.
*/

func (k Keeper) SaveObservedInboundInformation(ctx sdk.Context, cctx *types.CrossChainTx, eventIndex uint64) {
	EmitEventInboundFinalized(ctx, cctx)
	k.AddFinalizedInbound(ctx,
		cctx.GetInboundParams().ObservedHash,
		cctx.GetInboundParams().SenderChainId,
		eventIndex)
	// #nosec G115 always positive
	cctx.InboundParams.FinalizedZetaHeight = uint64(ctx.BlockHeight())
	cctx.InboundParams.TxFinalizationStatus = types.TxFinalizationStatus_Executed
	k.RemoveInboundTrackerIfExists(ctx, cctx.InboundParams.SenderChainId, cctx.InboundParams.ObservedHash)
	k.SetCrossChainTx(ctx, *cctx)
}
