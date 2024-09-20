package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// RemoveChainInfo removes the chain info for the specified chain id
func (k msgServer) RemoveChainInfo(
	goCtx context.Context,
	msg *types.MsgRemoveChainInfo,
) (*types.MsgRemoveChainInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check the authorization for this message against the authorization list
	err := k.CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnauthorized, err.Error())
	}

	chainInfo, found := k.GetChainInfo(ctx)
	if !found {
		return nil, types.ErrChainInfoNotFound
	}

	updatedChainInfo := RemoveChain(chainInfo, msg.ChainId)
	k.SetChainInfo(ctx, updatedChainInfo)
	return &types.MsgRemoveChainInfoResponse{}, nil
}

// RemoveChain removes the chain with the given chainID from the chainInfo
// If the chain does not exist, the chainInfo is returned as is
func RemoveChain(chainInfo types.ChainInfo, chainID int64) types.ChainInfo {
	updatedChainInfo := types.ChainInfo{}
	for _, chain := range chainInfo.Chains {
		if chain.ChainId != chainID {
			updatedChainInfo.Chains = append(updatedChainInfo.Chains, chain)
		}
	}
	return updatedChainInfo
}
