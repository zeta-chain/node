package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// UpdateOperationalChainParams updates the operational-related chain params
// Unlike MsgUpdateChainParams, this message doesn't allow updated sensitive values such as the gateway contract to listen to on connected chains
func (k msgServer) UpdateOperationalChainParams(
	goCtx context.Context,
	msg *types.MsgUpdateOperationalChainParams,
) (*types.MsgUpdateOperationalChainParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// can't only update existing params, not create the object
	chainParamsList, found := k.GetChainParamsList(ctx)
	if !found {
		return nil, errors.Wrap(types.ErrChainParamsNotFound, "chain params list not found")
	}

	// find chain params for the chain
	for i, cp := range chainParamsList.ChainParams {
		if cp.ChainId == msg.ChainId {
			// update values and save object
			chainParamsList.ChainParams[i].GasPriceTicker = msg.GasPriceTicker
			chainParamsList.ChainParams[i].InboundTicker = msg.InboundTicker
			chainParamsList.ChainParams[i].OutboundTicker = msg.OutboundTicker
			chainParamsList.ChainParams[i].WatchUtxoTicker = msg.WatchUtxoTicker
			chainParamsList.ChainParams[i].OutboundScheduleInterval = msg.OutboundScheduleInterval
			chainParamsList.ChainParams[i].OutboundScheduleLookahead = msg.OutboundScheduleLookahead
			chainParamsList.ChainParams[i].ConfirmationParams = &msg.ConfirmationParams
			k.SetChainParamsList(ctx, chainParamsList)

			return &types.MsgUpdateOperationalChainParamsResponse{}, nil
		}
	}

	return nil, errors.Wrap(types.ErrChainParamsNotFound, "chain params not found")
}
