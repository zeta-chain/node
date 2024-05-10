package keeper_test

import (
	"math"
	"testing"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_AddObserver(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)
		wctx := sdk.WrapSDKContext(ctx)

		srv := keeper.NewMsgServerImpl(*k)
		res, err := srv.AddObserver(wctx, &types.MsgAddObserver{
			Creator: admin,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if pub key not valid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		wctx := sdk.WrapSDKContext(ctx)

		srv := keeper.NewMsgServerImpl(*k)
		res, err := srv.AddObserver(wctx, &types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: "invalid",
		})
		require.Error(t, err)
		require.Equal(t, &types.MsgAddObserverResponse{}, res)
	})

	t.Run("should add if add node account only false", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		admin := sample.AccAddress()
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		wctx := sdk.WrapSDKContext(ctx)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
		srv := keeper.NewMsgServerImpl(*k)
		observerAddress := sdk.AccAddress(crypto.AddressHash([]byte("ObserverAddress")))
		res, err := srv.AddObserver(wctx, &types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: sample.PubKeyString(),
			AddNodeAccountOnly:      false,
			ObserverAddress:         observerAddress.String(),
		})
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
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		wctx := sdk.WrapSDKContext(ctx)

		_, found := k.GetLastObserverCount(ctx)
		require.False(t, found)
		srv := keeper.NewMsgServerImpl(*k)
		observerAddress := sdk.AccAddress(crypto.AddressHash([]byte("ObserverAddress")))
		_, found = k.GetKeygen(ctx)
		require.False(t, found)
		_, found = k.GetNodeAccount(ctx, observerAddress.String())
		require.False(t, found)

		res, err := srv.AddObserver(wctx, &types.MsgAddObserver{
			Creator:                 admin,
			ZetaclientGranteePubkey: sample.PubKeyString(),
			AddNodeAccountOnly:      true,
			ObserverAddress:         observerAddress.String(),
		})
		require.NoError(t, err)
		require.Equal(t, &types.MsgAddObserverResponse{}, res)

		_, found = k.GetLastObserverCount(ctx)
		require.False(t, found)

		keygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.Equal(t, types.Keygen{BlockNumber: math.MaxInt64}, keygen)

		_, found = k.GetNodeAccount(ctx, observerAddress.String())
		require.True(t, found)
	})
}
