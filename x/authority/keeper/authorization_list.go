package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

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

// CheckAuthorization uses both the authorization list and the policies to check if the signer is authorized
func (k Keeper) CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error {
	// Policy transactions must have only one signer
	if len(msg.GetSigners()) != 1 {
		return errors.Wrapf(types.ErrSigners, "msg: %v", sdk.MsgTypeURL(msg))
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
	if policyRequired == types.PolicyType_groupEmpty {
		return errors.Wrap(types.ErrInvalidPolicyType, fmt.Sprintf("Empty policy for msg: %v", msgURL))
	}

	policies, found := k.GetPolicies(ctx)
	if !found {
		return errors.Wrap(types.ErrPoliciesNotFound, fmt.Sprintf("msg: %v", msgURL))
	}

	return policies.CheckSigner(signer, policyRequired)
}
