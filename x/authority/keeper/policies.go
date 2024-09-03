package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// SetPolicies sets the policies to the store
func (k Keeper) SetPolicies(ctx sdk.Context, policies types.Policies) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoliciesKey))
	b := k.cdc.MustMarshal(&policies)
	store.Set([]byte{0}, b)
}

// GetPolicies returns the policies from the store
func (k Keeper) GetPolicies(ctx sdk.Context) (val types.Policies, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoliciesKey))
	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
