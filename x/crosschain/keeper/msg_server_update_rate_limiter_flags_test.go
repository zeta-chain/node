package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgServer_UpdateRateLimiterFlags(t *testing.T) {
	t.Run("can update rate limiter flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		flags := sample.RateLimiterFlags()
		msg := types.NewMsgUpdateRateLimiterFlags(
			admin,
			flags,
		)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		_, found := k.GetRateLimiterFlags(ctx)
		require.False(t, found)

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
		msg := types.NewMsgUpdateRateLimiterFlags(
			admin,
			flags,
		)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		_, err := msgServer.UpdateRateLimiterFlags(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
