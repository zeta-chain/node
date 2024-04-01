package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_GetParams(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	params := k.GetParams(ctx)
	// TODO: there is no get params method?
	k.SetParams(ctx, types.Params{Enabled: false})
	require.Equal(t, types.NewParams(), params)
	require.True(t, params.Enabled)
}

func TestKeeper_QueryParams(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.Params(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.Params(ctx, &types.QueryParamsRequest{})
		require.NoError(t, err)
		require.True(t, res.Params.Enabled)
	})
}
