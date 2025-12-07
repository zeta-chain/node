package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_RemoveObserver(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)

		msg := types.MsgRemoveObserver{
			Creator:         admin,
			ObserverAddress: observerAddress,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		// ACT
		res, err := srv.RemoveObserver(ctx, &msg)

		// ASSERT
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
		require.Nil(t, res)
	})

	t.Run("should remove observer from set and node account", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)

		k.SetObserverSet(ctx, types.ObserverSet{ObserverList: []string{observerAddress}})
		k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: 1})
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: observerAddress,
		})

		require.True(t, k.IsAddressPartOfObserverSet(ctx, observerAddress))
		_, found := k.GetNodeAccount(ctx, observerAddress)
		require.True(t, found)
		loc, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(1), loc.Count)

		msg := types.MsgRemoveObserver{
			Creator:         admin,
			ObserverAddress: observerAddress,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		res, err := srv.RemoveObserver(ctx, &msg)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, &types.MsgRemoveObserverResponse{}, res)
		require.False(t, k.IsAddressPartOfObserverSet(ctx, observerAddress))
		_, found = k.GetNodeAccount(ctx, observerAddress)
		require.False(t, found)
		loc, found = k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(0), loc.Count)
	})

	t.Run("should handle removing observer when observer set has multiple observers", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress1 := sample.AccAddress()
		observerAddress2 := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)

		k.SetObserverSet(ctx, types.ObserverSet{ObserverList: []string{observerAddress1, observerAddress2}})
		k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: 2})
		k.SetNodeAccount(ctx, types.NodeAccount{Operator: observerAddress1})
		k.SetNodeAccount(ctx, types.NodeAccount{Operator: observerAddress2})

		msg := types.MsgRemoveObserver{
			Creator:         admin,
			ObserverAddress: observerAddress1,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		res, err := srv.RemoveObserver(ctx, &msg)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, &types.MsgRemoveObserverResponse{}, res)
		require.False(t, k.IsAddressPartOfObserverSet(ctx, observerAddress1))
		require.True(t, k.IsAddressPartOfObserverSet(ctx, observerAddress2))
		_, found := k.GetNodeAccount(ctx, observerAddress1)
		require.False(t, found)
		_, found = k.GetNodeAccount(ctx, observerAddress2)
		require.True(t, found)
		loc, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(1), loc.Count)
	})

	t.Run("should handle removing non-existent observer gracefully", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)

		k.SetObserverSet(ctx, types.ObserverSet{ObserverList: []string{}})

		msg := types.MsgRemoveObserver{
			Creator:         admin,
			ObserverAddress: observerAddress,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		res, err := srv.RemoveObserver(ctx, &msg)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, &types.MsgRemoveObserverResponse{}, res)
	})
}
