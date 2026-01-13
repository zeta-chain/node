package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/zeta-chain/node/testutil/keeper"
)

func TestCounter(t *testing.T) {
	keeper, ctx, _, _ := testkeeper.CrosschainKeeper(t)
	initialCounter := keeper.GetCctxCounter(ctx)
	require.Zero(t, initialCounter)

	nextVal := keeper.GetNextCctxCounter(ctx)
	require.Greater(t, nextVal, initialCounter)
	require.Equal(t, nextVal, keeper.GetCctxCounter(ctx))

	// also test direct set
	nextVal += 1
	keeper.SetCctxCounter(ctx, nextVal)
	require.Equal(t, nextVal, keeper.GetCctxCounter(ctx))
}
