package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// SetCrosschainFlags set the crosschain flags in the store
func (k Keeper) SetCrosschainFlags(ctx sdk.Context, permissionFlags types.CrosschainFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CrosschainFlagsKey))
	b := k.cdc.MustMarshal(&permissionFlags)
	store.Set([]byte{0}, b)
}

// GetCrosschainFlags returns the crosschain flags
func (k Keeper) GetCrosschainFlags(ctx sdk.Context) (val types.CrosschainFlags, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CrosschainFlagsKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) IsInboundEnabled(ctx sdk.Context) (found bool) {
	flags, found := k.GetCrosschainFlags(ctx)
	if !found {
		return false
	}
	return flags.IsInboundEnabled
}

func (k Keeper) IsOutboundEnabled(ctx sdk.Context) (found bool) {
	flags, found := k.GetCrosschainFlags(ctx)
	if !found {
		return false
	}
	return flags.IsOutboundEnabled
}

// RemovePermissionFlags removes permissionFlags from the store
func (k Keeper) RemovePermissionFlags(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CrosschainFlagsKey))
	store.Delete([]byte{0})
}

func (k Keeper) DisableInboundOnly(ctx sdk.Context) {
	flags, found := k.GetCrosschainFlags(ctx)
	if !found {
		flags.IsOutboundEnabled = true
	}
	flags.IsInboundEnabled = false
	k.SetCrosschainFlags(ctx, flags)
}
