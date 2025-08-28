package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// EnableCCTX enables the IsInboundEnabled and IsOutboundEnabled flags.These flags control the creation of inbounds and outbounds.
// The flags are enabled by the policy account with the groupOperational policy type.
func (k msgServer) EnableCCTX(
	goCtx context.Context,
	msg *types.MsgEnableCCTX,
) (*types.MsgEnableCCTXResponse, error) {
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

	if msg.EnableInbound {
		flags.IsInboundEnabled = true
	}
	if msg.EnableOutbound {
		flags.IsOutboundEnabled = true
	}

	k.SetCrosschainFlags(ctx, flags)

	err = ctx.EventManager().EmitTypedEvents(&types.EventCCTXEnabled{
		MsgTypeUrl:        sdk.MsgTypeURL(&types.MsgEnableCCTX{}),
		IsInboundEnabled:  flags.IsInboundEnabled,
		IsOutboundEnabled: flags.IsOutboundEnabled,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting EventCCTXEnabled :", err)
	}

	return &types.MsgEnableCCTXResponse{}, nil
}
