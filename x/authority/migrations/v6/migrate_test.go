package v6_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v6 "github.com/zeta-chain/node/x/authority/migrations/v6"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("update authorization list with remove observer", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		list := types.DefaultAuthorizationsList()
		// Ensure the target authorization is missing so migration should add it
		list.RemoveAuthorization("/zetachain.zetacore.observer.MsgRemoveObserver")
		k.SetAuthorizationList(ctx, list)

		// Act
		err := v6.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)

		// After migration, default list should be restored for this item at least;
		// simplest is to compare with defaults since migration uses SetAuthorization on existing list
		// which results in same as default for that authorization
		policy, err := list.GetAuthorizedPolicy("/zetachain.zetacore.observer.MsgRemoveObserver")
		require.NoError(t, err)
		require.Equal(t, types.PolicyType_groupAdmin, policy)
	})

	t.Run("set default authorization list if list is not found", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		// Act
		err := v6.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultAuthorizationsList(), list)
	})

	t.Run("return error when list is invalid", func(t *testing.T) {
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
		err := v6.MigrateStore(ctx, *k)

		// Assert
		require.Error(t, err)
	})
}
