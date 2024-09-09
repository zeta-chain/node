package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/keeper"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMsgServer_AddAuthorization(t *testing.T) {
	const url = "/zetachain.zetacore.sample.ABC"
	var AddAuthorization = types.Authorization{
		MsgUrl:           "/zetachain.zetacore.authority.MsgAddAuthorization",
		AuthorizedPolicy: types.PolicyType_groupAdmin,
	}
	t.Run("successfully add authorization of type admin to existing authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		msgServer := keeper.NewMsgServerImpl(*k)
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)

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
		require.Equal(t, prevLen+1, len(authorizationList.Authorizations))
	})

	t.Run("successfully add authorization of type operational to existing authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)

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
		require.Equal(t, prevLen+1, len(authorizationList.Authorizations))
	})

	t.Run("successfully add authorization of type emergency to existing authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)

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
		require.Equal(t, prevLen+1, len(authorizationList.Authorizations))
	})

	t.Run(
		"successfully add authorization to list containing only authorization for AddAuthorization",
		func(t *testing.T) {
			k, ctx := keepertest.AuthorityKeeper(t)
			admin := keepertest.SetAdminPolicies(ctx, k)
			k.SetAuthorizationList(ctx, types.AuthorizationList{
				Authorizations: []types.Authorization{
					AddAuthorization,
				},
			})
			msgServer := keeper.NewMsgServerImpl(*k)

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
			require.Equal(t, 2, len(authorizationList.Authorizations))
		},
	)

	t.Run("unable to add authorization to empty authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.AuthorizationList{})
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrUnauthorized)
	})

	t.Run("update existing authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "/zetachain.zetacore.sample.ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			AddAuthorization,
		},
		}
		k.SetAuthorizationList(ctx, authorizationList)
		prevLen := len(authorizationList.Authorizations)

		msgServer := keeper.NewMsgServerImpl(*k)

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
		require.Equal(t, prevLen, len(authorizationList.Authorizations))
	})

	t.Run("fail to add authorization with invalid policy as creator", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := &types.MsgAddAuthorization{
			Creator:          sample.AccAddress(),
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrUnauthorized)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, prevLen, len(authorizationList.Authorizations))
	})

	// This scenario is not possible as the authorization list is always valid.But it is good to have in case the validation logic is changed in the future
	t.Run("fail to set invalid authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           url,
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           url,
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			AddAuthorization,
		}}
		k.SetAuthorizationList(ctx, authorizationList)
		prevLen := len(authorizationList.Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := &types.MsgAddAuthorization{
			Creator:          admin,
			MsgUrl:           url,
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}

		_, err := msgServer.AddAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrInvalidAuthorizationList)

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, prevLen, len(authorizationList.Authorizations))
	})
}
