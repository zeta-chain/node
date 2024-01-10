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

func TestMsgServer_UpdateChainParams(t *testing.T) {
	t.Run("can update chain params", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chain1 := common.ExternalChainList()[0].ChainId
		chain2 := common.ExternalChainList()[1].ChainId
		chain3 := common.ExternalChainList()[2].ChainId

		// set admin
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// check list initially empty
		_, found := k.GetChainParamsList(ctx)
		require.False(t, found)

		// a new chain params can be added
		chainParams1 := sample.ChainParams(chain1)
		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams1,
		})
		require.NoError(t, err)

		// check list has one chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])

		// a new chian params can be added
		chainParams2 := sample.ChainParams(chain2)
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams2,
		})
		require.NoError(t, err)

		// check list has two chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 2)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])

		// a new chain params can be added
		chainParams3 := sample.ChainParams(chain3)
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams3,
		})
		require.NoError(t, err)

		// check list has three chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 3)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])
		require.Equal(t, chainParams3, chainParamsList.ChainParams[2])

		// chain params can be updated
		chainParams2.ConfirmationCount = chainParams2.ConfirmationCount + 1
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams2,
		})
		require.NoError(t, err)

		// check list has three chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 3)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])
		require.Equal(t, chainParams3, chainParamsList.ChainParams[2])
	})

	t.Run("cannot update chain params if not authorized", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     sample.AccAddress(),
			ChainParams: sample.ChainParams(common.ExternalChainList()[0].ChainId),
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

		// group 1 should not be able to update chain params
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     sample.AccAddress(),
			ChainParams: sample.ChainParams(common.ExternalChainList()[0].ChainId),
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

	})
}
