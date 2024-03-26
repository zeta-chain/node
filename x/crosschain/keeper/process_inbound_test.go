package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_ProcessZEVMDeposit(t *testing.T) {
	t.Run("process zevm deposit successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		fungibleMock.On("DepositCoinZeta", mock.Anything, receiver, amount).
			Return(nil)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		k.ProcessZEVMDeposit(ctx, cctx)
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
		fungibleMock.On("DepositCoinZeta", mock.Anything, receiver, amount).
			Return(fmt.Errorf("deposit error"), false)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		k.ProcessZEVMDeposit(ctx, cctx)
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
		senderChain := getValidEthChain(t)
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// mock unsuccessful GetSupportedChainFromChainID
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("invalid sender chain id %d", cctx.InboundTxParams.SenderChainId), cctx.CctxStatus.StatusMessage)
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
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock unsuccessful GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{}, false)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
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
		senderChain := getValidEthChain(t)
		asset := ""

		// Setup expected calls
		errDeposit := fmt.Errorf("deposit failed")
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock unsuccessful PayGasInERC20AndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
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
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// Mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
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
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)
		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, *senderChain)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, errDeposit.Error(), cctx.CctxStatus.StatusMessage)
		require.Equal(t, updatedNonce, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
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
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.OutboundTxParams = append(cctx.OutboundTxParams, cctx.GetCurrentOutTxParam())
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, fmt.Sprintf("revert outbound error: %s", "cannot revert a revert tx"))
	})
}

func TestKeeper_ProcessCrosschainMsgPassing(t *testing.T) {
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
		receiverChain := getValidEthChain(t)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *receiverChain, "")

		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, *receiverChain)

		// call ProcessCrosschainMsgPassing
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessCrosschainMsgPassing(ctx, cctx)
		require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.Equal(t, updatedNonce, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
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
		receiverChain := getValidEthChain(t)

		// mock unsuccessful PayGasAndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, receiverChain.ChainId).
			Return(nil).Once()

		// call ProcessCrosschainMsgPassing
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessCrosschainMsgPassing(ctx, cctx)
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
		receiverChain := getValidEthChain(t)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *receiverChain, "")

		// mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, receiverChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		// call ProcessCrosschainMsgPassing
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessCrosschainMsgPassing(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, "cannot find receiver chain nonce")
	})
}
