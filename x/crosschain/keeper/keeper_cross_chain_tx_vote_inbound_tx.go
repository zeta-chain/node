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
func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_InBoundTx
	observationChain, found := k.zetaObserverKeeper.GetChainFromChainID(ctx, msg.SenderChain)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.SenderChain, observationType.String()))
	}
	receiverChain, found := k.zetaObserverKeeper.GetChainFromChainID(ctx, msg.ReceiverChain)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.ReceiverChain, observationType.String()))
	}
	// IsAuthorized does various checks against the list of observer mappers
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain, observationType)
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

	// ******************************************************************************
	// below only happens when ballot is finalized: exactly when threshold vote is in
	// ******************************************************************************

	// Inbound Ballot has been finalized , Create CCTX
	// New CCTX can only set either to Aborted or PendingOutbound
	cctx := k.CreateNewCCTX(ctx, msg, index, types.CctxStatus_PendingInbound, observationChain, receiverChain)
	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail
	switch receiverChain.ChainName {
	case common.ChainName_ZetaChain:
		err = k.HandleEVMDeposit(ctx, &cctx, *msg, observationChain)
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_OutboundMined, "First half of EVM transfer Completed", cctx.LogIdentifierForCCTX())
	default: // Cross Chain SWAP
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

func (k msgServer) FinalizeInbound(ctx sdk.Context, cctx *types.CrossChainTx, receiveChain common.Chain, numberofobservers int) error {
	cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	k.UpdateLastBlockHeight(ctx, cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % numberofobservers)

	err := k.UpdatePrices(ctx, receiveChain.ChainId, cctx)
	if err != nil {
		return err
	}
	err = k.UpdateNonce(ctx, receiveChain.ChainName.String(), cctx)
	if err != nil {
		return err
	}
	return nil
}

func (k msgServer) UpdateLastBlockHeight(ctx sdk.Context, msg *types.CrossChainTx) {
	lastblock, isFound := k.GetLastBlockHeight(ctx, msg.InBoundTxParams.SenderChain)
	if !isFound {
		lastblock = types.LastBlockHeight{
			Creator:           msg.Creator,
			Index:             msg.InBoundTxParams.SenderChain, // ?
			Chain:             msg.InBoundTxParams.SenderChain,
			LastSendHeight:    msg.InBoundTxParams.InBoundTxObservedExternalHeight,
			LastReceiveHeight: 0,
		}
	} else {
		lastblock.LastSendHeight = msg.InBoundTxParams.InBoundTxObservedExternalHeight
	}
	k.SetLastBlockHeight(ctx, lastblock)
}
