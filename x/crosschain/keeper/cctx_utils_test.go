package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
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
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Zeta,
			}})
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
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
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
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should fail if no gas coin found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
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
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})

	t.Run("should fail if no asset found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, &types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
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
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
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

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		fungibleMock.On("DepositCoinZeta", mock.Anything, receiver, amount).
			Return(nil)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit returns err without reverting", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup expected calls
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		fungibleMock.On("DepositCoinZeta", mock.Anything, receiver, amount).
			Return(fmt.Errorf("deposit error"), false)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
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

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChainId := getValidEthChainID(t)
		// Setup expected calls
		errDeposit := fmt.Errorf("deposit failed")
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChainId,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChainId).
			Return(nil)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = senderChainId
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, "invalid sender chain", cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at and GetRevertGasLimit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// Setup expected calls
		errDeposit := fmt.Errorf("deposit failed")
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain.ChainId,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{}, false)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(senderChain)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = senderChain.ChainId
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = asset
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("can't get revert tx gas limit,%s", types.ErrForeignCoinNotFound), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at PayGasInERC20AndUpdateCctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// Setup expected calls
		errDeposit := fmt.Errorf("deposit failed")
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain.ChainId,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).
			Return(big.NewInt(100), nil)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(senderChain).Once()
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = senderChain.ChainId
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = asset
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("deposit revert message: %s err : %s", errDeposit, observertypes.ErrSupportedChains), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit reverts fails at PayGasInERC20AndUpdateCctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// Setup expected calls
		errDeposit := fmt.Errorf("deposit failed")
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain.ChainId,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).
			Return(big.NewInt(100), nil)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(senderChain).Once()
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = senderChain.ChainId
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = asset
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("deposit revert message: %s err : %s", errDeposit, observertypes.ErrSupportedChains), cctx.CctxStatus.StatusMessage)
	})

}
