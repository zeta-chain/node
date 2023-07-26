package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// SetPermissionFlags set permissionFlags in the store
func (k Keeper) SetPermissionFlags(ctx sdk.Context, permissionFlags types.PermissionFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PermissionFlagsKey))
	b := k.cdc.MustMarshal(&permissionFlags)
	store.Set([]byte{0}, b)
}

// GetPermissionFlags returns permissionFlags
func (k Keeper) GetPermissionFlags(ctx sdk.Context) (val types.PermissionFlags, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PermissionFlagsKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) IsInboundAllowed(ctx sdk.Context) (found bool) {
	flags, found := k.GetPermissionFlags(ctx)
	if !found {
		return false
	}
	return flags.IsInboundEnabled
}

// RemovePermissionFlags removes permissionFlags from the store
func (k Keeper) RemovePermissionFlags(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PermissionFlagsKey))
	store.Delete([]byte{0})
}
