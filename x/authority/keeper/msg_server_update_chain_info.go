package keeper

import (
	"context"
	"fmt"

	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

// UpdateChainInfo updates the chain inffo structure that adds new static chain info or overwrite existing chain info
// on the hard-coded chain info
func (k msgServer) UpdateChainInfo(
	goCtx context.Context,
	msg *types.MsgUpdateChainInfo,
) (*types.MsgUpdateChainInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// This message is only allowed to be called by group admin
	// Group admin because this functionality would rarely be called
	// and overwriting false chain info can have undesired effects
	if !k.IsAuthorized(ctx, msg.Signer, types.PolicyType_groupAdmin) {
		return nil, cosmoserror.Wrap(types.ErrUnauthorized, fmt.Sprintf("Signer %s", msg.Signer))
	}

	// set chain info
	k.SetChainInfo(ctx, msg.ChainInfo)

	return &types.MsgUpdateChainInfoResponse{}, nil
}
