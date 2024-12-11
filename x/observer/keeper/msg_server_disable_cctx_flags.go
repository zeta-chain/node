package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// DisableCCTX disables the IsInboundEnabled and IsOutboundEnabled flags. These flags control the creation of inbounds and outbounds.
// The flags are disabled by the policy account with the groupEmergency policy type.
func (k msgServer) DisableCCTX(
	goCtx context.Context,
	msg *types.MsgDisableCCTX,
) (*types.MsgDisableCCTXResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}
	// check if the value exists,
	// if not, set the default value for the Inbound and Outbound flags only
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
		flags.GasPriceIncreaseFlags = nil
	}

	if msg.DisableInbound {
		flags.IsInboundEnabled = false
	}
	if msg.DisableOutbound {
		flags.IsOutboundEnabled = false
	}
	k.SetCrosschainFlags(ctx, flags)

	err = ctx.EventManager().EmitTypedEvents(&types.EventCCTXDisabled{
		MsgTypeUrl:        sdk.MsgTypeURL(&types.MsgDisableCCTX{}),
		IsInboundEnabled:  flags.IsInboundEnabled,
		IsOutboundEnabled: flags.IsOutboundEnabled,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting event EventCCTXDisabled :", err)
	}

	return &types.MsgDisableCCTXResponse{}, nil
}
