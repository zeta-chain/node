package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) SupportedChains(goCtx context.Context, _ *types.QuerySupportedChains) (*types.QuerySupportedChainsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chains := k.GetParams(ctx).GetSupportedChains()
	return &types.QuerySupportedChainsResponse{Chains: chains}, nil
}
