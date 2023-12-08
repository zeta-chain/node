package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) SetBlame(ctx sdk.Context, blame *types.Blame) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	b := k.cdc.MustMarshal(blame)
	store.Set([]byte(blame.Index), b)
}

func (k Keeper) GetBlame(ctx sdk.Context, index string) (val types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllBlame(ctx sdk.Context) (BlameRecords []*types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	found = false
	for ; iterator.Valid(); iterator.Next() {
		var val types.Blame
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		BlameRecords = append(BlameRecords, &val)
		found = true
	}
	return
}

func (k Keeper) GetBlamesByChainAndNonce(ctx sdk.Context, chainID int64, nonce int64) (BlameRecords []*types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	blamePrefix := types.GetBlamePrefix(chainID, nonce)
	iterator := sdk.KVStorePrefixIterator(store, []byte(blamePrefix))
	defer iterator.Close()
	found = false
	for ; iterator.Valid(); iterator.Next() {
		var val types.Blame
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		BlameRecords = append(BlameRecords, &val)
		found = true
	}
	return
}
