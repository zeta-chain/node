package keeper_test

import (
	"testing"
)

func TestKeeper_VerifyProof(t *testing.T) {
	//t.Run("should error if crosschain flags not found", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
	//		UseObserverMock: true,
	//	})
	//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
	//	observerMock.On("GetCrosschainFlags", mock.Anything).Return(observertypes.CrosschainFlags{}, false)
	//
	//	res, err := k.VerifyProof(ctx, &proofs.Proof{}, 5, sample.Hash().String(), 1)
	//	require.Error(t, err)
	//	require.Nil(t, res)
	//})
	//
	//t.Run("should error if verification not enabled for btc chain", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
	//		UseObserverMock: true,
	//	})
	//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
	//	observerMock.On("GetCrosschainFlags", mock.Anything).Return(observertypes.CrosschainFlags{
	//		BlockHeaderVerificationFlags: &observertypes.BlockHeaderVerificationFlags{
	//			IsBtcTypeChainEnabled: false,
	//		},
	//	}, true)
	//
	//	res, err := k.VerifyProof(ctx, &proofs.Proof{}, 18444, sample.Hash().String(), 1)
	//	require.Error(t, err)
	//	require.Nil(t, res)
	//})
	//
	//t.Run("should error if verification not enabled for evm chain", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
	//		UseObserverMock: true,
	//	})
	//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
	//	observerMock.On("GetCrosschainFlags", mock.Anything).Return(observertypes.CrosschainFlags{
	//		BlockHeaderVerificationFlags: &observertypes.BlockHeaderVerificationFlags{
	//			IsEthTypeChainEnabled: false,
	//		},
	//	}, true)
	//
	//	res, err := k.VerifyProof(ctx, &proofs.Proof{}, 5, sample.Hash().String(), 1)
	//	require.Error(t, err)
	//	require.Nil(t, res)
	//})
	//
	//t.Run("should error if block header-based verification not supported", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
	//		UseObserverMock: true,
	//	})
	//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
	//	observerMock.On("GetCrosschainFlags", mock.Anything).Return(observertypes.CrosschainFlags{
	//		BlockHeaderVerificationFlags: &observertypes.BlockHeaderVerificationFlags{
	//			IsEthTypeChainEnabled: false,
	//		},
	//	}, true)
	//
	//	res, err := k.VerifyProof(ctx, &proofs.Proof{}, 101, sample.Hash().String(), 1)
	//	require.Error(t, err)
	//	require.Nil(t, res)
	//})
	//
	//t.Run("should error if blockhash invalid", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
	//		UseObserverMock: true,
	//	})
	//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
	//	observerMock.On("GetCrosschainFlags", mock.Anything).Return(observertypes.CrosschainFlags{
	//		BlockHeaderVerificationFlags: &observertypes.BlockHeaderVerificationFlags{
	//			IsBtcTypeChainEnabled: true,
	//		},
	//	}, true)
	//
	//	res, err := k.VerifyProof(ctx, &proofs.Proof{}, 18444, "invalid", 1)
	//	require.Error(t, err)
	//	require.Nil(t, res)
	//})
	//
	//t.Run("should error if block header not found", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
	//		UseObserverMock: true,
	//	})
	//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
	//	observerMock.On("GetCrosschainFlags", mock.Anything).Return(observertypes.CrosschainFlags{
	//		BlockHeaderVerificationFlags: &observertypes.BlockHeaderVerificationFlags{
	//			IsEthTypeChainEnabled: true,
	//		},
	//	}, true)
	//
	//	observerMock.On("GetBlockHeader", mock.Anything, mock.Anything).Return(proofs.BlockHeader{}, false)
	//
	//	res, err := k.VerifyProof(ctx, &proofs.Proof{}, 5, sample.Hash().String(), 1)
	//	require.Error(t, err)
	//	require.Nil(t, res)
	//})
}
