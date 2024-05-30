package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateGasPriceIncreaseFlags updates the GasPriceIncreaseFlags. These flags control the increase of gas prices.
// The flags are updated by the policy account with the groupOperational policy type.
func (k msgServer) UpdateGasPriceIncreaseFlags(
	goCtx context.Context,
	msg *types.MsgUpdateGasPriceIncreaseFlags,
) (*types.MsgUpdateGasPriceIncreaseFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupOperational) {
		return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, authoritytypes.ErrUnauthorized.Wrap(
			"UpdateGasPriceIncreaseFlags can only be executed by the correct policy account",
		)
	}
	// check if the value exists,
	// if not, set the default value for the GasPriceIncreaseFlags only
	// Set Inbound and Outbound flags to false
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
		flags.IsInboundEnabled = false
		flags.IsOutboundEnabled = false
	}

	err := msg.GasPriceIncreaseFlags.Validate()
	if err != nil {
		return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, err
	}

	flags.GasPriceIncreaseFlags = &msg.GasPriceIncreaseFlags
	k.SetCrosschainFlags(ctx, flags)

	err = ctx.EventManager().EmitTypedEvents(&types.EventGasPriceIncreaseFlagsUpdated{
		MsgTypeUrl:            sdk.MsgTypeURL(&types.MsgUpdateGasPriceIncreaseFlags{}),
		GasPriceIncreaseFlags: flags.GasPriceIncreaseFlags,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting EventGasPriceIncreaseFlagsUpdated :", err)
	}

	return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, nil
}
