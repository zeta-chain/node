package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetBallot(ctx sdk.Context, voter types.Ballot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	voter.Index = fmt.Sprintf("%s", voter.BallotThreshold)
	b := k.cdc.MustMarshal(&voter)
	store.Set([]byte(voter.Index), b)
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
	voter, isFound := k.GetBallot(ctx, req.BallotIdentifier)
	if !isFound {
		return &types.QueryBallotByIdentifierResponse{Ballot: "Not Found"}, nil
	}
	return &types.QueryBallotByIdentifierResponse{Ballot: voter.String()}, nil
}
