package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/keeper"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgServer_AddAuthorization(t *testing.T) {
	t.Run("successfully add authorization of type admin to existing authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		policy, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupAdmin, policy)
	})

	t.Run("successfully add authorization of type operational to existing authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupOperational,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		policy, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupOperational, policy)
	})

	t.Run("successfully add authorization of type emergency to existing authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupEmergency,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		policy, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupEmergency, policy)
	})

	t.Run("successfully add authorization to empty authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.AuthorizationList{})
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		policy, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupAdmin, policy)
	})

	t.Run("successfully set authorization when list is not found ", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		policy, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupAdmin, policy)
	})

	t.Run("update existing authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "/zetachain.zetacore.sample.ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		},
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		policy, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupAdmin, policy)
	})

	t.Run("fail to add authorization with invalid policy as creator", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          sample.AccAddress(),
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrUnauthorized)
	})

	// This scenario is not possible as the authorization list is always valid.But it is good to have in case the validation logic is changed in the future
	t.Run("fail to set invalid authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "/zetachain.zetacore.sample.ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           "/zetachain.zetacore.sample.ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		}})
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "/zetachain.zetacore.sample.ABC"

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrInvalidAuthorizationList)
	})
}
