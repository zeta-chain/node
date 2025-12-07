package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestLastObserverCount(t *testing.T) {
	t.Run("should return false if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
	})

	t.Run("should set and get last observer count", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ARRANGE
		lbc := &types.LastObserverCount{
			Count:            10,
			LastChangeHeight: 100,
		}

		// ACT
		k.SetLastObserverCount(ctx, lbc)
		result, found := k.GetLastObserverCount(ctx)

		// ASSERT
		require.True(t, found)
		require.Equal(t, lbc.Count, result.Count)
		require.Equal(t, lbc.LastChangeHeight, result.LastChangeHeight)
	})

	t.Run("should overwrite existing value", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ARRANGE
		lbc1 := &types.LastObserverCount{
			Count:            10,
			LastChangeHeight: 100,
		}
		lbc2 := &types.LastObserverCount{
			Count:            20,
			LastChangeHeight: 200,
		}

		// ACT
		k.SetLastObserverCount(ctx, lbc1)
		k.SetLastObserverCount(ctx, lbc2)
		result, found := k.GetLastObserverCount(ctx)

		// ASSERT
		require.True(t, found)
		require.Equal(t, lbc2.Count, result.Count)
		require.Equal(t, lbc2.LastChangeHeight, result.LastChangeHeight)
	})
}

func TestDecrementLastObserverCount(t *testing.T) {
	t.Run("should do nothing if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ACT - should not panic
		k.DecrementLastObserverCount(ctx)

		// ASSERT
		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
	})

	t.Run("should decrement count", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ARRANGE
		lbc := &types.LastObserverCount{
			Count:            10,
			LastChangeHeight: 100,
		}
		k.SetLastObserverCount(ctx, lbc)

		// ACT
		k.DecrementLastObserverCount(ctx)
		result, found := k.GetLastObserverCount(ctx)

		// ASSERT
		require.True(t, found)
		require.Equal(t, uint64(9), result.Count)
	})

	t.Run("should not decrement below zero", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ARRANGE
		lbc := &types.LastObserverCount{
			Count:            0,
			LastChangeHeight: 100,
		}
		k.SetLastObserverCount(ctx, lbc)

		// ACT
		k.DecrementLastObserverCount(ctx)
		result, found := k.GetLastObserverCount(ctx)

		// ASSERT
		require.True(t, found)
		require.Equal(t, uint64(0), result.Count)
	})
}
