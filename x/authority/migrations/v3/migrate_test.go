package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v3 "github.com/zeta-chain/node/x/authority/migrations/v3"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMigrateStore(t *testing.T) {
	disableFastConfirmationAuthorization := types.Authorization{
		MsgUrl:           v3.MsgURLDisableFastConfirmation,
		AuthorizedPolicy: types.PolicyType_groupEmergency,
	}

	t.Run("update authorization list", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		oldList := types.DefaultAuthorizationsList()
		oldList.RemoveAuthorization(v3.MsgURLDisableFastConfirmation)
		k.SetAuthorizationList(ctx, oldList)

		// Act
		err := v3.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		newList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)

		// two lists should be equal if adds the removed authorization back
		oldList.SetAuthorization(disableFastConfirmationAuthorization)
		require.Equal(t, oldList, newList)
	})

	t.Run("set default authorization list if list is not found", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		// Act
		err := v3.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultAuthorizationsList(), list)
	})

	t.Run("return error if authorization list is invalid", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		k.SetAuthorizationList(ctx, types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
		}})

		// Act
		err := v3.MigrateStore(ctx, *k)

		// Assert
		require.Error(t, err)
	})
}
