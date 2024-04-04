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

func TestMsgServer_RemoveFromOutTxTracker(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: 1,
			Nonce:   1,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)

		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.RemoveFromOutTxTracker(ctx, &types.MsgRemoveFromOutTxTracker{
			Creator: admin,
		})
		require.Error(t, err)
		require.Empty(t, res)

		_, found := k.GetOutTxTracker(ctx, 1, 1)
		require.True(t, found)
	})

	t.Run("should remove if authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: 1,
			Nonce:   1,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.RemoveFromOutTxTracker(ctx, &types.MsgRemoveFromOutTxTracker{
			Creator: admin,
			ChainId: 1,
			Nonce:   1,
		})
		require.NoError(t, err)
		require.Empty(t, res)

		_, found := k.GetOutTxTracker(ctx, 1, 1)
		require.False(t, found)
	})
}
