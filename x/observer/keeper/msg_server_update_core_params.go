package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) UpdateCoreParams(goCtx context.Context, msg *types.MsgUpdateCoreParams) (*types.MsgUpdateCoreParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_update_client_params) {
	//	return &types.MsgUpdateClientParamsResponse{}, types.ErrNotAuthorizedPolicy
	//}
	if !k.GetParams(ctx).IsChainIDSupported(msg.ChainId) {
		return &types.MsgUpdateCoreParamsResponse{}, types.ErrSupportedChains
	}
	k.SetCoreParamsByChainID(ctx, msg.ChainId, *msg.CoreParams)
	return &types.MsgUpdateCoreParamsResponse{}, nil
}
