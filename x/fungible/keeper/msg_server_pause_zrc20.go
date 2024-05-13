package keeper

import (
	"context"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// PauseZRC20 pauses the ZRC20 token
func (k msgServer) PauseZRC20(
	goCtx context.Context,
	msg *types.MsgPauseZRC20,
) (*types.MsgPauseZRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check message validity
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// iterate all foreign coins and set paused status
	for _, zrc20 := range msg.Zrc20Addresses {
		fc, found := k.GetForeignCoins(ctx, zrc20)
		if !found {
			return nil, cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "foreign coin not found %s", zrc20)
		}
		// Set status to paused
		fc.Paused = true
		k.SetForeignCoins(ctx, fc)
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventZRC20Paused{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgPauseZRC20{}),
			Zrc20Addresses: msg.Zrc20Addresses,
			Signer:         msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventZRC20Paused",
			"error", err.Error(),
		)
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgPauseZRC20Response{}, nil
}
