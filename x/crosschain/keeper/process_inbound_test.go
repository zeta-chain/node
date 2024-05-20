package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_ProcessInboundZEVMDeposit(t *testing.T) {
	t.Run("process zevm deposit successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		fungibleMock.On("ZETADepositAndCallContract", mock.Anything,
			mock.Anything,
			receiver, int64(0), amount, mock.Anything, mock.Anything).Return(nil, nil)

		// call ProcessInbound
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.GetInboundParams().SenderChainId = 0
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
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

		fungibleMock.On("ZETADepositAndCallContract", mock.Anything, mock.Anything, receiver, int64(0), amount, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("deposit error"))

		// call ProcessInbound
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.GetInboundParams().SenderChainId = 0
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, "deposit error", cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit reverts fails at GetSupportedChainFromChainID", func(t *testing.T) {
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
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// mock unsuccessful GetSupportedChainFromChainID
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil)

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("invalid sender chain id %d", cctx.InboundParams.SenderChainId), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at and GetRevertGasLimit", func(t *testing.T) {
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

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("revert gas limit error: %s", types.ErrForeignCoinNotFound), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at PayGasInERC20AndUpdateCctx", func(t *testing.T) {
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock unsuccessful PayGasInERC20AndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("deposit revert message: %s err : %s", errDeposit, observertypes.ErrSupportedChains), cctx.CctxStatus.StatusMessage)
	})

	t.Run("uunable to process zevm deposit HandleEVMDeposit revert fails at PayGasInERC20AndUpdateCctx with gas limit is 0", func(t *testing.T) {
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 0)

		// mock unsuccessful PayGasInERC20AndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("deposit revert message: %s err : %s", errDeposit, observertypes.ErrSupportedChains), cctx.CctxStatus.StatusMessage)
	})

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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// Mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, "cannot find receiver chain nonce")
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

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)
		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, *senderChain)

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, errDeposit.Error(), cctx.CctxStatus.StatusMessage)
		require.Equal(t, updatedNonce, cctx.GetCurrentOutboundParam().TssNonce)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails as the cctx has already been reverted", func(t *testing.T) {
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaPrivnetChain.ChainId
		cctx.OutboundParams = append(cctx.OutboundParams, cctx.GetCurrentOutboundParam())
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, fmt.Sprintf("revert outbound error: %s", "cannot revert a revert tx"))
	})
}

func TestKeeper_ProcessInboundProcessCrosschainMsgPassing(t *testing.T) {
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
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *receiverChain, "")

		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, *receiverChain)

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.Equal(t, updatedNonce, cctx.GetCurrentOutboundParam().TssNonce)
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
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, receiverChain.ChainId).
			Return(nil).Once()

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, observertypes.ErrSupportedChains.Error(), cctx.CctxStatus.StatusMessage)
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
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *receiverChain, "")

		// mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, receiverChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		// call ProcessInbound
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessInbound(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, "cannot find receiver chain nonce")
	})
}
