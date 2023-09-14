package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// Updates permissions. Currently, this is only used to enable/disable the
// inbound transactions.
//
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) UpdatePermissionFlags(goCtx context.Context, msg *types.MsgUpdatePermissionFlags) (*types.MsgUpdatePermissionFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_stop_inbound_cctx) {
		return &types.MsgUpdatePermissionFlagsResponse{}, types.ErrNotAuthorizedPolicy
	}
	// Check if the value exists
	flags, isFound := k.GetPermissionFlags(ctx)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}
	flags.IsInboundEnabled = msg.IsInboundEnabled
	flags.IsOutboundEnabled = msg.IsOutboundEnabled
	k.SetPermissionFlags(ctx, flags)

	return &types.MsgUpdatePermissionFlagsResponse{}, nil
}
