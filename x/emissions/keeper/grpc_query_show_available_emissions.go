package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ShowAvailableEmissions(
	goCtx context.Context,
	req *types.QueryShowAvailableEmissionsRequest,
) (*types.QueryShowAvailableEmissionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	emissions, found := k.GetWithdrawableEmission(ctx, req.Address)
	if !found {
		return &types.QueryShowAvailableEmissionsResponse{
			Amount: sdk.NewCoin(config.BaseDenom, sdk.ZeroInt()).String(),
		}, nil
	}
	return &types.QueryShowAvailableEmissionsResponse{
		Amount: sdk.NewCoin(config.BaseDenom, emissions.Amount).String(),
	}, nil
}
