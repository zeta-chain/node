package v5

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func MigrateStore(ctx sdk.Context, observerStoreKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	var legacyObserverMappers []*types.ObserverMapper
	legacyObserverMapperStore := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.ObserverMapperKey))
	iterator := sdk.KVStorePrefixIterator(legacyObserverMapperStore, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ObserverMapper
		cdc.MustUnmarshal(iterator.Value(), &val)
		legacyObserverMappers = append(legacyObserverMappers, &val)
	}

	// We can safely assume that the observer list is the same for all the observer mappers
	observerList := legacyObserverMappers[0].ObserverList

	storelastBlockObserverCount := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.LastBlockObserverCountKey))
	b := cdc.MustMarshal(&types.LastObserverCount{Count: uint64(len(observerList)), LastChangeHeight: ctx.BlockHeight()})
	storelastBlockObserverCount.Set([]byte{0}, b)

	storeObserverSet := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.ObserverSetKey))
	b = cdc.MustMarshal(&types.ObserverSet{ObserverList: observerList})
	storeObserverSet.Set([]byte{0}, b)

	for _, legacyObserverMapper := range legacyObserverMappers {
		legacyObserverMapperStore.Delete(types.KeyPrefix(legacyObserverMapper.Index))
	}
	return nil
}
