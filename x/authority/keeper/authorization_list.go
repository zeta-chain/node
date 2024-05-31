package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/authority/types"
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
