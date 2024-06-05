package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_RemoveChainParams(t *testing.T) {
	t.Run("can update chain params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		chain1 := chains.ExternalChainList()[0].ChainId
		chain2 := chains.ExternalChainList()[1].ChainId
		chain3 := chains.ExternalChainList()[2].ChainId

		// set admin
		admin := sample.AccAddress()
		msg := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain2,
		}

		// add chain params
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(chain1),
				sample.ChainParams(chain2),
				sample.ChainParams(chain3),
			},
		})

		// remove chain params
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has two chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 2)
		require.Equal(t, chain1, chainParamsList.ChainParams[0].ChainId)
		require.Equal(t, chain3, chainParamsList.ChainParams[1].ChainId)

		msg2 := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain1,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg2, nil)
		// remove chain params
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg2)
		require.NoError(t, err)

		// check list has one chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chain3, chainParamsList.ChainParams[0].ChainId)
		msg3 := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chain3,
		}

		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg3, nil)

		// remove chain params
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg3)
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
		msg := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chains.ExternalChainList()[0].ChainId,
		}

		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("cannot remove if chain ID not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		// set admin
		admin := sample.AccAddress()
		msg := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chains.ExternalChainList()[0].ChainId,
		}

		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// not found if no chain params
		_, found := k.GetChainParamsList(ctx)
		require.False(t, found)

		_, err := srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// add chain params
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(chains.ExternalChainList()[0].ChainId),
				sample.ChainParams(chains.ExternalChainList()[1].ChainId),
				sample.ChainParams(chains.ExternalChainList()[2].ChainId),
			},
		})

		msg2 := types.MsgRemoveChainParams{
			Creator: admin,
			ChainId: chains.ExternalChainList()[3].ChainId,
		}

		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg2, nil)

		// not found if chain ID not in list
		_, err = srv.RemoveChainParams(sdk.WrapSDKContext(ctx), &msg2)
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})
}
