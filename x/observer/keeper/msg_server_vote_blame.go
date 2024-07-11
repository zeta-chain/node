package keeper

import (
	"context"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschainTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) VoteBlame(
	goCtx context.Context,
	msg *types.MsgVoteBlame,
) (*types.MsgVoteBlameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := types.ObservationType_TSSKeySign

	// GetChainFromChainID makes sure we are getting only supported chains , if a chain support has been turned on using gov proposal, this function returns nil
	observationChain, found := k.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, cosmoserrors.Wrap(
			crosschainTypes.ErrUnsupportedChain,
			fmt.Sprintf("ChainID %d, Blame vote", msg.ChainId),
		)
	}

	if ok := k.IsNonTombstonedObserver(ctx, msg.Creator); !ok {
		return nil, types.ErrNotObserver
	}

	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	ballot, isNew, err := k.FindBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
	}

	if isNew {
		EmitEventBallotCreated(ctx, ballot, msg.BlameInfo.Index, observationChain.String())
	}

	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, types.VoteType_SuccessObservation)
	if err != nil {
		return nil, err
	}

	_, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgVoteBlameResponse{}, nil
	}

	// ******************************************************************************
	// below only happens when ballot is finalized: exactly when threshold vote is in
	// ******************************************************************************

	k.SetBlame(ctx, msg.BlameInfo)
	return &types.MsgVoteBlameResponse{}, nil
}
