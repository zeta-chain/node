package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// FIXME: use more specific error types & codes

// Casts a vote on an inbound transaction observed on a connected chain. If this
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
// with the CCTX are swapped using a Uniswap pool on ZetaChain for the ZRC20 of
// the gas token of the receiver chain. The ZRC20 tokens are then burned. If
// everything is successful, the CCTX status is changed to `PendingOutbound`.
//
// ```mermaid
// stateDiagram-v2
//
//	state is_zetachain <<choice>>
//	state evm_deposit_success <<choice>>
//	state finalize_inbound <<choice>>
//	[*] --> PendingInbound: Create New CCTX
//	PendingInbound --> is_zetachain
//	is_zetachain --> evm_deposit_success: Receiver is ZetaChain
//	evm_deposit_success --> OutboundMined: EVM deposit success
//	evm_deposit_success --> Aborted: EVM deposit error
//	is_zetachain --> finalize_inbound: Receiver is connected chain
//	finalize_inbound --> Aborted: Finalize inbound error
//	finalize_inbound --> PendingOutbound: Finalize inbound success
//	Aborted --> [*]
//	PendingOutbound --> [*]
//	OutboundMined --> [*]
//
// ```
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_InBoundTx
	if !k.IsInboundAllowed(ctx) {
		return nil, types.ErrNotEnoughPermissions
	}
	// GetChainFromChainID makes sure we are getting only supported chains , if a chain support has been turned on using gov proposal, this function returns nil
	observationChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.SenderChainId)
	if observationChain == nil {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.SenderChainId, observationType.String()))
	}
	receiverChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ReceiverChain)
	if receiverChain == nil {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.ReceiverChain, observationType.String()))
	}
	// IsAuthorized does various checks against the list of observer mappers
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain)
	if !ok {
		return nil, err
	}

	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	ballot, isNew, err := k.GetBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	if isNew {
		EmitEventBallotCreated(ctx, ballot, msg.InTxHash, observationChain.String())
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.VoteType_SuccessObservation)
	if err != nil {
		return nil, err
	}
	// CheckIfBallotIsFinalized checks status and sets the ballot if finalized

	ballot, isFinalized := k.CheckIfBallotIsFinalized(ctx, ballot)
	if !isFinalized {
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	// Validation if we want to send ZETA to external chain, but there is no ZETA token.
	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, receiverChain.ChainId)
	if !found {
		return nil, types.ErrNotFoundCoreParams
	}
	if receiverChain.IsExternalChain() && coreParams.ZetaTokenContractAddress == "" && msg.CoinType == common.CoinType_Zeta {
		return nil, types.ErrUnableToSendCoinType
	}

	// ******************************************************************************
	// below only happens when ballot is finalized: exactly when threshold vote is in
	// ******************************************************************************

	// Inbound Ballot has been finalized , Create CCTX
	cctx := k.CreateNewCCTX(ctx, msg, index, types.CctxStatus_PendingInbound, observationChain, receiverChain)
	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail

	if receiverChain.IsZetaChain() {
		err = k.HandleEVMDeposit(ctx, &cctx, *msg, observationChain)
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_OutboundMined, "First half of EVM transfer Completed", cctx.LogIdentifierForCCTX())
	} else { // Cross Chain SWAP
		err = k.FinalizeInbound(ctx, &cctx, *receiverChain, len(ballot.VoterList))
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}

		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingOutbound, "Status Changed to Pending Outbound", cctx.LogIdentifierForCCTX())
	}
	EmitEventInboundFinalized(ctx, &cctx)
	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}

// TODO: is LastBlockHeight needed?
func (k msgServer) FinalizeInbound(ctx sdk.Context, cctx *types.CrossChainTx, receiveChain common.Chain, numberofobservers int) error {
	cctx.InboundTxParams.InboundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	//k.UpdateLastBlockHeight(ctx, cctx)

	err := k.UpdatePrices(ctx, receiveChain.ChainId, cctx)
	if err != nil {
		return err
	}
	err = k.UpdateNonce(ctx, receiveChain.ChainId, cctx)
	if err != nil {
		return err
	}
	return nil
}

//func (k msgServer) UpdateLastBlockHeight(ctx sdk.Context, msg *types.CrossChainTx) {
//	lastblock, isFound := k.GetLastBlockHeight(ctx, msg.InboundTxParams.SenderChain)
//	if !isFound {
//		lastblock = types.LastBlockHeight{
//			Creator:           msg.Creator,
//			Index:             msg.InboundTxParams.SenderChain, // ?
//			Chain:             msg.InboundTxParams.SenderChain,
//			LastSendHeight:    msg.InboundTxParams.InboundTxObservedExternalHeight,
//			LastReceiveHeight: 0,
//		}
//	} else {
//		lastblock.LastSendHeight = msg.InboundTxParams.InboundTxObservedExternalHeight
//	}
//	k.SetLastBlockHeight(ctx, lastblock)
//}
