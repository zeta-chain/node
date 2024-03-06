package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_AbortStuckCCTX(t *testing.T) {
	t.Run("can abort a cctx in pending inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingInbound,
			StatusMessage: "pending inbound",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		_, err := msgServer.AbortStuckCCTX(ctx, &crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		})

		require.NoError(t, err)
		cctxFound, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("cctx_index"))
		require.True(t, found)
		require.Equal(t, crosschaintypes.CctxStatus_Aborted, cctxFound.CctxStatus.Status)
		require.Equal(t, crosschainkeeper.AbortMessage, cctxFound.CctxStatus.StatusMessage)
	})

	t.Run("can abort a cctx in pending outbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingOutbound,
			StatusMessage: "pending outbound",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		_, err := msgServer.AbortStuckCCTX(ctx, &crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		})

		require.NoError(t, err)
		cctxFound, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("cctx_index"))
		require.True(t, found)
		require.Equal(t, crosschaintypes.CctxStatus_Aborted, cctxFound.CctxStatus.Status)
		require.Equal(t, crosschainkeeper.AbortMessage, cctxFound.CctxStatus.StatusMessage)
	})

	t.Run("can abort a cctx in pending revert", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingRevert,
			StatusMessage: "pending revert",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		_, err := msgServer.AbortStuckCCTX(ctx, &crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		})

		require.NoError(t, err)
		cctxFound, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("cctx_index"))
		require.True(t, found)
		require.Equal(t, crosschaintypes.CctxStatus_Aborted, cctxFound.CctxStatus.Status)
		require.Equal(t, crosschainkeeper.AbortMessage, cctxFound.CctxStatus.StatusMessage)
	})

	t.Run("cannot abort a cctx in pending outbound if not admin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, false)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingOutbound,
			StatusMessage: "pending outbound",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		_, err := msgServer.AbortStuckCCTX(ctx, &crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		})
		require.ErrorIs(t, err, observertypes.ErrNotAuthorized)
	})

	t.Run("cannot abort a cctx if doesn't exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// abort the cctx
		_, err := msgServer.AbortStuckCCTX(ctx, &crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		})
		require.ErrorIs(t, err, crosschaintypes.ErrCannotFindCctx)
	})

	t.Run("cannot abort a cctx if not pending", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_OutboundMined,
			StatusMessage: "outbound mined",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		_, err := msgServer.AbortStuckCCTX(ctx, &crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		})
		require.ErrorIs(t, err, crosschaintypes.ErrStatusNotPending)
	})
}
