package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_RemoveChainParams(t *testing.T) {
	t.Run("can update chain params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		chainList := chains.ExternalChainList([]chains.Chain{})

		chain1 := chainList[0].ChainId
		chain2 := chainList[1].ChainId
		chain3 := chainList[2].ChainId

		// set admin
		admin := sample.AccAddress()
		// add chain params
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(chain1),
				sample.ChainParams(chain2),
				sample.ChainParams(chain3),
			},
		})

		// remove chain params
		msg := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain2,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has two chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 2)
		require.Equal(t, chain1, chainParamsList.ChainParams[0].ChainId)
		require.Equal(t, chain3, chainParamsList.ChainParams[1].ChainId)

		// remove chain params
		msg = types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain1,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has one chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chain3, chainParamsList.ChainParams[0].ChainId)

		// remove chain params
		msg = types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain3,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has no chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 0)
	})

	t.Run("cannot remove chain params if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chains.ExternalChainList([]chains.Chain{})[0].ChainId,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("cannot remove if chain ID not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// set admin
		admin := sample.AccAddress()

		// not found if no chain params
		_, found := k.GetChainParamsList(ctx)
		require.False(t, found)

		chainList := chains.ExternalChainList([]chains.Chain{})

		msg := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chainList[0].ChainId,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// add chain params
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(chainList[0].ChainId),
				sample.ChainParams(chainList[1].ChainId),
				sample.ChainParams(chainList[2].ChainId),
			},
		})

		// not found if chain ID not in list
		msg = types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chainList[3].ChainId,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})
}
