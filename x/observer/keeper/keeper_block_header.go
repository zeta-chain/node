package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// SetKeygen set keygen in the store
func (k Keeper) SetBlockHeaderState(ctx sdk.Context, blockHeaderState types.BlockHeaderState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderStateKey))
	b := k.cdc.MustMarshal(&blockHeaderState)
	key := strconv.FormatInt(blockHeaderState.ChainId, 10)
	store.Set(types.KeyPrefix(key), b)
}

func (k Keeper) GetBlockHeaderState(ctx sdk.Context, chainID int64) (val types.BlockHeaderState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderStateKey))

	b := store.Get(types.KeyPrefix(strconv.FormatInt(chainID, 10)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
