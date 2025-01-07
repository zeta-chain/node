package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

func (k Keeper) SetOperationalFlags(ctx sdk.Context, operationalFlags types.OperationalFlags) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&operationalFlags)
	key := types.KeyPrefix(types.OperationalFlagsKey)
	store.Set(key, b)
}

func (k Keeper) GetOperationalFlags(ctx sdk.Context) (val types.OperationalFlags, found bool) {
	found = false
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyPrefix(types.OperationalFlagsKey))
	if b == nil {
		return
	}
	found = true
	k.cdc.MustUnmarshal(b, &val)
	return
}
