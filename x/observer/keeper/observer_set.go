package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
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

// AddObserverToSet adds an observer to the observer set.It makes sure the updated observer set is valid.
// It also sets the observer count and returns the updated length of the observer set.
func (k Keeper) AddObserverToSet(ctx sdk.Context, address string) (uint64, error) {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		observerSet = types.ObserverSet{
			ObserverList: []string{},
		}
	}

	observerSet.ObserverList = append(observerSet.ObserverList, address)
	if err := observerSet.Validate(); err != nil {
		return 0, err
	}

	k.SetObserverSet(ctx, observerSet)
	newCount := observerSet.LenUint()
	k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: newCount})

	return newCount, nil
}

// RemoveObserverFromSet removes an observer from the observer set.
// Returns the observer count after the operation
func (k Keeper) RemoveObserverFromSet(ctx sdk.Context, address string) uint64 {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return 0
	}
	for i, addr := range observerSet.ObserverList {
		if addr == address {
			observerSet.ObserverList = append(observerSet.ObserverList[:i], observerSet.ObserverList[i+1:]...)
			k.SetObserverSet(ctx, observerSet)
			break
		}
	}
	return observerSet.LenUint()
}

// UpdateObserverAddress updates an observer address in the observer set.It makes sure the updated observer set is valid.
func (k Keeper) UpdateObserverAddress(ctx sdk.Context, oldObserverAddress, newObserverAddress string) error {
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return types.ErrObserverSetNotFound
	}
	found = false
	for i, addr := range observerSet.ObserverList {
		if addr == oldObserverAddress {
			observerSet.ObserverList[i] = newObserverAddress
			found = true
			break
		}
	}
	if !found {
		return errors.Wrapf(types.ErrObserverNotFound, "observer %s", oldObserverAddress)
	}

	err := observerSet.Validate()
	if err != nil {
		return errors.Wrap(types.ErrUpdateObserver, err.Error())
	}
	k.SetObserverSet(ctx, observerSet)
	return nil
}
