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

func (k Keeper) SetVoter(ctx sdk.Context, voter types.Voter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	voter.Index = fmt.Sprintf("%s", voter.VoteIdentifier)
	b := k.cdc.MustMarshal(&voter)
	store.Set([]byte(voter.Index), b)
}

func (k Keeper) GetVoter(ctx sdk.Context, index string) (val types.Voter, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// Queries

func (k Keeper) VoterByIdentifier(goCtx context.Context, req *types.QueryVoterByIdentifierRequest) (*types.QueryVoterByIdentifierResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	voter, isFound := k.GetVoter(ctx, req.VoteIdentifier)
	if !isFound {
		return &types.QueryVoterByIdentifierResponse{Voter: "Not Found"}, nil
	}
	return &types.QueryVoterByIdentifierResponse{Voter: voter.String()}, nil
}
