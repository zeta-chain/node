package keeper

import (
	"context"

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
	for _, cp := range chainParamsList.ChainParams {
		if cp.ConfirmationParams == nil {
			return nil, types.ErrInvalidChainParams
		}

		// setting fast confirmation count to same value as safe confirmation count
		// will effectively disable fast confirmation
		if cp.ChainId == msg.ChainId {
			foundChain = true
			cp.ConfirmationParams.FastInboundCount = cp.ConfirmationParams.SafeInboundCount
			cp.ConfirmationParams.FastOutboundCount = cp.ConfirmationParams.SafeOutboundCount
		}
	}
	if !foundChain {
		return nil, cosmoserror.Wrap(types.ErrChainParamsNotFound, "no matching chain ID found")
	}

	// set the updated chain params list
	k.SetChainParamsList(ctx, chainParamsList)

	return &types.MsgDisableFastConfirmationResponse{}, nil
}
