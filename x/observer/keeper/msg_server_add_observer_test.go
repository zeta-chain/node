package keeper_test

import (
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_AddObserver(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		wctx := sdk.WrapSDKContext(ctx)
		srv := keeper.NewMsgServerImpl(*k)

		msg := types.MsgAddObserver{
			Creator: admin,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		res, err := srv.AddObserver(wctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
		require.Nil(t, res)
	})

	t.Run("should error if pub key not valid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		wctx := sdk.WrapSDKContext(ctx)
		srv := keeper.NewMsgServerImpl(*k)

		msg := types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: "invalid",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		res, err := srv.AddObserver(wctx, &msg)
		require.Error(t, err)
		require.Equal(t, &types.MsgAddObserverResponse{}, res)
	})

	t.Run("unable to add observer if observer already exists", func(t *testing.T) {
		//ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress := sample.AccAddress()
		wctx := sdk.WrapSDKContext(ctx)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
		srv := keeper.NewMsgServerImpl(*k)
		k.SetObserverSet(ctx, types.ObserverSet{ObserverList: []string{observerAddress}})

		msg := types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: sample.PubKeyString(),
			AddNodeAccountOnly:      false,
			ObserverAddress:         observerAddress,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		_, err := srv.AddObserver(wctx, &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrDuplicateObserver)
	})

	t.Run("should add observer if add node account only false", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress := sample.AccAddress()
		wctx := sdk.WrapSDKContext(ctx)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
		srv := keeper.NewMsgServerImpl(*k)

		msg := types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: sample.PubKeyString(),
			AddNodeAccountOnly:      false,
			ObserverAddress:         observerAddress,
		}
		fmt.Println("msg", msg.ZetaclientGranteePubkey)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		res, err := srv.AddObserver(wctx, &msg)
		require.NoError(t, err)
		require.Equal(t, &types.MsgAddObserverResponse{}, res)

		loc, found := k.GetLastObserverCount(ctx)
		require.True(t, found)
		require.Equal(t, uint64(1), loc.Count)
	})

	t.Run("should add to node account if add node account only true", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		observerAddress := sample.AccAddress()

		wctx := sdk.WrapSDKContext(ctx)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
		srv := keeper.NewMsgServerImpl(*k)

		_, found = k.GetKeygen(ctx)
		require.False(t, found)
		_, found = k.GetNodeAccount(ctx, observerAddress)
		require.False(t, found)

		msg := types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: sample.PubKeyString(),
			AddNodeAccountOnly:      true,
			ObserverAddress:         observerAddress,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		res, err := srv.AddObserver(wctx, &msg)
		require.NoError(t, err)
		require.Equal(t, &types.MsgAddObserverResponse{}, res)

		_, found = k.GetLastObserverCount(ctx)
		require.False(t, found)

		keygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.Equal(t, types.Keygen{BlockNumber: math.MaxInt64}, keygen)

		_, found = k.GetNodeAccount(ctx, observerAddress)
		require.True(t, found)
	})
}
