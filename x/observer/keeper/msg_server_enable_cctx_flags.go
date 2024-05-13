package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) EnableCCTXFlags(goCtx context.Context, msg *types.MsgEnableCCTXFlags) (*types.MsgEnableCCTXFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check if the value exists,
	// if not, set the default value for the Inbound and Outbound flags only
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
		flags.GasPriceIncreaseFlags = nil
		flags.BlockHeaderVerificationFlags = nil
	}

	if msg.EnableInbound {
		flags.IsInboundEnabled = true
	}
	if msg.EnableOutbound {
		flags.IsOutboundEnabled = true
	}

	k.SetCrosschainFlags(ctx, flags)

	err = ctx.EventManager().EmitTypedEvents(&types.EventCCTXFlagsEnabled{
		MsgTypeUrl:        sdk.MsgTypeURL(&types.MsgEnableCCTXFlags{}),
		IsInboundEnabled:  flags.IsInboundEnabled,
		IsOutboundEnabled: flags.IsOutboundEnabled,
	})

	if err != nil {
		ctx.Logger().Error("Error emitting event EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgEnableCCTXFlagsResponse{}, nil
}
