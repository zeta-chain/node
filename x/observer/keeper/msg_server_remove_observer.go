package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// RemoveObserver removes an observer address from the observer set
func (k msgServer) RemoveObserver(
	goCtx context.Context,
	msg *types.MsgRemoveObserver,
) (*types.MsgRemoveObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.RemoveNodeAccount(ctx, msg.ObserverAddress)
	k.RemoveObserverFromSet(ctx, msg.ObserverAddress)
	k.DecrementLastObserverCount(ctx)

	return &types.MsgRemoveObserverResponse{}, nil
}
