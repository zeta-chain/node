package keeper

import (
	"context"

	types2 "github.com/coinbase/rosetta-sdk-go/types"

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

func (k Keeper) GetBallot(ctx sdk.Context, index string) (val types.Ballot, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	b := store.Get(types.KeyPrefix(index))
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

// Queries

func (k Keeper) BallotByIdentifier(goCtx context.Context, req *types.QueryBallotByIdentifierRequest) (*types.QueryBallotByIdentifierResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	ballot, found := k.GetBallot(ctx, req.BallotIdentifier)
	if !found {
		return nil, status.Error(codes.NotFound, "not found ballot")
	}
	type voters struct {
		VoterAddress string `json:"voter_address"`
		VoteType     string `json:"vote_type"`
	}
	votersList := make([]voters, len(ballot.VoterList))
	for i, voterAddress := range ballot.VoterList {
		votersList[i].VoterAddress = voterAddress
		ballot.GetIndex()
		votersList[i].VoteType = ballot.Votes[ballot.GetVoterIndex(voterAddress)].String()
	}

	outputString := types2.PrettyPrintStruct(votersList)
	ballot.VoterList = []string{outputString}
	ballot.Votes = nil
	return &types.QueryBallotByIdentifierResponse{Ballot: &ballot}, nil
}
