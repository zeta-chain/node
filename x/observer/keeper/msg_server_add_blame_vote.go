package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crosschainTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) AddBlameVote(goCtx context.Context, vote *types.MsgAddBlameVote) (*types.MsgAddBlameVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := types.ObservationType_TSSKeySign
	// GetChainFromChainID makes sure we are getting only supported chains , if a chain support has been turned on using gov proposal, this function returns nil
	observationChain := k.GetParams(ctx).GetChainFromChainID(vote.ChainId)
	if observationChain == nil {
		return nil, sdkerrors.Wrap(crosschainTypes.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Blame vote", vote.ChainId))
	}
	// IsAuthorized does various checks against the list of observer mappers
	ok, err := k.IsAuthorized(ctx, vote.Creator, observationChain)
	if !ok {
		return nil, err
	}

	index := vote.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	ballot, isNew, err := k.FindBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
	}

	if isNew {
		types.EmitEventBallotCreated(ctx, ballot, vote.BlameInfo.Index, observationChain.String())
	}

	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, vote.Creator, types.VoteType_SuccessObservation)
	if err != nil {
		return nil, err
	}

	ballot, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgAddBlameVoteResponse{}, nil
	}

	// ******************************************************************************
	// below only happens when ballot is finalized: exactly when threshold vote is in
	// ******************************************************************************

	k.SetBlame(ctx, vote.BlameInfo)
	return &types.MsgAddBlameVoteResponse{}, nil
}
