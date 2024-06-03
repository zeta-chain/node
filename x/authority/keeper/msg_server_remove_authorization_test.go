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

func TestMsgServer_RemoveAuthorization(t *testing.T) {
	t.Run("successfully remove authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
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
	})

	t.Run("unable to remove authorization if creator is not the correct policy", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
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
	})

	t.Run("unable to remove authorization if authorization list does not exist", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		url := types.OperationPolicyMessages[0]

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  url,
		}

		_, err := msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrAuthorizationListNotFound)
	})

	t.Run("unable to remove authorization if authorization does not exist", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
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
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)

		authorizationListNew, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, authorizationListNew)
	})

	t.Run("unable to remove authorization if authorization list is invalid", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := keepertest.SetAdminPolices(ctx, k)
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
		}}
		k.SetAuthorizationList(ctx, authorizationList)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := &types.MsgRemoveAuthorization{
			Creator: admin,
			MsgUrl:  "ABC",
		}

		_, err := msgServer.RemoveAuthorization(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrInvalidAuthorizationList)
	})
}
