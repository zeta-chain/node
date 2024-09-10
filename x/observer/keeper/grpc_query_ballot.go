package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/observer/types"
)

func (k Keeper) HasVoted(goCtx context.Context, req *types.QueryHasVotedRequest) (*types.QueryHasVotedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	ballot, found := k.GetBallot(ctx, req.BallotIdentifier)
	if !found {
		return &types.QueryHasVotedResponse{
			HasVoted: false,
		}, nil
	}
	hasVoted := ballot.HasVoted(req.VoterAddress)

	return &types.QueryHasVotedResponse{
		HasVoted: hasVoted,
	}, nil
}

func (k Keeper) BallotByIdentifier(
	goCtx context.Context,
	req *types.QueryBallotByIdentifierRequest,
) (*types.QueryBallotByIdentifierResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	ballot, found := k.GetBallot(ctx, req.BallotIdentifier)
	if !found {
		return nil, status.Error(codes.NotFound, "not found ballot")
	}

	return &types.QueryBallotByIdentifierResponse{
		BallotIdentifier: ballot.BallotIdentifier,
		Voters:           ballot.GenerateVoterList(),
		ObservationType:  ballot.ObservationType,
		BallotStatus:     ballot.BallotStatus,
	}, nil
}
