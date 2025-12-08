package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// RemoveObserver removes an observer address from the observer set and node account list
func (k msgServer) RemoveObserver(
	goCtx context.Context,
	msg *types.MsgRemoveObserver,
) (*types.MsgRemoveObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// We remove it from both the node account list and the observer set to effectively observing and signing
	k.RemoveNodeAccount(ctx, msg.ObserverAddress)
	newCount := k.RemoveObserverFromSet(ctx, msg.ObserverAddress)
	k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: newCount})

	return &types.MsgRemoveObserverResponse{}, nil
}
