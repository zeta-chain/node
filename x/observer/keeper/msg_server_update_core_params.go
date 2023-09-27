package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateCoreParams updates core parameters for a specific chain. Core parameters include
// confirmation count, outbound transaction schedule interval, ZETA token,
// connector and ERC20 custody contract addresses, etc.
//
// Throws an error if the chain ID is not supported.
//
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) UpdateCoreParams(goCtx context.Context, msg *types.MsgUpdateCoreParams) (*types.MsgUpdateCoreParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_group2) {
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
