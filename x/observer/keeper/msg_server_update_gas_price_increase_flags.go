package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) UpdateGasPriceIncreaseFlags(goCtx context.Context, msg *types.MsgUpdateGasPriceIncreaseFlags) (*types.MsgUpdateGasPriceIncreaseFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg) {
		return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, authoritytypes.ErrUnauthorized
	}

	// check if the value exists,
	// This will also set the default values for such as GasPriceIncreaseFlags
	// We can still use the default values as all flags are part of the same struct ,
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	err := msg.GasPriceIncreaseFlags.Validate()
	if err != nil {
		return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, err
	}

	flags.GasPriceIncreaseFlags = &msg.GasPriceIncreaseFlags
	k.SetCrosschainFlags(ctx, flags)

	err = ctx.EventManager().EmitTypedEvents(&types.EventCrosschainFlagsUpdated{
		MsgTypeUrl:                   sdk.MsgTypeURL(&types.MsgUpdateGasPriceIncreaseFlags{}),
		IsInboundEnabled:             flags.IsInboundEnabled,
		IsOutboundEnabled:            flags.IsOutboundEnabled,
		GasPriceIncreaseFlags:        flags.GasPriceIncreaseFlags,
		BlockHeaderVerificationFlags: flags.BlockHeaderVerificationFlags,
		Signer:                       msg.Creator,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting event EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, nil
}
