package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// RemoveObserver removes an observer address from the observer set
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

	k.RemoveNodeAccount(ctx, msg.ObserverAddress)
	k.RemoveObserverFromSet(ctx, msg.ObserverAddress)
	k.DecrementLastObserverCount(ctx)

	return &types.MsgRemoveObserverResponse{}, nil
}
