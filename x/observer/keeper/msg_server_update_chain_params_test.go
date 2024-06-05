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

func TestMsgServer_UpdateChainParams(t *testing.T) {
	t.Run("can update chain params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		chain1 := chains.ExternalChainList()[0].ChainId
		chain2 := chains.ExternalChainList()[1].ChainId
		chain3 := chains.ExternalChainList()[2].ChainId

		// set admin
		admin := sample.AccAddress()
		chainParams1 := sample.ChainParams(chain1)
		msg := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams1,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// check list initially empty
		_, found := k.GetChainParamsList(ctx)
		require.False(t, found)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// a new chain params can be added
		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check list has one chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		chainParams2 := sample.ChainParams(chain2)
		msg2 := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams2,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg2, nil)

		// a new chian params can be added
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg2)
		require.NoError(t, err)

		// check list has two chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 2)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])
		chainParams3 := sample.ChainParams(chain3)
		msg3 := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams3,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg3, nil)

		// a new chain params can be added
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg3)
		require.NoError(t, err)

		// check list has three chain params
		chainParamsList, found = k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 3)
		require.Equal(t, chainParams1, chainParamsList.ChainParams[0])
		require.Equal(t, chainParams2, chainParamsList.ChainParams[1])
		require.Equal(t, chainParams3, chainParamsList.ChainParams[2])

		msg4 := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: chainParams2,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg4, nil)

		// chain params can be updated
		chainParams2.ConfirmationCount = chainParams2.ConfirmationCount + 1
		_, err = srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg4)
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
		msg := types.MsgUpdateChainParams{
			Creator:     admin,
			ChainParams: sample.ChainParams(chains.ExternalChainList()[0].ChainId),
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		_, err := srv.UpdateChainParams(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
