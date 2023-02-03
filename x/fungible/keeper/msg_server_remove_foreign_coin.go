package keeper

import (
	"context"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k msgServer) RemoveForeignCoin(goCtx context.Context, msg *types.MsgRemoveForeignCoin) (*types.MsgRemoveForeignCoinResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//
	//if msg.Creator != types.AdminAddress {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only admin can remove foreign coin")
	//}
	//index := msg.Name

	//_, found := k.GetForeignCoins(ctx, index)
	//if !found {
	//	return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "foreign coin not found")
	//}
	//k.RemoveForeignCoins(ctx, index)
	// TODO : FIX THIS MSG
	return &types.MsgRemoveForeignCoinResponse{}, nil
}
