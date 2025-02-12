package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// DisableFastConfirmation disables fast confirmation for the given chain IDs
// Inbound and outbound will be only confirmed using SAFE confirmation count on disabled chains
func (k msgServer) DisableFastConfirmation(
	goCtx context.Context,
	msg *types.MsgDisableFastConfirmation,
) (*types.MsgDisableFastConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// find current chain params list
	chainParamsList, found := k.GetChainParamsList(ctx)
	if !found {
		return &types.MsgDisableFastConfirmationResponse{}, types.ErrChainParamsNotFound
	}

	// create a lookup map for the chain IDs
	disableAll := len(msg.ChainIdList) == 0
	chainIDMap := make(map[int64]bool)
	for _, chainID := range msg.ChainIdList {
		chainIDMap[chainID] = true
	}

	// disable fast confirmation by setting fast confirmation count to zero
	for _, cp := range chainParamsList.ChainParams {
		if cp.ConfirmationParams == nil {
			continue // should never happen
		}

		if disableAll || chainIDMap[cp.ChainId] {
			cp.ConfirmationParams.FastInboundCount = 0
			cp.ConfirmationParams.FastOutboundCount = 0
		}
	}

	// validate the updated chain params list
	if err := chainParamsList.Validate(); err != nil {
		return &types.MsgDisableFastConfirmationResponse{}, errorsmod.Wrap(types.ErrInvalidChainParams, err.Error())
	}

	// set the updated chain params list
	k.SetChainParamsList(ctx, chainParamsList)

	return &types.MsgDisableFastConfirmationResponse{}, nil
}
