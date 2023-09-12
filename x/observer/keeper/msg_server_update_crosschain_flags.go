package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateCrosschainFlags updates the crosschain related flags.
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) UpdateCrosschainFlags(goCtx context.Context, msg *types.MsgUpdateCrosschainFlags) (*types.MsgUpdateCrosschainFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_stop_inbound_cctx) {
		return &types.MsgUpdateCrosschainFlagsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// check if the value exists
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	flags.IsInboundEnabled = msg.IsInboundEnabled
	flags.IsOutboundEnabled = msg.IsOutboundEnabled

	if msg.GasPriceIncreaseFlags != nil {
		flags.GasPriceIncreaseFlags = msg.GasPriceIncreaseFlags
	}

	k.SetCrosschainFlags(ctx, flags)

	return &types.MsgUpdateCrosschainFlagsResponse{}, nil
}
