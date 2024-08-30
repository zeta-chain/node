package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/emissions/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestMsgServer_UpdateParams(t *testing.T) {
	t.Run("successfully update params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: k.GetAuthority(),
			Params:    types.DefaultParams(),
		})

		require.NoError(t, err)
		require.Empty(t, res)
		params, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultParams(), params)
	})

	t.Run("fail for wrong authority", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    types.DefaultParams(),
		})

		require.Error(t, err)
	})

	t.Run("fail for invalid params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		params := types.DefaultParams()
		params.ValidatorEmissionPercentage = "-1.5"
		_, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: k.GetAuthority(),
			Params:    params,
		})

		require.ErrorIs(t, err, types.ErrUnableToSetParams)
	})
}
