package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// RemoveOutboundTracker removes a record from the outbound transaction tracker by chain ID and nonce.
//
// Authorized: admin policy group 1.
func (k msgServer) RemoveOutboundTracker(goCtx context.Context, msg *types.MsgRemoveOutboundTracker) (*types.MsgRemoveOutboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency) {
		return &types.MsgRemoveOutboundTrackerResponse{}, authoritytypes.ErrUnauthorized
	}

	k.RemoveOutboundTrackerFromStore(ctx, msg.ChainId, msg.Nonce)
	return &types.MsgRemoveOutboundTrackerResponse{}, nil
}
