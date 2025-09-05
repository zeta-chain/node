package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_AddInboundTracker(t *testing.T) {
	t.Run("fail normal user submit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		nonAdmin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))

		msg := types.MsgAddInboundTracker{
			Creator:  nonAdmin,
			ChainId:  chainID,
			TxHash:   txHash,
			CoinType: coin.CoinType_Zeta,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("fail for unsupported chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, false)

		msg := types.MsgAddInboundTracker{
			Creator:  sample.AccAddress(),
			ChainId:  chainID + 1,
			TxHash:   txHash,
			CoinType: coin.CoinType_Zeta,
		}
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("admin add tx tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))

		setSupportedChain(ctx, zk, chainID)

		msg := types.MsgAddInboundTracker{
			Creator:  admin,
			ChainId:  chainID,
			TxHash:   txHash,
			CoinType: coin.CoinType_Zeta,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.AddInboundTracker(ctx, &msg)

		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})

	t.Run("observer is able to add inbound tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		msg := types.MsgAddInboundTracker{
			Creator:  admin,
			ChainId:  chainID,
			TxHash:   txHash,
			CoinType: coin.CoinType_Zeta,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})
}
