package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_InBoundTx
	observationChain := zetaObserverTypes.ParseCommonChaintoObservationChain(msg.SenderChain)
	//Check is msg.Creator is authorized to vote
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain, observationType.String())
	if !ok {
		return nil, err
	}

	index := msg.Digest()
	// Add votes and Set Ballot
	ballot, err := k.GetBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
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
	// Inbound Ballot has been finalized , Create CCTX
	// New CCTX can only set either to Aborted or PendingOutbound
	cctx := k.CreateNewCCTX(ctx, msg, index)
	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail
	err = k.FinalizeInbound(ctx, &cctx, msg.ReceiverChain, len(ballot.VoterList))
	if err != nil {
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
		ctx.Logger().Error(err.Error())
		k.SetCrossChainTx(ctx, cctx)
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingOutbound, "Status Changed to Pending Outbound", cctx.LogIdentifierForCCTX())
	EmitEventInboundFinalized(ctx, &cctx)
	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil

}

func (k msgServer) FinalizeInbound(ctx sdk.Context, cctx *types.CrossChainTx, receiveChain string, numberofobservers int) error {
	cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	k.UpdateLastBlockHeight(ctx, cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % numberofobservers)

	err := k.UpdatePrices(ctx, receiveChain, cctx)
	if err != nil {
		return err
	}
	err = k.UpdateNonce(ctx, receiveChain, cctx)
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
