package keeper

import (
	"strconv"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/lightclient/types"
)

// GetAllChainStates returns all chain states
func (k Keeper) GetAllChainStates(ctx sdk.Context) (list []types.ChainState) {
	p := types.KeyPrefix(types.ChainStateKey)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ChainState
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// SetChainState set a specific chain state in the store from its index
func (k Keeper) SetChainState(ctx sdk.Context, chainState types.ChainState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainStateKey))
	b := k.cdc.MustMarshal(&chainState)
	key := strconv.FormatInt(chainState.ChainId, 10)
	store.Set(types.KeyPrefix(key), b)
}

// GetChainState returns a chain state from its chainID
func (k Keeper) GetChainState(ctx sdk.Context, chainID int64) (val types.ChainState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainStateKey))

	b := store.Get(types.KeyPrefix(strconv.FormatInt(chainID, 10)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
