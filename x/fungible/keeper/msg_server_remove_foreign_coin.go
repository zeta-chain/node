package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// RemoveForeignCoin removes a coin from the list of foreign coins in the module's state.
//
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) RemoveForeignCoin(goCtx context.Context, msg *types.MsgRemoveForeignCoin) (*types.MsgRemoveForeignCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_group2) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Removal can only be executed by the correct policy account")
	}
	index := msg.Name
	_, found := k.GetForeignCoins(ctx, index)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "foreign coin not found")
	}
	k.RemoveForeignCoins(ctx, index)
	return &types.MsgRemoveForeignCoinResponse{}, nil
}
