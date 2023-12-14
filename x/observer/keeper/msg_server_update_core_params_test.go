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

func TestMsgServer_UpdateCoreParams(t *testing.T) {
	t.Run("can update core params", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chain1 := common.ExternalChainList()[0].ChainId
		chain2 := common.ExternalChainList()[1].ChainId
		chain3 := common.ExternalChainList()[2].ChainId

		// set admin
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// check list initially empty
		_, found := k.GetCoreParamsList(ctx)
		require.False(t, found)

		// a new core params can be added
		coreParams1 := sample.CoreParams(chain1)
		_, err := srv.UpdateCoreParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateCoreParams{
			Creator:    admin,
			CoreParams: coreParams1,
		})
		require.NoError(t, err)

		// check list has one core params
		coreParamsList, found := k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 1)
		require.Equal(t, coreParams1, coreParamsList.CoreParams[0])

		// a new core params can be added
		coreParams2 := sample.CoreParams(chain2)
		_, err = srv.UpdateCoreParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateCoreParams{
			Creator:    admin,
			CoreParams: coreParams2,
		})
		require.NoError(t, err)

		// check list has two core params
		coreParamsList, found = k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 2)
		require.Equal(t, coreParams1, coreParamsList.CoreParams[0])
		require.Equal(t, coreParams2, coreParamsList.CoreParams[1])

		// a new core params can be added
		coreParams3 := sample.CoreParams(chain3)
		_, err = srv.UpdateCoreParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateCoreParams{
			Creator:    admin,
			CoreParams: coreParams3,
		})
		require.NoError(t, err)

		// check list has three core params
		coreParamsList, found = k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 3)
		require.Equal(t, coreParams1, coreParamsList.CoreParams[0])
		require.Equal(t, coreParams2, coreParamsList.CoreParams[1])
		require.Equal(t, coreParams3, coreParamsList.CoreParams[2])

		// core params can be updated
		coreParams2.ConfirmationCount = coreParams2.ConfirmationCount + 1
		_, err = srv.UpdateCoreParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateCoreParams{
			Creator:    admin,
			CoreParams: coreParams2,
		})
		require.NoError(t, err)

		// check list has three core params
		coreParamsList, found = k.GetCoreParamsList(ctx)
		require.True(t, found)
		require.Len(t, coreParamsList.CoreParams, 3)
		require.Equal(t, coreParams1, coreParamsList.CoreParams[0])
		require.Equal(t, coreParams2, coreParamsList.CoreParams[1])
		require.Equal(t, coreParams3, coreParamsList.CoreParams[2])
	})

	t.Run("cannot update core params if not authorized", func(t *testing.T) {
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
}
