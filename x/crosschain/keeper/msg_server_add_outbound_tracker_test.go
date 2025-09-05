package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func getEthereumChainID() int64 {
	return 5 // Goerli
}

func TestMsgServer_AddToOutboundTracker(t *testing.T) {
	t.Run("admin can add tracker", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		hash := sample.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))
		cctxIndex := keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		msg := types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  hash,
			Nonce:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &msg)
		require.NoError(t, err)

		// assert
		tracker, found := k.GetOutboundTracker(ctx, chainID, 0)
		require.True(t, found)
		require.Equal(t, hash, tracker.HashList[0].TxHash)

		// related cctx outbound hash is set
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		require.True(t, found)
		require.Equal(t, hash, cctx.GetCurrentOutboundParam().Hash)
	})

	t.Run("observer can add tracker", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		hash := sample.Hash().Hex()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)
		cctxIndex := keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		msg := types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  hash,
			Nonce:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &msg)
		require.NoError(t, err)

		// assert
		tracker, found := k.GetOutboundTracker(ctx, chainID, 0)
		require.True(t, found)
		require.Equal(t, hash, tracker.HashList[0].TxHash)

		// related cctx outbound hash is set
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		require.True(t, found)
		require.Equal(t, hash, cctx.GetCurrentOutboundParam().Hash)
	})

	t.Run("can add hash to existing tracker", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		existinghHash := sample.Hash().Hex()
		newHash := sample.Hash().Hex()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))
		cctxIndex := keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: chainID,
			Nonce:   42,
			HashList: []*types.TxHash{
				{
					TxHash: existinghHash,
				},
			},
		})

		msg := types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  newHash,
			Nonce:   42,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &msg)
		require.NoError(t, err)

		// assert
		tracker, found := k.GetOutboundTracker(ctx, chainID, 42)
		require.True(t, found)
		require.Len(t, tracker.HashList, 2)
		require.EqualValues(t, existinghHash, tracker.HashList[0].TxHash)
		require.EqualValues(t, newHash, tracker.HashList[1].TxHash)

		// related cctx outbound hash is set
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		require.True(t, found)
		require.Equal(t, newHash, cctx.GetCurrentOutboundParam().Hash)
	})

	t.Run("should return early if cctx not pending", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		msg := types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  sample.Hash().Hex(),
			Nonce:   0,
		}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)

		// set cctx status to outbound mined
		keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_OutboundMined, false)

		// act
		res, err := msgServer.AddOutboundTracker(ctx, &msg)
		require.NoError(t, err)

		// assert
		require.Equal(t, &types.MsgAddOutboundTrackerResponse{IsRemoved: true}, res)

		// check if tracker is removed
		_, found := k.GetOutboundTracker(ctx, chainID, 0)
		require.False(t, found)
	})

	t.Run("should error for unsupported chain", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, false)

		chainID := getEthereumChainID()

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  sample.Hash().Hex(),
			Nonce:   0,
		})

		// assert
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should error if no CctxByNonce", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, true)

		chainID := getEthereumChainID()

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  sample.Hash().Hex(),
			Nonce:   0,
		})

		// assert
		require.ErrorIs(t, err, types.ErrCannotFindCctx)
	})

	t.Run("should fail if max tracker hashes reached", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		newHash := sample.Hash().Hex()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))
		keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		hashes := make([]*types.TxHash, types.MaxOutboundTrackerHashes)
		for i := 0; i < types.MaxOutboundTrackerHashes; i++ {
			hashes[i] = &types.TxHash{
				TxHash: sample.Hash().Hex(),
			}
		}

		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId:  chainID,
			Nonce:    42,
			HashList: hashes,
		})

		msg := types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  newHash,

			Nonce: 42,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &msg)

		// assert
		require.ErrorIs(t, err, types.ErrMaxTxOutTrackerHashesReached)
	})

	t.Run("no hash added if already exist", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		chainID := getEthereumChainID()
		existinghHash := sample.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))
		keepertest.MockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: chainID,
			Nonce:   42,
			HashList: []*types.TxHash{
				{
					TxHash: existinghHash,
				},
			},
		})

		msg := types.MsgAddOutboundTracker{
			Creator: admin,
			ChainId: chainID,
			TxHash:  existinghHash,
			Nonce:   42,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// act
		_, err := msgServer.AddOutboundTracker(ctx, &msg)
		require.NoError(t, err)

		// assert
		tracker, found := k.GetOutboundTracker(ctx, chainID, 42)
		require.True(t, found)
		require.Len(t, tracker.HashList, 1)
		require.EqualValues(t, existinghHash, tracker.HashList[0].TxHash)
	})
}
