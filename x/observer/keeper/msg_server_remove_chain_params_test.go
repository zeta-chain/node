package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_RemoveChainParams(t *testing.T) {
	t.Run("can update chain params", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chain1 := common.ExternalChainList()[0].ChainId
		chain2 := common.ExternalChainList()[1].ChainId
		chain3 := common.ExternalChainList()[2].ChainId

		// set admin
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// add chain params
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(chain1),
				sample.ChainParams(chain2),
				sample.ChainParams(chain3),
			},
		})

		// remove chain params
		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain2,
		})
		assert.NoError(t, err)

		// check list has two chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		assert.True(t, found)
		assert.Len(t, chainParamsList.ChainParams, 2)
		assert.Equal(t, chain1, chainParamsList.ChainParams[0].ChainId)
		assert.Equal(t, chain3, chainParamsList.ChainParams[1].ChainId)

		// remove chain params
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain1,
		})
		assert.NoError(t, err)

		// check list has one chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		assert.True(t, found)
		assert.Len(t, chainParamsList.ChainParams, 1)
		assert.Equal(t, chain3, chainParamsList.ChainParams[0].ChainId)

		// remove chain params
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain3,
		})
		assert.NoError(t, err)

		// check list has no chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		assert.True(t, found)
		assert.Len(t, chainParamsList.ChainParams, 0)
	})

	t.Run("cannot remove chain params if not authorized", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     sample.AccAddress(),
			ChainParams: sample.ChainParams(common.ExternalChainList()[0].ChainId),
		})
		assert.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

		// group 1 should not be able to update core params
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainParams{
			Creator:     sample.AccAddress(),
			ChainParams: sample.ChainParams(common.ExternalChainList()[0].ChainId),
		})
		assert.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

	})

	t.Run("cannot remove if chain ID not found", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// set admin
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// not found if no chain params
		_, found := k.GetChainParamsList(ctx)
		assert.False(t, found)

		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: common.ExternalChainList()[0].ChainId,
		})
		assert.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// add chain params
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(common.ExternalChainList()[0].ChainId),
				sample.ChainParams(common.ExternalChainList()[1].ChainId),
				sample.ChainParams(common.ExternalChainList()[2].ChainId),
			},
		})

		// not found if chain ID not in list
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: common.ExternalChainList()[3].ChainId,
		})
		assert.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})
}
