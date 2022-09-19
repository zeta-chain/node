package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
)

// SetZetaDepositAndCallContract set zetaDepositAndCallContract in the store
func (k Keeper) SetZetaDepositAndCallContract(ctx sdk.Context, zetaDepositAndCallContract types.ZetaDepositAndCallContract) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaDepositAndCallContractKey))
	b := k.cdc.MustMarshal(&zetaDepositAndCallContract)
	store.Set([]byte{0}, b)
}

// GetZetaDepositAndCallContract returns zetaDepositAndCallContract
func (k Keeper) GetZetaDepositAndCallContract(ctx sdk.Context) (val types.ZetaDepositAndCallContract, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaDepositAndCallContractKey))

	b := store.Get([]byte{0})
    if b == nil {
        return val, false
    }

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveZetaDepositAndCallContract removes zetaDepositAndCallContract from the store
func (k Keeper) RemoveZetaDepositAndCallContract(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaDepositAndCallContractKey))
	store.Delete([]byte{0})
}
