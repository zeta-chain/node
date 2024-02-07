package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// RemoveChainParams removes chain parameters for a specific chain.
func (k msgServer) RemoveChainParams(goCtx context.Context, msg *types.MsgRemoveChainParams) (*types.MsgRemoveChainParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_group2) {
		return &types.MsgRemoveChainParamsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// find current core params list or initialize a new one
	chainParamsList, found := k.GetChainParamsList(ctx)
	if !found {
		return &types.MsgRemoveChainParamsResponse{}, types.ErrChainParamsNotFound
	}

	// remove the core param from the list
	newChainParamsList := types.ChainParamsList{}
	found = false
	for _, cp := range chainParamsList.ChainParams {
		if cp.ChainId != msg.ChainId {
			newChainParamsList.ChainParams = append(newChainParamsList.ChainParams, cp)
		} else {
			found = true
		}
	}
	if !found {
		return &types.MsgRemoveChainParamsResponse{}, types.ErrChainParamsNotFound
	}

	k.SetChainParamsList(ctx, newChainParamsList)
	return &types.MsgRemoveChainParamsResponse{}, nil
}
