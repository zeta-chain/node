package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_AbortStuckCCTX(t *testing.T) {
	t.Run("can abort a cctx in pending inbound", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingInbound,
			StatusMessage: "pending inbound",
		}
		cctx.GetCurrentOutboundParam().TssNonce = 1

		k.SetCrossChainTx(ctx, *cctx)
		k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  0,
			NonceHigh: 10,
			ChainId:   cctx.GetCurrentOutboundParam().ReceiverChainId,
			Tss:       cctx.GetCurrentOutboundParam().TssPubkey,
		})

		// abort the cctx
		msg := crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		// Act
		_, err := msgServer.AbortStuckCCTX(ctx, &msg)

		// Assert
		require.NoError(t, err)
		cctxFound, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("cctx_index"))
		require.True(t, found)
		require.Equal(t, crosschaintypes.CctxStatus_Aborted, cctxFound.CctxStatus.Status)
		require.Contains(t, cctxFound.CctxStatus.StatusMessage, crosschainkeeper.AbortMessage)
		pendingNonces, found := k.GetObserverKeeper().
			GetPendingNonces(ctx, cctx.GetCurrentOutboundParam().TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId)
		require.True(t, found)
		require.Equal(t, pendingNonces.NonceLow, int64(cctx.GetCurrentOutboundParam().TssNonce+1))
	})

	t.Run("can abort a cctx in pending outbound", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingOutbound,
			StatusMessage: "pending outbound",
		}
		cctx.GetCurrentOutboundParam().TssNonce = 1

		k.SetCrossChainTx(ctx, *cctx)
		k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  0,
			NonceHigh: 10,
			ChainId:   cctx.GetCurrentOutboundParam().ReceiverChainId,
			Tss:       cctx.GetCurrentOutboundParam().TssPubkey,
		})

		msg := crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// Act
		_, err := msgServer.AbortStuckCCTX(ctx, &msg)

		// Assert
		require.NoError(t, err)
		cctxFound, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("cctx_index"))
		require.True(t, found)
		require.Equal(t, crosschaintypes.CctxStatus_Aborted, cctxFound.CctxStatus.Status)
		require.Contains(t, cctxFound.CctxStatus.StatusMessage, crosschainkeeper.AbortMessage)
		// ensure the last update timestamp is updated
		require.Equal(t, cctxFound.CctxStatus.LastUpdateTimestamp, ctx.BlockTime().Unix())
		pendingNonces, found := k.GetObserverKeeper().
			GetPendingNonces(ctx, cctx.GetCurrentOutboundParam().TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId)
		require.True(t, found)
		require.Equal(t, pendingNonces.NonceLow, int64(cctx.GetCurrentOutboundParam().TssNonce+1))

	})

	t.Run("can abort a cctx in pending revert", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingRevert,
			StatusMessage: "pending revert",
		}
		cctx.GetCurrentOutboundParam().TssNonce = 1

		k.SetCrossChainTx(ctx, *cctx)
		k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  0,
			NonceHigh: 10,
			ChainId:   cctx.GetCurrentOutboundParam().ReceiverChainId,
			Tss:       cctx.GetCurrentOutboundParam().TssPubkey,
		})

		// abort the cctx
		msg := crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		// Act
		_, err := msgServer.AbortStuckCCTX(ctx, &msg)

		// Assert
		require.NoError(t, err)
		cctxFound, found := k.GetCrossChainTx(ctx, sample.GetCctxIndexFromString("cctx_index"))
		require.True(t, found)
		require.Equal(t, crosschaintypes.CctxStatus_Aborted, cctxFound.CctxStatus.Status)
		require.Contains(t, cctxFound.CctxStatus.StatusMessage, crosschainkeeper.AbortMessage)
		pendingNonces, found := k.GetObserverKeeper().
			GetPendingNonces(ctx, cctx.GetCurrentOutboundParam().TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId)
		require.True(t, found)
		require.Equal(t, pendingNonces.NonceLow, int64(cctx.GetCurrentOutboundParam().TssNonce+1))
	})

	t.Run("cannot abort a cctx in pending outbound if not admin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_PendingOutbound,
			StatusMessage: "pending outbound",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		msg := crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AbortStuckCCTX(ctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("cannot abort a cctx if doesn't exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		// abort the cctx
		msg := crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.AbortStuckCCTX(ctx, &msg)
		require.ErrorIs(t, err, crosschaintypes.ErrCannotFindCctx)
	})

	t.Run("cannot abort a cctx if not pending", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		// create a cctx
		cctx := sample.CrossChainTx(t, "cctx_index")
		cctx.CctxStatus = &crosschaintypes.Status{
			Status:        crosschaintypes.CctxStatus_OutboundMined,
			StatusMessage: "outbound mined",
		}
		k.SetCrossChainTx(ctx, *cctx)

		// abort the cctx
		msg := crosschaintypes.MsgAbortStuckCCTX{
			Creator:   admin,
			CctxIndex: sample.GetCctxIndexFromString("cctx_index"),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.AbortStuckCCTX(ctx, &msg)
		require.ErrorIs(t, err, crosschaintypes.ErrStatusNotPending)
	})
}
