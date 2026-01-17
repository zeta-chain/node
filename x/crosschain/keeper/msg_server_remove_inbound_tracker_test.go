package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgServer_RemoveInboundTracker(t *testing.T) {
	t.Run("fail if creator is not admin", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		nonAdmin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		txHash := "hash"
		chainID := int64(1)
		msg := types.MsgRemoveInboundTracker{
			Creator: nonAdmin,
			ChainId: chainID,
			TxHash:  txHash,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		// Act
		_, err := msgServer.RemoveInboundTracker(ctx, &msg)

		// Assert
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("successfully remove inbound tracker", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		txHash := "hash"
		chainID := int64(1)
		msg := types.MsgRemoveInboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  txHash,
		}
		k.SetInboundTracker(ctx, types.InboundTracker{ChainId: chainID, TxHash: txHash, CoinType: coin.CoinType_Gas})
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// Act
		_, err := msgServer.RemoveInboundTracker(ctx, &msg)

		// Assert
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("do nothing if inbound tracker does not exist", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		txHash := "hash"
		chainID := int64(1)
		msg := types.MsgRemoveInboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  txHash,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// Act
		_, err := msgServer.RemoveInboundTracker(ctx, &msg)

		// Assert
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})
}
