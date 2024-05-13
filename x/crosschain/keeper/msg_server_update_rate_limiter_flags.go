package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// UpdateRateLimiterFlags updates the rate limiter flags.
// Authorized: admin policy operational.
func (k msgServer) UpdateRateLimiterFlags(goCtx context.Context, msg *types.MsgUpdateRateLimiterFlags) (*types.MsgUpdateRateLimiterFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	k.SetRateLimiterFlags(ctx, msg.RateLimiterFlags)

	return &types.MsgUpdateRateLimiterFlagsResponse{}, nil
}
