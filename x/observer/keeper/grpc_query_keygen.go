package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Keygen(c context.Context, _ *types.QueryGetKeygenRequest) (*types.QueryGetKeygenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, found := k.GetKeygen(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}
	return &types.QueryGetKeygenResponse{Keygen: &val}, nil
}
