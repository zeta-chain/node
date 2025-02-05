package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v2 "github.com/zeta-chain/node/x/authority/migrations/v2"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("set authorization list", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		list := types.DefaultAuthorizationsList()
		list.RemoveAuthorization("/zetachain.zetacore.crosschain.MsgRemoveInboundTracker")
		k.SetAuthorizationList(ctx, list)

		// Act
		err := v2.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultAuthorizationsList(), list)
	})
}
