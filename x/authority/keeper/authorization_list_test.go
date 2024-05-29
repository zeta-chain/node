package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestKeeper_GetAuthorizationList(t *testing.T) {
	t.Run("successfully get authorizations list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		authorizationList := sample.AuthorizationList("sample")
		require.NoError(t, k.SetAuthorizationList(ctx, authorizationList))
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, list)
	})

	t.Run("get authorizations list not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		_, found := k.GetAuthorizationList(ctx)
		require.False(t, found)
	})
}

func TestKeeper_SetAuthorizationList(t *testing.T) {
	t.Run("successfully set authorizations list when a list already exists", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		authorizationList := sample.AuthorizationList("sample")
		require.NoError(t, k.SetAuthorizationList(ctx, authorizationList))

		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, list)

		newAuthorizationList := sample.AuthorizationList("sample2")
		require.NotEqual(t, authorizationList, newAuthorizationList)
		require.NoError(t, k.SetAuthorizationList(ctx, newAuthorizationList))

		list, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, newAuthorizationList, list)
	})

	t.Run("unable to set invalid authorizations list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		authorizationsList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
		}}

		require.ErrorIs(t, k.SetAuthorizationList(ctx, authorizationsList), types.ErrInValidAuthorizationList)
	})
}
