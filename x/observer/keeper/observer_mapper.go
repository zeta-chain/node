package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func GetObserverMapperIndex(chain *common.Chain) string {
	return fmt.Sprintf("%d", chain.ChainId)
}

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

func (k Keeper) SetObserverMapper(ctx sdk.Context, om *types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	om.Index = GetObserverMapperIndex(om.ObserverChain)
	b := k.cdc.MustMarshal(om)
	store.Set([]byte(om.Index), b)
}

func (k Keeper) GetObserverMapper(ctx sdk.Context, chain *common.Chain) (val types.ObserverMapper, found bool) {
	index := GetObserverMapperIndex(chain)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllObserverMappers(ctx sdk.Context) (mappers []*types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ObserverMapper
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		mappers = append(mappers, &val)
	}
	return
}
func (k Keeper) GetAllObserverMappersForAddress(ctx sdk.Context, address string) (mappers []*types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ObserverMapper
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		addToList := false
		for _, addr := range val.ObserverList {
			if addr == address {
				addToList = true
			}
		}
		if addToList {
			mappers = append(mappers, &val)
		}
	}
	return
}

// Utils

func (k Keeper) GetAllObserverAddresses(ctx sdk.Context) []string {
	var val []string
	mappers := k.GetAllObserverMappers(ctx)
	for _, mapper := range mappers {
		val = append(val, mapper.ObserverList...)
	}
	allKeys := make(map[string]bool)
	var dedupedList []string
	for _, item := range val {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			dedupedList = append(dedupedList, item)
		}
	}
	return dedupedList
}

func (k Keeper) AddObserverToMapper(ctx sdk.Context, chain *common.Chain, address string) {
	mapper, found := k.GetObserverMapper(ctx, chain)
	if !found {
		k.SetObserverMapper(ctx, &types.ObserverMapper{
			Index:         "",
			ObserverChain: chain,
			ObserverList:  []string{address},
		})
		return
	}
	// Return if duplicate
	for _, addr := range mapper.ObserverList {
		if addr == address {
			return
		}
	}
	mapper.ObserverList = append(mapper.ObserverList, address)
	k.SetObserverMapper(ctx, &mapper)
}
