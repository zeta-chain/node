package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) GetEmmisonsFactors(goCtx context.Context, _ *types.QueryGetEmmisonsFactorsRequest) (*types.QueryGetEmmisonsFactorsResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	reservesFactor, bondFactor, durationFactor := k.GetBlockRewardComponents(ctx)
	return &types.QueryGetEmmisonsFactorsResponse{
		ReservesFactor: reservesFactor.String(),
		BondFactor:     bondFactor.String(),
		DurationFactor: durationFactor.String(),
	}, nil
}
