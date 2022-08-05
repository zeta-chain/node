package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

// SetERC20TokenPairs set eRC20TokenPairs in the store
func (k Keeper) SetERC20TokenPairs(ctx sdk.Context, eRC20TokenPairs types.ERC20TokenPairs) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ERC20TokenPairsKey))
	b := k.cdc.MustMarshal(&eRC20TokenPairs)
	store.Set([]byte{0}, b)
}

// GetERC20TokenPairs returns eRC20TokenPairs
func (k Keeper) GetERC20TokenPairs(ctx sdk.Context) (val types.ERC20TokenPairs, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ERC20TokenPairsKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveERC20TokenPairs removes eRC20TokenPairs from the store
func (k Keeper) RemoveERC20TokenPairs(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ERC20TokenPairsKey))
	store.Delete([]byte{0})
}
