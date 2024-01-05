package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// RemoveFromOutTxTracker removes a record from the outbound transaction tracker by chain ID and nonce.
//
// Authorized: admin policy group 1.
func (k msgServer) RemoveFromOutTxTracker(goCtx context.Context, msg *types.MsgRemoveFromOutTxTracker) (*types.MsgRemoveFromOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group1) {
		return &types.MsgRemoveFromOutTxTrackerResponse{}, observertypes.ErrNotAuthorizedPolicy
	}

	k.RemoveOutTxTracker(ctx, msg.ChainId, msg.Nonce)
	return &types.MsgRemoveFromOutTxTrackerResponse{}, nil
}
