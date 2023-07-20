package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetBlame(ctx sdk.Context, blame *types.Blame) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	b := k.cdc.MustMarshal(blame)
	store.Set([]byte(blame.Index), b)
}

func (k Keeper) GetBlame(ctx sdk.Context, index string) (val types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// Query

func (k Keeper) BlameByIdentifier(goCtx context.Context, request *types.QueryBlameByIdentifierRequest) (*types.QueryBlameByIdentifierResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	blame, found := k.GetBlame(ctx, request.BlameIdentifier)
	if !found {
		return nil, status.Error(codes.NotFound, "blame info not found")
	}

	return &types.QueryBlameByIdentifierResponse{
		BlameInfo: &blame,
	}, nil
}
