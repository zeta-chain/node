package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateCrosschainFlags updates the crosschain related flags.
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) UpdateCrosschainFlags(goCtx context.Context, msg *types.MsgUpdateCrosschainFlags) (*types.MsgUpdateCrosschainFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_stop_inbound_cctx) {
		return &types.MsgUpdateCrosschainFlagsResponse{}, types.ErrNotAuthorizedPolicy
	}
	// Check if the value exists
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}
	flags.IsInboundEnabled = msg.IsInboundEnabled
	flags.IsOutboundEnabled = msg.IsOutboundEnabled
	k.SetCrosschainFlags(ctx, flags)

	return &types.MsgUpdateCrosschainFlagsResponse{}, nil
}
