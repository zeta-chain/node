package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) DisableCCTXFlags(goCtx context.Context, msg *types.MsgDisableCCTXFlags) (*types.MsgDisableCCTXFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg) {
		return &types.MsgDisableCCTXFlagsResponse{}, authoritytypes.ErrUnauthorized
	}

	// check if the value exists,
	// This will also set the default values for such as GasPriceIncreaseFlags
	// We can still use the default values as all flags are part of the same struct ,
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	if msg.DisableInbound {
		flags.IsInboundEnabled = false
	}
	if msg.DisableOutbound {
		flags.IsOutboundEnabled = false
	}

	k.SetCrosschainFlags(ctx, flags)

	err := ctx.EventManager().EmitTypedEvents(&types.EventCrosschainFlagsUpdated{
		MsgTypeUrl:                   sdk.MsgTypeURL(&types.MsgDisableCCTXFlags{}),
		IsInboundEnabled:             flags.IsInboundEnabled,
		IsOutboundEnabled:            flags.IsOutboundEnabled,
		GasPriceIncreaseFlags:        flags.GasPriceIncreaseFlags,
		BlockHeaderVerificationFlags: flags.BlockHeaderVerificationFlags,
		Signer:                       msg.Creator,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting event EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgDisableCCTXFlagsResponse{}, nil
}
