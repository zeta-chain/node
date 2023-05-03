package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) UpdateCoreParams(goCtx context.Context, msg *types.MsgUpdateCoreParams) (*types.MsgUpdateCoreParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_update_client_params) {
		return &types.MsgUpdateCoreParamsResponse{}, types.ErrNotAuthorizedPolicy
	}
	if !k.GetParams(ctx).IsChainIDSupported(msg.CoreParams.ChainId) {
		return &types.MsgUpdateCoreParamsResponse{}, types.ErrSupportedChains
	}
	coreParams, found := k.GetAllCoreParams(ctx)
	if !found {
		return &types.MsgUpdateCoreParamsResponse{}, types.ErrCoreParamsNotSet
	}
	newCoreParams := make([]*types.CoreParams, len(coreParams.CoreParams))
	for i, cp := range coreParams.CoreParams {
		if cp.ChainId == msg.CoreParams.ChainId {
			newCoreParams[i] = msg.CoreParams
			continue
		}
		newCoreParams[i] = cp
	}
	k.SetCoreParams(ctx, types.CoreParamsList{CoreParams: newCoreParams})
	return &types.MsgUpdateCoreParamsResponse{}, nil
}
