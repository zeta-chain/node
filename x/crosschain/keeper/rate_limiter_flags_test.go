package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetRateLimiterFlags(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)

	// not found
	_, found := k.GetRateLimiterFlags(ctx)
	require.False(t, found)

	flags := sample.RateLimiterFlags()

	k.SetRateLimiterFlags(ctx, flags)
	r, found := k.GetRateLimiterFlags(ctx)
	require.True(t, found)
	require.Equal(t, flags, r)
}
