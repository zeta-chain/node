package keeper

import (
	"context"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) ShowAvailableEmissions(goCtx context.Context, req *types.QueryShowAvailableEmissionsRequest) (*types.QueryShowAvailableEmissionsResponse, error) {
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
