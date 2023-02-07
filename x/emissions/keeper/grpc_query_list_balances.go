package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListBalances(goCtx context.Context, req *types.QueryListBalancesRequest) (*types.QueryListBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	trackers := k.GetAllEmissionTrackers(ctx)
	return &types.QueryListBalancesResponse{Trackers: trackers, EmissionModuleAddress: types.EmissionsModuleAddress.String()}, nil
}
