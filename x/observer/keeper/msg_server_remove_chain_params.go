package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// RemoveChainParams removes chain parameters for a specific chain.
func (k msgServer) RemoveChainParams(goCtx context.Context, msg *types.MsgRemoveChainParams) (*types.MsgRemoveChainParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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
