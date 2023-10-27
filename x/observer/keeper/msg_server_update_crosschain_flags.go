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

	requiredGroup := types.Policy_Type_group1
	if msg.IsInboundEnabled || msg.IsOutboundEnabled || msg.GasPriceIncreaseFlags != nil {
		requiredGroup = types.Policy_Type_group2
	}

	// check permission
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(requiredGroup) {
		return &types.MsgUpdateCrosschainFlagsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// check if the value exists
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	// update values
	flags.IsInboundEnabled = msg.IsInboundEnabled
	flags.IsOutboundEnabled = msg.IsOutboundEnabled

	if msg.GasPriceIncreaseFlags != nil {
		flags.GasPriceIncreaseFlags = msg.GasPriceIncreaseFlags
	}

	if msg.BlockHeaderVerificationFlags != nil {
		flags.BlockHeaderVerificationFlags = msg.BlockHeaderVerificationFlags
	}

	k.SetCrosschainFlags(ctx, flags)

	err := ctx.EventManager().EmitTypedEvents(&types.EventCrosschainFlagsUpdated{
		MsgTypeUrl:                   sdk.MsgTypeURL(&types.MsgUpdateCrosschainFlags{}),
		IsInboundEnabled:             msg.IsInboundEnabled,
		IsOutboundEnabled:            msg.IsOutboundEnabled,
		GasPriceIncreaseFlags:        msg.GasPriceIncreaseFlags,
		BlockHeaderVerificationFlags: msg.BlockHeaderVerificationFlags,
		Signer:                       msg.Creator,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgUpdateCrosschainFlagsResponse{}, nil
}
