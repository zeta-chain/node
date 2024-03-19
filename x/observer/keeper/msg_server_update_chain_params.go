package keeper

import (
	"context"
	"github.com/zeta-chain/zetacore/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateChainParams updates chain parameters for a specific chain, or add a new one.
// Chain parameters include: confirmation count, outbound transaction schedule interval, ZETA token,
// connector and ERC20 custody contract addresses, etc.
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) UpdateChainParams(goCtx context.Context, msg *types.MsgUpdateChainParams) (*types.MsgUpdateChainParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupAdmin) {
		return &types.MsgUpdateChainParamsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// find current chain params list or initialize a new one
	chainParamsList, found := k.GetChainParamsList(ctx)
	if !found {
		chainParamsList = types.ChainParamsList{}
	}

	if msg.ChainParams.ChainId == common.SepoliaChain().ChainId {
		_, err := k.ResetChainNonces(ctx, &MsgResetChainNonces{
			Creator: msg.Creator,
			ChainID: msg.ChainParams.ChainId,
		})
		if err != nil {
			return nil, err
		}
		ctx.Logger().Info("ResetChainNonces", "ChainID", msg.ChainParams.ChainId)
	}

	// find chain params for the chain
	for i, cp := range chainParamsList.ChainParams {
		if cp.ChainId == msg.ChainParams.ChainId {
			chainParamsList.ChainParams[i] = msg.ChainParams
			k.SetChainParamsList(ctx, chainParamsList)
			return &types.MsgUpdateChainParamsResponse{}, nil
		}
	}

	// add new chain params
	chainParamsList.ChainParams = append(chainParamsList.ChainParams, msg.ChainParams)
	k.SetChainParamsList(ctx, chainParamsList)

	return &types.MsgUpdateChainParamsResponse{}, nil
}
