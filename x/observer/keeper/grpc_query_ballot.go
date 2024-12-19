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
	voterList, err := ballot.GenerateVoterList()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBallotByIdentifierResponse{
		BallotIdentifier: ballot.BallotIdentifier,
		Voters:           voterList,
		ObservationType:  ballot.ObservationType,
		BallotStatus:     ballot.BallotStatus,
	}, nil
}

func (k Keeper) Ballots(goCtx context.Context, req *types.QueryBallotsRequest) (*types.QueryBallotsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	ballotList := k.GetAllBallots(ctx)
	ballots := make([]types.Ballot, len(ballotList))
	for i, ballot := range ballotList {
		ballots[i] = *ballot
	}
	return &types.QueryBallotsResponse{Ballots: ballots}, nil
}
