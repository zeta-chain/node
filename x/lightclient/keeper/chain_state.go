package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

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
