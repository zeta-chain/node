package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetObserverSet(t *testing.T) {
	t.Run("get observer set", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObserverSet(ctx, os)
		tfm, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, os, tfm)
	})
}

func TestKeeper_IsAddressPartOfObserverSet(t *testing.T) {
	t.Run("address is part of observer set", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObserverSet(ctx, os)
		require.True(t, k.IsAddressPartOfObserverSet(ctx, os.ObserverList[0]))
		require.False(t, k.IsAddressPartOfObserverSet(ctx, sample.AccAddress()))
	})
}

func TestKeeper_AddObserverToSet(t *testing.T) {
	t.Run("add observer to set", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObserverSet(ctx, os)
		newObserver := sample.AccAddress()
		k.AddObserverToSet(ctx, newObserver)
		require.True(t, k.IsAddressPartOfObserverSet(ctx, newObserver))
		require.False(t, k.IsAddressPartOfObserverSet(ctx, sample.AccAddress()))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, len(os.ObserverList)+1)
	})
}

func TestKeeper_RemoveObserverFromSet(t *testing.T) {
	t.Run("remove observer from set", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObserverSet(ctx, os)
		k.RemoveObserverFromSet(ctx, os.ObserverList[0])
		require.False(t, k.IsAddressPartOfObserverSet(ctx, os.ObserverList[0]))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, len(os.ObserverList)-1)
	})
}

func TestKeeper_UpdateObserverAddress(t *testing.T) {
	t.Run("update observer address", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(10)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		k.SetObserverSet(ctx, observerSet)
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.NoError(t, err)
		observerSet, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, newObserverAddress, observerSet.ObserverList[len(observerSet.ObserverList)-1])
	})
	t.Run("update observer address long observerList", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(10000)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		k.SetObserverSet(ctx, observerSet)
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.NoError(t, err)
		observerMappers, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, newObserverAddress, observerMappers.ObserverList[len(observerMappers.ObserverList)-1])
	})
	t.Run("update observer address short observerList", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(1)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		k.SetObserverSet(ctx, observerSet)
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.NoError(t, err)
		observerMappers, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, newObserverAddress, observerMappers.ObserverList[len(observerMappers.ObserverList)-1])
	})
}
