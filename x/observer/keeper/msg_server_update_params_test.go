package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_UpdateParams(t *testing.T) {
	t.Run("successfully update params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
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
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    types.DefaultParams(),
		})

		require.Error(t, err)
		require.Nil(t, res)
	})
}
