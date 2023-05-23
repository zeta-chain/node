package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) RemoveForeignCoin(goCtx context.Context, msg *types.MsgRemoveForeignCoin) (*types.MsgRemoveForeignCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaobserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
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
