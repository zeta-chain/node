package keeper

import (
	"context"
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateZRC20PausedStatus updates the paused status of a ZRC20
// The list of ZRC20s are either paused or unpaused
func (k Keeper) UpdateZRC20PausedStatus(
	goCtx context.Context,
	msg *types.MsgUpdateZRC20PausedStatus,
) (*types.MsgUpdateZRC20PausedStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check message validity
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// check if the sender is the admin
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
		return nil, cosmoserrors.Wrap(sdkerrors.ErrUnauthorized, "Update can only be executed by the correct policy account")
	}

	pausedStatus := true
	if msg.Action == types.UpdatePausedStatusAction_UNPAUSE {
		pausedStatus = false
	}

	// iterate all foreign coins
	for _, zrc20 := range msg.Zrc20Addresses {
		fc, found := k.GetForeignCoins(ctx, zrc20)
		if !found {
			return nil, cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "foreign coin not found %s", zrc20)
		}

		fc.Paused = pausedStatus
		k.SetForeignCoins(ctx, fc)
	}

	err := ctx.EventManager().EmitTypedEvent(
		&types.EventZRC20PausedStatusUpdated{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgUpdateZRC20PausedStatus{}),
			Action:         msg.Action,
			Zrc20Addresses: msg.Zrc20Addresses,
			Signer:         msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventZRC20PausedStatusUpdated",
			"error", err.Error(),
		)
		return nil, sdkerrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgUpdateZRC20PausedStatusResponse{}, nil
}
