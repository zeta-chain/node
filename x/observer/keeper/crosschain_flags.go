package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

// SetCrosschainFlags set the crosschain flags in the store
func (k Keeper) SetCrosschainFlags(ctx sdk.Context, crosschainFlags types.CrosschainFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CrosschainFlagsKey))
	b := k.cdc.MustMarshal(&crosschainFlags)
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

// RemoveCrosschainFlags removes crosschain flags from the store
func (k Keeper) RemoveCrosschainFlags(ctx sdk.Context) {
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

// IsV2ZetaEnabled returns true if V2 ZETA gateway flows are enabled
func (k Keeper) IsV2ZetaEnabled(ctx sdk.Context) bool {
	flags, found := k.GetCrosschainFlags(ctx)
	if !found {
		return false
	}
	return flags.IsV2ZetaEnabled
}
