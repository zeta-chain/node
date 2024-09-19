package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// AddAuthorization defines a method to add an authorization.If the authorization already exists, it will be overwritten with the provided policy.
// This should be called by the admin policy account.
func (k msgServer) AddAuthorization(
	goCtx context.Context,
	msg *types.MsgAddAuthorization,
) (*types.MsgAddAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if the caller is authorized to add an authorization
	err := k.CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnauthorized, err.Error())
	}

	authorizationList, found := k.GetAuthorizationList(ctx)
	if !found {
		authorizationList = types.AuthorizationList{Authorizations: []types.Authorization{}}
	}
	authorizationList.SetAuthorization(types.Authorization{MsgUrl: msg.MsgUrl, AuthorizedPolicy: msg.AuthorizedPolicy})

	// validate the authorization list after adding the authorization as a precautionary measure.
	err = authorizationList.Validate()
	if err != nil {
		return nil, errorsmod.Wrap(err, "authorization list is invalid")
	}

	k.SetAuthorizationList(ctx, authorizationList)
	return &types.MsgAddAuthorizationResponse{}, nil
}
