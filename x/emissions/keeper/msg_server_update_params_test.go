package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
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
		require.Equal(t, types.DefaultParams(), k.GetParams(ctx))
	})

	t.Run("fail for wrong authority", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    types.DefaultParams(),
		})

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("fail for invalid params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		params := types.DefaultParams()
		params.ValidatorEmissionPercentage = "-1.5"
		res, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: k.GetAuthority(),
			Params:    params,
		})

		require.Error(t, err)
		require.Nil(t, res)
	})
}
