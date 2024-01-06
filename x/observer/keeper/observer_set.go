package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) SetObserverSet(ctx sdk.Context, om types.ObserverSet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverSetKey))
	b := k.cdc.MustMarshal(&om)
	store.Set([]byte{0}, b)
}

func (k Keeper) GetObserverSet(ctx sdk.Context) (val types.ObserverSet, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverSetKey))
	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) IsAddressPartOfObserverSet(ctx sdk.Context, address string) bool {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return false
	}
	for _, addr := range observerSet.ObserverList {
		if addr == address {
			return true
		}
	}
	return false

}

func (k Keeper) AddObserverToSet(ctx sdk.Context, address string) {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{address},
		})
		return
	}
	for _, addr := range observerSet.ObserverList {
		if addr == address {
			return
		}
	}
	observerSet.ObserverList = append(observerSet.ObserverList, address)
	k.SetObserverSet(ctx, observerSet)
}

func (k Keeper) RemoveObserverFromSet(ctx sdk.Context, address string) {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return
	}
	for i, addr := range observerSet.ObserverList {
		if addr == address {
			observerSet.ObserverList = append(observerSet.ObserverList[:i], observerSet.ObserverList[i+1:]...)
			k.SetObserverSet(ctx, observerSet)
			return
		}
	}
}

func (k Keeper) UpdateObserverAddress(ctx sdk.Context, oldObserverAddress, newObserverAddress string) error {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return types.ErrObserverSetNotFound
	}
	for i, addr := range observerSet.ObserverList {
		if addr == oldObserverAddress {
			observerSet.ObserverList[i] = newObserverAddress
			k.SetObserverSet(ctx, observerSet)
			return nil
		}
	}
	return types.ErrUpdateObserver
}
