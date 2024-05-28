package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) GetEmissionsFactors(
	goCtx context.Context,
	_ *types.QueryGetEmissionsFactorsRequest,
) (*types.QueryGetEmissionsFactorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, found := k.GetParams(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "params not found")
	}
	reservesFactor, bondFactor, durationFactor := k.GetBlockRewardComponents(ctx, params)
	return &types.QueryGetEmissionsFactorsResponse{
		ReservesFactor: reservesFactor.String(),
		BondFactor:     bondFactor.String(),
		DurationFactor: durationFactor.String(),
	}, nil
}
