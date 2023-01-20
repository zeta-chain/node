package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) SetEmissionTracker(ctx sdk.Context, tracker *types.EmissionTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EmissionsTrackerKey))
	key := fmt.Sprintf("%s", tracker.Type.String())
	b := k.cdc.MustMarshal(tracker)
	store.Set([]byte(key), b)
}

func (k Keeper) GetEmissionTracker(ctx sdk.Context, category types.EmissionCategory) (val *types.EmissionTracker, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EmissionsTrackerKey))
	b := store.Get(types.KeyPrefix(category.String()))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, val)
	return val, true
}
