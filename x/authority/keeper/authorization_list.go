package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

// SetAuthorizationList sets the authorization list to the store
func (k Keeper) SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList) error {
	err := list.Validate()
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&list)
	store.Set([]byte{0}, b)
	return nil
}

// GetAuthorizationList returns the authorization list from the store
func (k Keeper) GetAuthorizationList(ctx sdk.Context) (val types.AuthorizationList, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) UpdateAuthorizationList(ctx sdk.Context, addList types.AuthorizationList, removeList types.AuthorizationList) types.AuthorizationList {
	list, found := k.GetAuthorizationList(ctx)
	if !found {
		return addList
	}
	list.AddAuthorizations(addList)
	list.RemoveAuthorizations(removeList)
	return list
}
