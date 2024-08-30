package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgServer_UpdateRateLimiterFlags(t *testing.T) {
	t.Run("can update rate limiter flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		flags := sample.RateLimiterFlags()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		_, found := k.GetRateLimiterFlags(ctx)
		require.False(t, found)

		msg := types.NewMsgUpdateRateLimiterFlags(
			admin,
			flags,
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateRateLimiterFlags(ctx, msg)
		require.NoError(t, err)

		storedFlags, found := k.GetRateLimiterFlags(ctx)
		require.True(t, found)
		require.Equal(t, flags, storedFlags)
	})

	t.Run("cannot update rate limiter flags if unauthorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		flags := sample.RateLimiterFlags()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msg := types.NewMsgUpdateRateLimiterFlags(
			admin,
			flags,
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UpdateRateLimiterFlags(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
