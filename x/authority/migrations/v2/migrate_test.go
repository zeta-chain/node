package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v2 "github.com/zeta-chain/node/x/authority/migrations/v2"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("Set authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		list, found := k.GetAuthorizationList(ctx)
		require.False(t, found)
		require.Equal(t, types.AuthorizationList{}, list)

		err := v2.MigrateStore(ctx, *k)
		require.NoError(t, err)

		list, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultAuthorizationsList(), list)
	})
}
