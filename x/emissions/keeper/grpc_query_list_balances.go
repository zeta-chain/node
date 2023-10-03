package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/emissions/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListPoolAddresses(_ context.Context, req *types.QueryListPoolAddressesRequest) (*types.QueryListPoolAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	return &types.QueryListPoolAddressesResponse{UndistributedObserverBalancesAddress: types.UndistributedObserverRewardsPoolAddress.String(),
		EmissionModuleAddress:           types.EmissionsModuleAddress.String(),
		UndistributedTssBalancesAddress: types.UndistributedTssRewardsPoolAddress.String()}, nil
}
