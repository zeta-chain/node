package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// RemoveInboundTracker removes the inbound tracker if it exists.
func (k msgServer) RemoveInboundTracker(
	goCtx context.Context,
	msg *types.MsgRemoveInboundTracker,
) (*types.MsgRemoveInboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	k.RemoveInboundTrackerIfExists(ctx, msg.ChainId, msg.TxHash)
	return &types.MsgRemoveInboundTrackerResponse{}, nil
}
