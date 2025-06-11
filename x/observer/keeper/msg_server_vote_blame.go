package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/observer/types"
)

const voteBlameID = "Vote Blame"

func (k msgServer) VoteBlame(
	goCtx context.Context,
	msg *types.MsgVoteBlame,
) (*types.MsgVoteBlameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// GetChainFromChainID makes sure we are getting only supported chains , if a chain support has been turned on using gov proposal, this function returns nil
	observationChain, found := k.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, sdkerrors.Wrapf(cctypes.ErrUnsupportedChain, "%s, ChainID %d", voteBlameID, msg.ChainId)
	}

	err := k.CheckObserverCanVote(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}

	ballot, isFinalized, isNew, err := k.VoteOnBallot(
		ctx,
		observationChain,
		msg.Digest(),
		types.ObservationType_TSSKeySign,
		msg.Creator,
		types.VoteType_SuccessObservation,
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			err,
			"%s, BallotIdentifier %v", voteBlameID, ballot.BallotIdentifier)
	}

	if isNew {
		EmitEventBallotCreated(ctx, ballot, msg.BlameInfo.Index, observationChain.String())
	}

	if !isFinalized {
		// Return nil here to add vote to ballot and commit state.
		return &types.MsgVoteBlameResponse{}, nil
	}

	// Ballot is finalized: exactly when threshold vote is in.
	k.SetBlame(ctx, msg.BlameInfo)
	return &types.MsgVoteBlameResponse{}, nil
}
