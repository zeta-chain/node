package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgServer_RemoveFromOutboundTracker(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
			Nonce:   1,
		})

		admin := sample.AccAddress()
		msg := types.MsgRemoveOutboundTracker{
			Creator: admin,
		}
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.RemoveOutboundTracker(ctx, &msg)
		require.Error(t, err)
		require.Empty(t, res)

		_, found := k.GetOutboundTracker(ctx, 1, 1)
		require.True(t, found)
	})

	t.Run("should remove if authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
			Nonce:   1,
		})

		admin := sample.AccAddress()
		msg := types.MsgRemoveOutboundTracker{
			Creator: admin,
			ChainId: 1,
			Nonce:   1,
		}
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.RemoveOutboundTracker(ctx, &msg)
		require.NoError(t, err)
		require.Empty(t, res)

		_, found := k.GetOutboundTracker(ctx, 1, 1)
		require.False(t, found)
	})
}
