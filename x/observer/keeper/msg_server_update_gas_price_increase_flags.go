package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) UpdateGasPriceIncreaseFlags(goCtx context.Context, msg *types.MsgUpdateGasPriceIncreaseFlags) (*types.MsgUpdateGasPriceIncreaseFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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

	err = msg.GasPriceIncreaseFlags.Validate()
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
		ctx.Logger().Error("Error emitting event EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgUpdateGasPriceIncreaseFlagsResponse{}, nil
}
