package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgServer_HandleEVMDeposit(t *testing.T) {
	t.Run("can process Zeta deposit calling fungible method", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		fungibleMock.On("DepositCoinZeta", ctx, receiver, amount).Return(nil)

		// call HandleEVMDeposit
		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_Zeta,
			},
			nil,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with non-reverted if deposit Zeta fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		errDeposit := errors.New("deposit failed")
		fungibleMock.On("DepositCoinZeta", ctx, receiver, amount).Return(errDeposit)

		// call HandleEVMDeposit
		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_Zeta,
			},
			nil,
		)
		require.ErrorIs(t, err, errDeposit)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("can process ERC20 deposit calling fungible method", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		// ZRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.CoinType, msg.Asset)
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, nil)

		// call HandleEVMDeposit
		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_ERC20,
				Message:  "",
				Asset:    "",
			},
			senderChain,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with non-reverted if deposit ERC20 fails with tx non-failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		// ZRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.CoinType, msg.Asset)
		errDeposit := errors.New("deposit failed")
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, errDeposit)

		// call HandleEVMDeposit
		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_ERC20,
				Message:  "",
				Asset:    "",
			},
			senderChain,
		)
		require.ErrorIs(t, err, errDeposit)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with tx failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		// ZRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.CoinType, msg.Asset)
		errDeposit := errors.New("deposit failed")
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, errDeposit)

		// call HandleEVMDeposit
		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_ERC20,
				Message:  "",
				Asset:    "",
			},
			senderChain,
		)
		require.ErrorIs(t, err, errDeposit)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with liquidity cap reached", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		// ZRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.CoinType, msg.Asset)
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, fungibletypes.ErrForeignCoinCapReached)

		// call HandleEVMDeposit
		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_ERC20,
				Message:  "",
				Asset:    "",
			},
			senderChain,
		)
		require.ErrorIs(t, err, fungibletypes.ErrForeignCoinCapReached)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})
}
