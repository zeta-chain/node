package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) UpdateClientParams(goCtx context.Context, msg *types.MsgUpdateClientParams) (*types.MsgUpdateClientParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_update_client_params) {
	//	return &types.MsgUpdateClientParamsResponse{}, types.ErrNotAuthorizedPolicy
	//}
	if !k.GetParams(ctx).IsChainIDSupported(msg.ChainId) {
		return &types.MsgUpdateClientParamsResponse{}, types.ErrSupportedChains
	}
	k.SetClientParamsByChainID(ctx, msg.ChainId, *msg.ClientParams)
	return &types.MsgUpdateClientParamsResponse{}, nil
}
