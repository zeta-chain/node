package keeper

import (
	"context"
	"fmt"

	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// DisableFastConfirmation disables fast confirmation for the given chain ID
// Inbound and outbound will be only confirmed using SAFE confirmation count on disabled chains
func (k msgServer) DisableFastConfirmation(
	goCtx context.Context,
	msg *types.MsgDisableFastConfirmation,
) (*types.MsgDisableFastConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserror.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// find current chain params list
	chainParamsList, found := k.GetChainParamsList(ctx)
	if !found {
		return nil, types.ErrChainParamsNotFound
	}

	// disable fast confirmation
	foundChain := false
	for i := range chainParamsList.ChainParams {
		if chainParamsList.ChainParams[i].ChainId == msg.ChainId {
			if chainParamsList.ChainParams[i].ConfirmationParams == nil {
				return nil, types.ErrInvalidChainParams
			}

			// setting fast confirmation count to same value as safe confirmation count
			// will effectively disable fast confirmation
			foundChain = true
			chainParamsList.ChainParams[i].ConfirmationParams.FastInboundCount = chainParamsList.ChainParams[i].ConfirmationParams.SafeInboundCount
			chainParamsList.ChainParams[i].ConfirmationParams.FastOutboundCount = chainParamsList.ChainParams[i].ConfirmationParams.SafeOutboundCount

			// set the updated chain params list
			k.SetChainParamsList(ctx, chainParamsList)

			// there should be only one chain with the same chain ID
			break
		}
	}
	if !foundChain {
		return nil, cosmoserror.Wrap(
			types.ErrChainParamsNotFound,
			fmt.Sprintf("no matching chain ID found (%d)", msg.ChainId),
		)
	}

	return &types.MsgDisableFastConfirmationResponse{}, nil
}
