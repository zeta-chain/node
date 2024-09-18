package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// UpdateChainInfo updates the chain info object
// If the provided chain does not exist in the chain info object, it is added
// If the chain already exists in the chain info object, it is updated
func (k msgServer) UpdateChainInfo(
	goCtx context.Context,
	msg *types.MsgUpdateChainInfo,
) (*types.MsgUpdateChainInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// This message is only allowed to be called by group admin
	// Group admin because this functionality would rarely be called
	// and overwriting false chain info can have undesired effects
	err := k.CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnauthorized, err.Error())
	}

	chainInfo := types.ChainInfo{}
	chainInfoExist := false

	existingChainInfo, found := k.GetChainInfo(ctx)
	if found {
		chainInfo = existingChainInfo
	}
	// try to update a chain if the chain info already exists
	for i, chain := range chainInfo.Chains {
		if chain.ChainId == msg.Chain.ChainId {
			chainInfo.Chains[i] = msg.Chain
			chainInfoExist = true
		}
	}

	// if the chain info does not exist, add the chain to the chain info object
	if !chainInfoExist {
		chainInfo.Chains = append(chainInfo.Chains, msg.Chain)
	}

	k.SetChainInfo(ctx, chainInfo)
	return &types.MsgUpdateChainInfoResponse{}, nil
}
