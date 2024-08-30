package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// RemoveOutboundTracker removes a record from the outbound transaction tracker by chain ID and nonce.
//
// Authorized: admin policy group 1.
func (k msgServer) RemoveOutboundTracker(
	goCtx context.Context,
	msg *types.MsgRemoveOutboundTracker,
) (*types.MsgRemoveOutboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	k.RemoveOutboundTrackerFromStore(ctx, msg.ChainId, msg.Nonce)
	return &types.MsgRemoveOutboundTrackerResponse{}, nil
}
