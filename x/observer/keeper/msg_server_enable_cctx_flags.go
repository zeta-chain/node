package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) EnableCCTXFlags(goCtx context.Context, msg *types.MsgEnableCCTXFlags) (*types.MsgEnableCCTXFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg) {
		return &types.MsgEnableCCTXFlagsResponse{}, authoritytypes.ErrUnauthorized
	}

	// check if the value exists,
	// This will also set the default values for such as GasPriceIncreaseFlags
	// We can still use the default values as all flags are part of the same struct ,
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	if msg.EnableInbound {
		flags.IsInboundEnabled = true
	}
	if msg.EnableOutbound {
		flags.IsOutboundEnabled = true
	}

	k.SetCrosschainFlags(ctx, flags)

	err := ctx.EventManager().EmitTypedEvents(&types.EventCrosschainFlagsUpdated{
		MsgTypeUrl:                   sdk.MsgTypeURL(&types.MsgEnableCCTXFlags{}),
		IsInboundEnabled:             flags.IsInboundEnabled,
		IsOutboundEnabled:            flags.IsOutboundEnabled,
		GasPriceIncreaseFlags:        flags.GasPriceIncreaseFlags,
		BlockHeaderVerificationFlags: flags.BlockHeaderVerificationFlags,
		Signer:                       msg.Creator,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting event EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgEnableCCTXFlagsResponse{}, nil
}
