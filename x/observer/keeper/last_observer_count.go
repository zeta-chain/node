package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

func (k Keeper) SetLastObserverCount(ctx sdk.Context, lbc *types.LastObserverCount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockObserverCountKey))
	b := k.cdc.MustMarshal(lbc)
	store.Set([]byte{0}, b)
}

func (k Keeper) GetLastObserverCount(ctx sdk.Context) (val types.LastObserverCount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockObserverCountKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
