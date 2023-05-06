package keeper

import (
	"context"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k msgServer) UpdateKeygen(goCtx context.Context, msg *types.MsgUpdateKeygen) (*types.MsgUpdateKeygenResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_update_keygen_block) {
		return &types.MsgUpdateKeygenResponse{}, zetaObserverTypes.ErrNotAuthorizedPolicy
	}
	keygen, found := k.GetKeygen(ctx)
	if !found {
		return nil, types.ErrKeygenNotFound
	}
	if msg.Block <= (ctx.BlockHeight() + 10) {
		return nil, types.ErrKeygenBlockTooLow
	}
	keygen.BlockNumber = msg.Block
	keygen.Status = types.KeygenStatus_PendingKeygen
	k.SetKeygen(ctx, keygen)
	EmitEventKeyGenBlockUpdated(ctx, &keygen)
	return &types.MsgUpdateKeygenResponse{}, nil
}
