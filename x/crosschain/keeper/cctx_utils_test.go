package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGetRevertGasLimit(t *testing.T) {
	t.Run("should return 0 if no inbound tx params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return 0 if coin type is not gas or erc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_Zeta})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return the gas limit of the gas token", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gas := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "FOO")

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, gas, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_Gas,
			InboundTxParams: &types.InboundTxParams{
				SenderChainId: chainID,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should return the gas limit of the associated asset", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		asset := sample.EthAddress().String()
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, zrc20Addr, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_ERC20,
			InboundTxParams: &types.InboundTxParams{
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should fail if no gas coin found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_Gas,
			InboundTxParams: &types.InboundTxParams{
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for gas coin fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             common.CoinType_Gas,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_Gas,
			InboundTxParams: &types.InboundTxParams{
				SenderChainId: chainID,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})

	t.Run("should fail if no asset found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_ERC20,
			InboundTxParams: &types.InboundTxParams{
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for asset fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		asset := sample.EthAddress().String()

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             common.CoinType_ERC20,
			Asset:                asset,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			CoinType: common.CoinType_ERC20,
			InboundTxParams: &types.InboundTxParams{
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})
}

func TestGetAbortedAmount(t *testing.T) {
	amount := sdkmath.NewUint(100)
	t.Run("should return the inbound amount if outbound not present", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Amount: amount,
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, amount, a)
	})
	t.Run("should return the amount outbound amount", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{Amount: amount},
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, amount, a)
	})
	t.Run("should return the zero if outbound amount is not present and inbound is 0", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Amount: sdkmath.ZeroUint(),
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, sdkmath.ZeroUint(), a)
	})
	t.Run("should return the zero if no amounts are present", func(t *testing.T) {
		cctx := types.CrossChainTx{}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, sdkmath.ZeroUint(), a)
	})
}

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
		cctx.CoinType = common.CoinType_Zeta
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
		cctx.CoinType = common.CoinType_Zeta
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
		require.Equal(t, "invalid sender chain", cctx.CctxStatus.StatusMessage)
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
		require.Equal(t, fmt.Sprintf("can't get revert tx gas limit,%s", types.ErrForeignCoinNotFound), cctx.CctxStatus.StatusMessage)
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

func TestKeeper_GetInbound(t *testing.T) {
	t.Run("should return a cctx with correct values", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		senderChain := getValidEthChain(t)
		sender := sample.EthAddress()
		receiverChain := getValidEthChain(t)
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		intxBlockHeight := uint64(420)
		intxHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := common.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteOnObservedInboundTx{
			Creator:       creator,
			Sender:        sender.String(),
			SenderChainId: senderChain.ChainId,
			Receiver:      receiver.String(),
			ReceiverChain: receiverChain.ChainId,
			Amount:        amount,
			Message:       message,
			InTxHash:      intxHash.String(),
			InBlockHeight: intxBlockHeight,
			GasLimit:      gasLimit,
			CoinType:      cointType,
			TxOrigin:      sender.String(),
			Asset:         asset,
			EventIndex:    eventIndex,
		}
		zk.ObserverKeeper.SetTSS(ctx, tss)
		cctx := k.GetInbound(ctx, &msg)
		require.Equal(t, receiver.String(), cctx.GetCurrentOutTxParam().Receiver)
		require.Equal(t, receiverChain.ChainId, cctx.GetCurrentOutTxParam().ReceiverChainId)
		require.Equal(t, sender.String(), cctx.GetInboundTxParams().Sender)
		require.Equal(t, senderChain.ChainId, cctx.GetInboundTxParams().SenderChainId)
		require.Equal(t, amount, cctx.GetInboundTxParams().Amount)
		require.Equal(t, message, cctx.RelayedMessage)
		require.Equal(t, intxHash.String(), cctx.GetInboundTxParams().InboundTxObservedHash)
		require.Equal(t, intxBlockHeight, cctx.GetInboundTxParams().InboundTxObservedExternalHeight)
		require.Equal(t, gasLimit, cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
		require.Equal(t, asset, cctx.GetInboundTxParams().Asset)
		require.Equal(t, eventIndex, cctx.EventIndex)
		require.Equal(t, cointType, cctx.CoinType)
		require.Equal(t, uint64(0), cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
		require.Equal(t, sdkmath.ZeroUint(), cctx.GetCurrentOutTxParam().Amount)
		require.Equal(t, types.CctxStatus_PendingInbound, cctx.CctxStatus.Status)
		require.Equal(t, false, cctx.CctxStatus.IsAbortRefunded)
	})
}

func Test_IsPending(t *testing.T) {
	tt := []struct {
		status   types.CctxStatus
		expected bool
	}{
		{types.CctxStatus_PendingInbound, false},
		{types.CctxStatus_PendingOutbound, true},
		{types.CctxStatus_PendingRevert, true},
		{types.CctxStatus_Reverted, false},
		{types.CctxStatus_Aborted, false},
		{types.CctxStatus_OutboundMined, false},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("status %s", tc.status), func(t *testing.T) {
			require.Equal(t, tc.expected, crosschainkeeper.IsPending(types.CrossChainTx{CctxStatus: &types.Status{Status: tc.status}}))
		})
	}
}

func TestKeeper_SaveInbound(t *testing.T) {
	t.Run("should save the cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		k.SaveInbound(ctx, cctx)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.InboundTxParams.TxFinalizationStatus)
		require.True(t, k.IsFinalizedInbound(ctx, cctx.GetInboundTxParams().InboundTxObservedHash, cctx.GetInboundTxParams().SenderChainId, cctx.EventIndex))
		_, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
	})

	t.Run("should save the cctx and remove tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		hash := sample.Hash()
		cctx.InboundTxParams.InboundTxObservedHash = hash.String()
		k.SetInTxTracker(ctx, types.InTxTracker{
			ChainId:  senderChain.ChainId,
			TxHash:   hash.String(),
			CoinType: 0,
		})
		k.SaveInbound(ctx, cctx)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.InboundTxParams.TxFinalizationStatus)
		require.True(t, k.IsFinalizedInbound(ctx, cctx.GetInboundTxParams().InboundTxObservedHash, cctx.GetInboundTxParams().SenderChainId, cctx.EventIndex))
		_, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		_, found = k.GetInTxTracker(ctx, senderChain.ChainId, hash.String())
		require.False(t, found)
	})
}

func GetERC20Cctx(t *testing.T, receiver ethcommon.Address, senderChain common.Chain, asset string, amount *big.Int) *types.CrossChainTx {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
	cctx.GetCurrentOutTxParam().Receiver = receiver.String()
	cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
	cctx.CoinType = common.CoinType_Zeta
	cctx.GetInboundTxParams().SenderChainId = senderChain.ChainId
	cctx.GetCurrentOutTxParam().ReceiverChainId = senderChain.ChainId
	cctx.CoinType = common.CoinType_ERC20
	cctx.RelayedMessage = ""
	cctx.GetInboundTxParams().Asset = asset
	cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
	return cctx
}
