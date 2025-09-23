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

func TestMsgServer_UpdateChainParams(t *testing.T) {
	t.Run("can update chain params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		chainList := chains.ExternalChainList([]chains.Chain{})

		chain1 := chainList[0].ChainId
		chain2 := chainList[1].ChainId
		chain3 := chainList[2].ChainId

		// set admin
		admin := sample.AccAddress()
		chainParams1 := sample.ChainParams(chain1)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// check list initially empty
		_, found := k.GetChainParamsList(ctx)
		require.False(t, found)

		// a new chain params can be added
		msg := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams1,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has one chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		chainParams2 := sample.ChainParams(chain2)

		// a new chain params can be added
		msg = types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams2,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has two chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 2)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])
		chainParams3 := sample.ChainParams(chain3)

		// a new chain params can be added
		msg = types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams3,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has three chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 3)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])
		require.Equal(t, chainParams3, chainParamsList.ChainParams[2])

		// chain params can be updated
		chainParams2.ConfirmationParams.SafeInboundCount = chainParams2.ConfirmationParams.SafeInboundCount + 1
		msg = types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams2,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)
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
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: sample.ChainParams(chains.ExternalChainList([]chains.Chain{})[0].ChainId),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)

		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
