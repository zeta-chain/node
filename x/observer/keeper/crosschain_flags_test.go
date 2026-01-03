package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_IsInboundEnabled(t *testing.T) {
	t.Run("should return false if flags not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		enabled := k.IsInboundEnabled(ctx)
		require.False(t, enabled)
	})

	t.Run("should return if flags found and set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: false,
		})
		enabled := k.IsInboundEnabled(ctx)
		require.False(t, enabled)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})

		enabled = k.IsInboundEnabled(ctx)
		require.True(t, enabled)
	})
}

func TestKeeper_IsOutboundEnabled(t *testing.T) {
	t.Run("should return false if flags not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		enabled := k.IsOutboundEnabled(ctx)
		require.False(t, enabled)
	})

	t.Run("should return if flags found and set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsOutboundEnabled: false,
		})
		enabled := k.IsOutboundEnabled(ctx)
		require.False(t, enabled)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsOutboundEnabled: true,
		})

		enabled = k.IsOutboundEnabled(ctx)
		require.True(t, enabled)
	})
}

func TestKeeper_DisableInboundOnly(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	k.DisableInboundOnly(ctx)
	enabled := k.IsOutboundEnabled(ctx)
	require.True(t, enabled)

	enabled = k.IsInboundEnabled(ctx)
	require.False(t, enabled)

	k.SetCrosschainFlags(ctx, types.CrosschainFlags{
		IsOutboundEnabled: false,
	})

	k.DisableInboundOnly(ctx)
	enabled = k.IsOutboundEnabled(ctx)
	require.False(t, enabled)

	enabled = k.IsInboundEnabled(ctx)
	require.False(t, enabled)
}

func TestKeeper_IsV2ZetaEnabled(t *testing.T) {
	t.Run("should return false if flags not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		enabled := k.IsV2ZetaEnabled(ctx)
		require.False(t, enabled)
	})

	t.Run("should return if flags found and set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsV2ZetaEnabled: false,
		})
		enabled := k.IsV2ZetaEnabled(ctx)
		require.False(t, enabled)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsV2ZetaEnabled: true,
		})

		enabled = k.IsV2ZetaEnabled(ctx)
		require.True(t, enabled)
	})
}
