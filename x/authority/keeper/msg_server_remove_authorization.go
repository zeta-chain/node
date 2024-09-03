package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// RemoveAuthorization defines a method to remove an authorization.
// This should be called by the admin policy account.
func (k msgServer) RemoveAuthorization(
	goCtx context.Context,
	msg *types.MsgRemoveAuthorization,
) (*types.MsgRemoveAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if the caller is authorized to remove an authorization
	err := k.CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnauthorized, err.Error())
	}

	// check if the authorization list exists, we can return early if there is no list.
	authorizationList, found := k.GetAuthorizationList(ctx)
	if !found {
		return nil, types.ErrAuthorizationListNotFound
	}

	// check if the authorization exists, we can return early if the authorization does not exist.
	_, err = authorizationList.GetAuthorizedPolicy(msg.MsgUrl)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("msg url %s", msg.MsgUrl))
	}

	// remove the authorization
	authorizationList.RemoveAuthorization(msg.MsgUrl)

	// validate the authorization list after adding the authorization as a precautionary measure.
	err = authorizationList.Validate()
	if err != nil {
		return nil, errorsmod.Wrap(err, "authorization list is invalid")
	}
	k.SetAuthorizationList(ctx, authorizationList)

	return &types.MsgRemoveAuthorizationResponse{}, nil
}
