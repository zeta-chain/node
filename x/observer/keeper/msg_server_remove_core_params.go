package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// RemoveCoreParams removes core parameters for a specific chain.
func (k msgServer) RemoveCoreParams(goCtx context.Context, msg *types.MsgRemoveCoreParams) (*types.MsgRemoveCoreParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_group2) {
		return &types.MsgRemoveCoreParamsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// find current core params list or initialize a new one
	coreParamsList, found := k.GetCoreParamsList(ctx)
	if !found {
		return &types.MsgRemoveCoreParamsResponse{}, types.ErrCoreParamsNotFound
	}

	// remove the core param from the list
	newCoreParamsList := types.CoreParamsList{}
	found = false
	for _, cp := range coreParamsList.CoreParams {
		if cp.ChainId != msg.ChainId {
			newCoreParamsList.CoreParams = append(newCoreParamsList.CoreParams, cp)
		} else {
			found = true
		}
	}
	if !found {
		return &types.MsgRemoveCoreParamsResponse{}, types.ErrCoreParamsNotFound
	}

	k.SetCoreParamsList(ctx, newCoreParamsList)
	return &types.MsgRemoveCoreParamsResponse{}, nil
}
