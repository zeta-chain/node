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
		k.SetAuthorizationList(ctx, authorizationList)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, list)
	})

	t.Run("get authorizations list not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		list, found := k.GetAuthorizationList(ctx)
		require.False(t, found)
		require.Equal(t, types.AuthorizationList{}, list)
	})
}

func TestKeeper_SetAuthorizationList(t *testing.T) {
	t.Run("successfully set authorizations list when a list already exists", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		authorizationList := sample.AuthorizationList("sample")
		k.SetAuthorizationList(ctx, authorizationList)

		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, list)

		newAuthorizationList := sample.AuthorizationList("sample2")
		require.NotEqual(t, authorizationList, newAuthorizationList)
		k.SetAuthorizationList(ctx, newAuthorizationList)

		list, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, newAuthorizationList, list)
	})
}
