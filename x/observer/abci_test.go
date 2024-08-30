package observer_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestBeginBlocker(t *testing.T) {
	t.Run("should not update LastObserverCount if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		observer.BeginBlocker(ctx, *k)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)

		_, found = k.GetKeygen(ctx)
		require.False(t, found)
	})

	t.Run("should not update LastObserverCount if observer set not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		count := 1
		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: uint64(count),
		})

		observer.BeginBlocker(ctx, *k)

		lastObserverCount, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(count), lastObserverCount.Count)
		require.Equal(t, int64(0), lastObserverCount.LastChangeHeight)

		_, found = k.GetKeygen(ctx)
		require.False(t, found)
	})

	t.Run("should not update LastObserverCount if observer set count equal last observed count", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		count := 1
		os := sample.ObserverSet(count)
		k.SetObserverSet(ctx, os)
		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: uint64(count),
		})

		observer.BeginBlocker(ctx, *k)

		lastObserverCount, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(count), lastObserverCount.Count)
		require.Equal(t, int64(0), lastObserverCount.LastChangeHeight)

		_, found = k.GetKeygen(ctx)
		require.False(t, found)
	})

	t.Run("should update LastObserverCount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		observeSetLen := 10
		count := 1
		os := sample.ObserverSet(observeSetLen)
		k.SetObserverSet(ctx, os)
		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: uint64(count),
		})

		keygen, found := k.GetKeygen(ctx)
		require.False(t, found)
		require.Equal(t, types.Keygen{}, keygen)

		observer.BeginBlocker(ctx, *k)

		keygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.Empty(t, keygen.GranteePubkeys)
		require.Equal(t, types.KeygenStatus_PendingKeygen, keygen.Status)
		require.Equal(t, int64(math.MaxInt64), keygen.BlockNumber)

		inboundEnabled := k.IsInboundEnabled(ctx)
		require.False(t, inboundEnabled)

		lastObserverCount, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(observeSetLen), lastObserverCount.Count)
		require.Equal(t, ctx.BlockHeight(), lastObserverCount.LastChangeHeight)
	})
}
