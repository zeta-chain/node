package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// UpdateChainInfo updates the chain info structure that adds new static chain info or overwrite existing chain info
// on the hard-coded chain info
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
	// set chain info
	k.SetChainInfo(ctx, msg.ChainInfo)

	return &types.MsgUpdateChainInfoResponse{}, nil
}
