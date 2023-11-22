package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetBallot(ctx sdk.Context, ballot *types.Ballot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	ballot.Index = ballot.BallotIdentifier
	b := k.cdc.MustMarshal(ballot)
	store.Set([]byte(ballot.Index), b)
}

func (k Keeper) SetBallotList(ctx sdk.Context, ballotlist *types.BallotListForHeight) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BallotListKey))
	b := k.cdc.MustMarshal(ballotlist)
	store.Set(types.BallotListKeyPrefix(ballotlist.Height), b)
}

func (k Keeper) GetBallot(ctx sdk.Context, index string) (val types.Ballot, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetBallotList(ctx sdk.Context, height int64) (val types.BallotListForHeight, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BallotListKey))
	b := store.Get(types.BallotListKeyPrefix(height))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllBallots(ctx sdk.Context) (voters []*types.Ballot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Ballot
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		voters = append(voters, &val)
	}
	return
}

// AddBallotToList adds a ballot to the list of ballots for a given height.
func (k Keeper) AddBallotToList(ctx sdk.Context, ballot types.Ballot) {
	list, found := k.GetBallotList(ctx, ballot.BallotCreationHeight)
	if !found {
		list = types.BallotListForHeight{Height: ballot.BallotCreationHeight, BallotsIndexList: []string{}}
	}
	list.BallotsIndexList = append(list.BallotsIndexList, ballot.BallotIdentifier)
	k.SetBallotList(ctx, &list)
}

// GetMaturedBallotList Returns a list of ballots which are matured at current height
func (k Keeper) GetMaturedBallotList(ctx sdk.Context) []string {
	maturityBlocks := k.GetParams(ctx).BallotMaturityBlocks
	list, found := k.GetBallotList(ctx, ctx.BlockHeight()-maturityBlocks)
	if !found {
		return []string{}
	}
	return list.BallotsIndexList
}

// Queries

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

func (k Keeper) BallotByIdentifier(goCtx context.Context, req *types.QueryBallotByIdentifierRequest) (*types.QueryBallotByIdentifierResponse, error) {
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
