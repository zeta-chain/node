package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_OperationalFlags(t *testing.T) {
	t.Run("should return operational flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		// should return zero value by default
		res, err := k.OperationalFlags(wctx, &types.QueryOperationalFlagsRequest{})
		require.NoError(t, err)
		require.Equal(t, types.OperationalFlags{}, res.OperationalFlags)

		// set the value and ensure it's returned by the query
		restartHeight := int64(100)
		k.SetOperationalFlags(ctx, types.OperationalFlags{
			RestartHeight: restartHeight,
		})
		res, err = k.OperationalFlags(wctx, &types.QueryOperationalFlagsRequest{})
		require.NoError(t, err)
		require.Equal(t, restartHeight, res.OperationalFlags.RestartHeight)

	})
}
