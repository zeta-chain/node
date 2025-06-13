package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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

	votersList := make([]*types.VoterList, len(ballot.VoterList))
	for i, voterAddress := range ballot.VoterList {
		voter := types.VoterList{
			VoterAddress: voterAddress,
			VoteType:     ballot.Votes[ballot.GetVoterIndex(voterAddress)],
		}
		votersList[i] = &voter
	}

	return &types.QueryBallotByIdentifierResponse{
		BallotIdentifier: ballot.BallotIdentifier,
		Voters:           votersList,
		ObservationType:  ballot.ObservationType,
		BallotStatus:     ballot.BallotStatus,
	}, nil
}

// Ballots queries all the ballots. It is a paginated query
func (k Keeper) Ballots(goCtx context.Context, req *types.QueryBallotsRequest) (*types.QueryBallotsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ballots := make([]types.Ballot, 0)
	ctx := sdk.UnwrapSDKContext(goCtx)

	ballotStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{}
	}

	pageRes, err := query.Paginate(ballotStore, req.Pagination, func(_ []byte, value []byte) error {
		var ballot types.Ballot
		if err := k.cdc.Unmarshal(value, &ballot); err != nil {
			return err
		}
		ballots = append(ballots, ballot)
		return nil
	})

	sort.Slice(ballots, func(i, j int) bool {
		return ballots[i].BallotCreationHeight < ballots[j].BallotCreationHeight
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryBallotsResponse{Ballots: ballots, Pagination: pageRes}, nil
}

func (k Keeper) BallotListForHeight(
	goCtx context.Context,
	req *types.QueryBallotListForHeightRequest,
) (*types.QueryBallotListForHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	ballotList, found := k.GetBallotListForHeight(ctx, req.Height)
	if !found {
		return nil, status.Error(codes.NotFound, "not found ballot list")
	}
	return &types.QueryBallotListForHeightResponse{
		BallotList: ballotList,
	}, nil
}
