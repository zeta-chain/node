package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_InitiateOutboundZEVMDeposit(t *testing.T) {
	t.Run("process zevm deposit successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		fungibleMock.On("LegacyZETADepositAndCallContract", mock.Anything,
			mock.Anything,
			receiver, int64(0), amount, mock.Anything, mock.Anything).Return(nil, nil)

		// call InitiateOutbound
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
		cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.GetInboundParams().SenderChainId = 0
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_OutboundMined, newStatus)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit returns err without reverting", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// mock unsuccessful HandleEVMDeposit which does not revert

		fungibleMock.On("LegacyZETADepositAndCallContract", mock.Anything, mock.Anything, receiver, int64(0), amount, mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("deposit error"))

		// call InitiateOutbound
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
		cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.GetInboundParams().SenderChainId = 0
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.ErrorContains(t, err, "deposit error")
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, "deposit error")
	})

	t.Run(
		"unable to process zevm deposit HandleEVMDeposit reverts fails at GetSupportedChainFromChainID",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
				UseObserverMock: true,
			})

			// Setup mock data
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			observerMock := keepertest.GetCrosschainObserverMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)
			senderChain := getValidEthChain()
			errDeposit := fmt.Errorf("deposit failed")

			// Setup expected calls
			// mock unsuccessful HandleEVMDeposit which reverts, i.e returns err and isContractReverted = true
			keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

			// mock unsuccessful GetSupportedChainFromChainID
			keepertest.MockFailedGetSupportedChainFromChainID(observerMock, senderChain)

			// call InitiateOutbound
			cctx := GetERC20Cctx(t, receiver, senderChain, "", amount)
			cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
			newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
			require.NoError(t, err)
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Equal(t, types.CctxStatus_Aborted, newStatus)
			require.Contains(
				t,
				cctx.CctxStatus.ErrorMessageRevert,
				"chain not supported",
			)
		},
	)

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at GetRevertGasLimit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock unsuccessful GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{}, false)

		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)

		require.Contains(
			t,
			cctx.CctxStatus.ErrorMessageRevert,
			"GetRevertGasLimit: foreign coin not found for sender chain",
		)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at PayGasInERC20AndUpdateCctx",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
				UseObserverMock: true,
			})

			// Setup mock data
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)
			senderChain := getValidEthChain()
			asset := ""

			// Setup expected calls
			errDeposit := fmt.Errorf("deposit failed")
			keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

			observerMock := keepertest.GetCrosschainObserverMock(t, k)

			// Mock successful GetSupportedChainFromChainID
			keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

			// mock successful GetRevertGasLimit for ERC20
			keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

			// mock unsuccessful PayGasInERC20AndUpdateCctx
			keepertest.MockFailedGetSupportedChainFromChainID(observerMock, senderChain)

			// call InitiateOutbound
			cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
			cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
			newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
			require.NoError(t, err)
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Equal(t, types.CctxStatus_Aborted, newStatus)
			require.Contains(
				t,
				cctx.CctxStatus.ErrorMessageRevert,
				"chain not supported",
			)
		},
	)

	t.Run(
		"uunable to process zevm deposit HandleEVMDeposit revert fails at PayGasInERC20AndUpdateCctx with gas limit is 0",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
				UseObserverMock: true,
			})

			// Setup mock data
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)
			senderChain := getValidEthChain()
			asset := ""

			// Setup expected calls
			errDeposit := fmt.Errorf("deposit failed")
			keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

			observerMock := keepertest.GetCrosschainObserverMock(t, k)

			// Mock successful GetSupportedChainFromChainID
			keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

			// mock successful GetRevertGasLimit for ERC20
			keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 0)

			// mock unsuccessful PayGasInERC20AndUpdateCctx
			keepertest.MockFailedGetSupportedChainFromChainID(observerMock, senderChain)

			// call InitiateOutbound
			cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
			cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
			newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
			require.NoError(t, err)
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Equal(t, types.CctxStatus_Aborted, newStatus)
			require.Contains(
				t,
				cctx.CctxStatus.ErrorMessageRevert,
				"chain not supported",
			)
		},
	)

	t.Run("unable to process zevm deposit HandleEVMDeposit reverts fails at UpdateNonce", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)

		// Mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainId).
			Return(observertypes.ChainNonces{}, false)

		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessageRevert, "cannot find receiver chain nonce")
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)
		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, senderChain)

		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_PendingRevert, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, errDeposit.Error())
		require.Equal(t, updatedNonce, cctx.GetCurrentOutboundParam().TssNonce)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails as the cctx has already been reverted",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
				UseObserverMock: true,
			})

			// Setup mock data
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			observerMock := keepertest.GetCrosschainObserverMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)
			senderChain := getValidEthChain()
			asset := ""
			errDeposit := fmt.Errorf("deposit failed")

			// Setup expected calls
			// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
			keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

			// Mock successful GetSupportedChainFromChainID
			keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

			// mock successful GetRevertGasLimit for ERC20
			keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

			// call InitiateOutbound
			cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
			cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
			cctx.OutboundParams = append(cctx.OutboundParams, cctx.GetCurrentOutboundParam())
			newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
			require.NoError(t, err)
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Equal(t, types.CctxStatus_Aborted, newStatus)
			require.Contains(
				t,
				cctx.CctxStatus.ErrorMessageRevert,
				"cannot revert a revert tx",
			)
		},
	)
}

func TestKeeper_InitiateOutboundProcessCrosschainMsgPassing(t *testing.T) {
	t.Run("process crosschain msg passing successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain()

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, receiverChain, "")

		// mock successful UpdateNonce
		nonce := uint64(1)
		tss := sample.Tss()
		observerMock.On("GetChainNonces", mock.Anything, receiverChain.ChainId).
			Return(observertypes.ChainNonces{Nonce: nonce}, true)
		observerMock.On("GetTSS", mock.Anything).
			Return(tss, true)
		observerMock.On("GetPendingNonces", mock.Anything, tss.TssPubkey, mock.Anything).
			Return(observertypes.PendingNonces{NonceHigh: int64(nonce)}, true)
		observerMock.On("SetChainNonces", mock.Anything, mock.Anything)
		observerMock.On("SetPendingNonces", mock.Anything, mock.Anything)

		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, receiverChain, "", amount)
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_PendingOutbound, newStatus)
		require.Equal(t, nonce, cctx.GetCurrentOutboundParam().TssNonce)
	})

	t.Run("unable to process crosschain msg passing PayGasAndUpdateCctx fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain()

		// mock unsuccessful PayGasAndUpdateCctx
		keepertest.MockFailedGetSupportedChainFromChainID(observerMock, receiverChain)

		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, receiverChain, "", amount)
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, observertypes.ErrSupportedChains.Error())
	})

	t.Run("unable to process crosschain msg passing UpdateNonce fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain()

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, receiverChain, "")

		// mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, receiverChain.ChainId).
			Return(observertypes.ChainNonces{}, false)

		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, receiverChain, "", amount)
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.ErrorContains(t, err, "cannot find receiver chain nonce")
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, "cannot find receiver chain nonce")
	})
}

func TestKeeper_InitiateOutboundFailures(t *testing.T) {
	t.Run("should fail if chain info can not be found for receiver chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain()
		receiverChain.ChainId = 123
		// call InitiateOutbound
		cctx := GetERC20Cctx(t, receiver, receiverChain, "", amount)
		newStatus, err := k.InitiateOutbound(ctx, keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true})
		require.Error(t, err)
		require.Equal(t, types.CctxStatus_PendingInbound, newStatus)
		require.ErrorContains(t, err, "chain info not found")
	})
}
