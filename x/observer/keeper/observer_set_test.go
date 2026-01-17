package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_GetObserverSet(t *testing.T) {
	t.Run("get observer set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		_, found := k.GetObserverSet(ctx)
		require.False(t, found)
		k.SetObserverSet(ctx, os)
		tfm, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, os, tfm)
	})
}

func TestKeeper_IsAddressPartOfObserverSet(t *testing.T) {
	t.Run("address is part of observer set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		require.False(t, k.IsAddressPartOfObserverSet(ctx, os.ObserverList[0]))
		k.SetObserverSet(ctx, os)
		require.True(t, k.IsAddressPartOfObserverSet(ctx, os.ObserverList[0]))
		require.False(t, k.IsAddressPartOfObserverSet(ctx, sample.AccAddress()))
	})
}

func TestKeeper_AddObserverToSet(t *testing.T) {
	t.Run("add observer to set", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObserverSet(ctx, os)
		newObserver := sample.AccAddress()

		// ACT
		countReturned, err := k.AddObserverToSet(ctx, newObserver)

		// ASSERT
		require.NoError(t, err)
		require.True(t, k.IsAddressPartOfObserverSet(ctx, newObserver))
		require.False(t, k.IsAddressPartOfObserverSet(ctx, sample.AccAddress()))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, len(os.ObserverList)+1)
		count, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, osNew.LenUint(), count.Count)
		require.Equal(t, osNew.LenUint(), countReturned)
	})

	t.Run("add observer to set if set doesn't exist", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		newObserver := sample.AccAddress()

		// ACT
		countReturned, err := k.AddObserverToSet(ctx, newObserver)

		// ASSERT
		require.NoError(t, err)
		require.True(t, k.IsAddressPartOfObserverSet(ctx, newObserver))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, 1)
		count, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, osNew.LenUint(), count.Count)
		require.Equal(t, osNew.LenUint(), countReturned)
	})

	t.Run("cannot add observer to set the address is already part of the set", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		newObserver := sample.AccAddress()
		_, err := k.AddObserverToSet(ctx, newObserver)
		require.NoError(t, err)
		require.True(t, k.IsAddressPartOfObserverSet(ctx, newObserver))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, 1)

		// ACT
		_, err = k.AddObserverToSet(ctx, newObserver)

		// ASSERT
		require.ErrorIs(t, err, types.ErrDuplicateObserver)
	})
}

func TestKeeper_RemoveObserverFromSet(t *testing.T) {
	t.Run("remove observer from set", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObserverSet(ctx, os)
		observerToRemove := os.ObserverList[0]

		// ACT
		count := k.RemoveObserverFromSet(ctx, observerToRemove)

		// ASSERT
		require.Equal(t, uint64(9), count)
		require.False(t, k.IsAddressPartOfObserverSet(ctx, observerToRemove))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, 9)
	})

	t.Run("returns 0 when observer set not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ACT
		count := k.RemoveObserverFromSet(ctx, sample.AccAddress())

		// ASSERT
		require.Equal(t, uint64(0), count)
	})

	t.Run("returns existing count when observer not in set", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(5)
		k.SetObserverSet(ctx, os)
		nonExistentObserver := sample.AccAddress()

		// ACT
		count := k.RemoveObserverFromSet(ctx, nonExistentObserver)

		// ASSERT
		require.Equal(t, uint64(5), count)
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, 5)
	})

	t.Run("remove last observer from set returns 0", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(1)
		k.SetObserverSet(ctx, os)
		observerToRemove := os.ObserverList[0]

		// ACT
		count := k.RemoveObserverFromSet(ctx, observerToRemove)

		// ASSERT
		require.Equal(t, uint64(0), count)
		require.False(t, k.IsAddressPartOfObserverSet(ctx, observerToRemove))
		osNew, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Len(t, osNew.ObserverList, 0)
	})
}

func TestKeeper_UpdateObserverAddress(t *testing.T) {
	t.Run("update observer address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(10)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.Error(t, err)
		k.SetObserverSet(ctx, observerSet)
		err = k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.NoError(t, err)
		observerSet, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, newObserverAddress, observerSet.ObserverList[len(observerSet.ObserverList)-1])
	})
	t.Run("unable to update observer list observe set not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()

		// ACT
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)

		// ASSERT
		require.ErrorIs(t, err, types.ErrObserverSetNotFound)
	})
	t.Run("unable to update observer list if the new list is not valid", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(10)
		observerSet.ObserverList = append(observerSet.ObserverList, []string{oldObserverAddress, newObserverAddress}...)
		k.SetObserverSet(ctx, observerSet)

		// ACT
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)

		// ASSERT
		require.ErrorContains(t, err, types.ErrDuplicateObserver.Error())
	})
	t.Run("should error if observer address not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(10)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		k.SetObserverSet(ctx, observerSet)
		err := k.UpdateObserverAddress(ctx, sample.AccAddress(), newObserverAddress)
		require.ErrorIs(t, err, types.ErrObserverNotFound)
	})
	t.Run("update observer address long observerList", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(10000)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		k.SetObserverSet(ctx, observerSet)
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.NoError(t, err)
		observerSet, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, newObserverAddress, observerSet.ObserverList[len(observerSet.ObserverList)-1])
	})
	t.Run("update observer address short observerList", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		observerSet := sample.ObserverSet(1)
		observerSet.ObserverList = append(observerSet.ObserverList, oldObserverAddress)
		k.SetObserverSet(ctx, observerSet)
		err := k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		require.NoError(t, err)
		observerSet, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, newObserverAddress, observerSet.ObserverList[len(observerSet.ObserverList)-1])
	})
}
