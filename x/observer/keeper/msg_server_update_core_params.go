package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateCoreParams updates core parameters for a specific chain, or add a new one.
// Core parameters include: confirmation count, outbound transaction schedule interval, ZETA token,
// connector and ERC20 custody contract addresses, etc.
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) UpdateCoreParams(goCtx context.Context, msg *types.MsgUpdateCoreParams) (*types.MsgUpdateCoreParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_group2) {
		return &types.MsgUpdateCoreParamsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// find current core params list or initialize a new one
	coreParamsList, found := k.GetCoreParamsList(ctx)
	if !found {
		coreParamsList = types.CoreParamsList{}
	}

	// find core params for the chain
	for i, cp := range coreParamsList.CoreParams {
		if cp.ChainId == msg.CoreParams.ChainId {
			coreParamsList.CoreParams[i] = msg.CoreParams
			k.SetCoreParamsList(ctx, coreParamsList)
			return &types.MsgUpdateCoreParamsResponse{}, nil
		}
	}

	// add new core params
	coreParamsList.CoreParams = append(coreParamsList.CoreParams, msg.CoreParams)
	k.SetCoreParamsList(ctx, coreParamsList)

	return &types.MsgUpdateCoreParamsResponse{}, nil
}
