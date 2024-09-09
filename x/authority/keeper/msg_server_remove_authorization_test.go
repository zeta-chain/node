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

func TestMsgServer_RemoveAuthorization(t *testing.T) {
	var removeAuthorization = types.Authorization{
		MsgUrl:           "/zetachain.zetacore.authority.MsgRemoveAuthorization",
		AuthorizedPolicy: types.PolicyType_groupAdmin,
	}
	t.Run("successfully remove operational policy authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := types.OperationPolicyMessages[0]

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  url,
		}

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err, types.ErrAuthorizationNotFound)

		_, err = msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err = authorizationList.GetAuthorizedPolicy(url)
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)
		require.Equal(t, prevLen-1, len(authorizationList.Authorizations))
	})

	t.Run("successfully remove admin policy authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := types.AdminPolicyMessages[0]

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  url,
		}

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err, types.ErrAuthorizationNotFound)

		_, err = msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err = authorizationList.GetAuthorizedPolicy(url)
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)
		require.Equal(t, prevLen-1, len(authorizationList.Authorizations))
	})

	t.Run("successfully remove emergency policy authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := types.EmergencyPolicyMessages[0]

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  url,
		}

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err, types.ErrAuthorizationNotFound)

		_, err = msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		authorizationList, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err = authorizationList.GetAuthorizedPolicy(url)
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)
		require.Equal(t, prevLen-1, len(authorizationList.Authorizations))
	})

	t.Run("unable to remove authorization if creator is not the correct policy", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
		prevLen := len(types.DefaultAuthorizationsList().Authorizations)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := types.OperationPolicyMessages[0]

		msg := &types.MsgRemoveAuthorization{
			Creator: sample.AccAddress(),
			MsgUrl:  url,
		}

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err := authorizationList.GetAuthorizedPolicy(url)
		require.NoError(t, err, types.ErrAuthorizationNotFound)

		_, err = msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrUnauthorized)

		authorizationList, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultAuthorizationsList(), authorizationList)
		require.Equal(t, prevLen, len(authorizationList.Authorizations))
	})

	t.Run("unable to remove authorization if authorization list does not exist", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := types.OperationPolicyMessages[0]

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  url,
		}

		_, err := msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorContains(t, err, types.ErrAuthorizationListNotFound.Error())
	})

	t.Run("unable to remove authorization if authorization does not exist", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		}}
		k.SetAuthorizationList(ctx, authorizationList)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := "invalid"

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  url,
		}

		authorizationList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		_, err := authorizationList.GetAuthorizedPolicy(url)
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)

		_, err = msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorContains(t, err, types.ErrAuthorizationNotFound.Error())

		authorizationListNew, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, authorizationListNew)
	})

	t.Run("unable to remove authorization if authorization list is invalid", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolicies(ctx, k)
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			removeAuthorization,
		}}
		k.SetAuthorizationList(ctx, authorizationList)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  "ABC",
		}

		_, err := msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorContains(t, err, types.ErrInvalidAuthorizationList.Error())

		authorizationListNew, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, authorizationListNew)
	})
}
