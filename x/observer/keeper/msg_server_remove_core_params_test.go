package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_RemoveCoreParams(t *testing.T) {
	t.Run("can update core params", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chain1 := common.ExternalChainList()[0].ChainId
		chain2 := common.ExternalChainList()[1].ChainId
		chain3 := common.ExternalChainList()[2].ChainId

		// set admin
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// add core params
		k.SetCoreParamsList(ctx, types.CoreParamsList{
			CoreParams: []*types.CoreParams{
				sample.CoreParams(chain1),
				sample.CoreParams(chain2),
				sample.CoreParams(chain3),
			},
		})

		// remove core params
		_, err := srv.RemoveCoreParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveCoreParams{
			Creator: admin,
			ChainId: chain2,
		})
		require.NoError(t, err)

		// check list has two core params
		coreParamsList, found := k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 2)
		require.Equal(t, chain1, coreParamsList.CoreParams[0].ChainId)
		require.Equal(t, chain3, coreParamsList.CoreParams[1].ChainId)

		// remove core params
		_, err = srv.RemoveCoreParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveCoreParams{
			Creator: admin,
			ChainId: chain1,
		})
		require.NoError(t, err)

		// check list has one core params
		coreParamsList, found = k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 1)
		require.Equal(t, chain3, coreParamsList.CoreParams[0].ChainId)

		// remove core params
		_, err = srv.RemoveCoreParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveCoreParams{
			Creator: admin,
			ChainId: chain3,
		})
		require.NoError(t, err)

		// check list has no core params
		coreParamsList, found = k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 0)
	})

	t.Run("cannot remove core params if not authorized", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		_, err := srv.UpdateCoreParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateCoreParams{
			Creator:    sample.AccAddress(),
			CoreParams: sample.CoreParams(common.ExternalChainList()[0].ChainId),
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

		// group 1 should not be able to update core params
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		_, err = srv.UpdateCoreParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateCoreParams{
			Creator:    sample.AccAddress(),
			CoreParams: sample.CoreParams(common.ExternalChainList()[0].ChainId),
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

	})

	t.Run("cannot remove if chain ID not found", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// set admin
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// not found if no core params
		_, found := k.GetCoreParamsList(ctx)
		require.False(t, found)

		_, err := srv.RemoveCoreParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveCoreParams{
			Creator: admin,
			ChainId: common.ExternalChainList()[0].ChainId,
		})
		require.ErrorIs(t, err, types.ErrCoreParamsNotFound)

		// add core params
		k.SetCoreParamsList(ctx, types.CoreParamsList{
			CoreParams: []*types.CoreParams{
				sample.CoreParams(common.ExternalChainList()[0].ChainId),
				sample.CoreParams(common.ExternalChainList()[1].ChainId),
				sample.CoreParams(common.ExternalChainList()[2].ChainId),
			},
		})

		// not found if chain ID not in list
		_, err = srv.RemoveCoreParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveCoreParams{
			Creator: admin,
			ChainId: common.ExternalChainList()[3].ChainId,
		})
		require.ErrorIs(t, err, types.ErrCoreParamsNotFound)
	})
}
