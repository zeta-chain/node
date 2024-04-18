package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// UpdateRateLimiterFlags updates the rate limiter flags.
// Authorized: admin policy operational.
func (k msgServer) UpdateRateLimiterFlags(goCtx context.Context, msg *types.MsgUpdateRateLimiterFlags) (*types.MsgUpdateRateLimiterFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupOperational) {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.SetRateLimiterFlags(ctx, msg.RateLimiterFlags)

	return &types.MsgUpdateRateLimiterFlagsResponse{}, nil
}
