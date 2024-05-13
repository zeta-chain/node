package keeper

import (
	"context"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// UnpauseZRC20 unpauses the ZRC20 token
func (k msgServer) UnpauseZRC20(
	goCtx context.Context,
	msg *types.MsgUnpauseZRC20,
) (*types.MsgUnpauseZRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check message validity
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// iterate all foreign coins and set unpaused status
	for _, zrc20 := range msg.Zrc20Addresses {
		fc, found := k.GetForeignCoins(ctx, zrc20)
		if !found {
			return nil, cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "foreign coin not found %s", zrc20)
		}
		// Set status to unpaused
		fc.Paused = false
		k.SetForeignCoins(ctx, fc)
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventZRC20UnPaused{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgUnpauseZRC20{}),
			Zrc20Addresses: msg.Zrc20Addresses,
			Signer:         msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventZRC20UnPaused",
			"error", err.Error(),
		)
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgUnpauseZRC20Response{}, nil
}
