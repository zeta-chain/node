package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

// TODO : Refactor this file to authorization_list.go

// SetAuthorizationList sets the authorization list to the store.It returns an error if the list is invalid.
func (k Keeper) SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuthorizationListKey))
	b := k.cdc.MustMarshal(&list)
	store.Set([]byte{0}, b)
}

// GetAuthorizationList returns the authorization list from the store
func (k Keeper) GetAuthorizationList(ctx sdk.Context) (val types.AuthorizationList, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuthorizationListKey))
	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// CheckAuthorization checks if the signer is authorized to sign the message
func (k Keeper) CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error {
	// Policy transactions must have only one signer
	if len(msg.GetSigners()) != 1 {
		return errors.Wrap(types.ErrSigners, fmt.Sprintf("msg: %v", sdk.MsgTypeURL(msg)))
	}
	signer := msg.GetSigners()[0].String()
	msgURL := sdk.MsgTypeURL(msg)
	authorizationsList, found := k.GetAuthorizationList(ctx)
	if !found {
		return types.ErrAuthorizationListNotFound
	}
	policyRequired, err := authorizationsList.GetAuthorizedPolicy(msgURL)
	if err != nil {
		return errors.Wrap(types.ErrAuthorizationNotFound, fmt.Sprintf("msg: %v", msgURL))
	}
	//// TODO : check for empty policy
	//if policyRequired == types.PolicyType_groupOperational {
	//	return errors.Wrap(types.ErrMsgNotAuthorized, fmt.Sprintf("msg: %v", sdk.MsgTypeURL(msg)))
	//}
	policies, found := k.GetPolicies(ctx)
	if !found {
		return errors.Wrap(types.ErrPoliciesNotFound, fmt.Sprintf("msg: %v", msgURL))
	}

	return policies.CheckSigner(signer, policyRequired)
}
