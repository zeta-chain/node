package keeper_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		assert.NoError(t, err)
		assert.False(t, reverted)
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
		assert.ErrorIs(t, err, errDeposit)
		assert.False(t, reverted)
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
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

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
		assert.NoError(t, err)
		assert.False(t, reverted)
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
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, errDeposit)

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
		assert.ErrorIs(t, err, errDeposit)
		assert.False(t, reverted)
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
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)

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
		assert.ErrorIs(t, err, errDeposit)
		assert.True(t, reverted)
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
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrForeignCoinCapReached)

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
		assert.ErrorIs(t, err, fungibletypes.ErrForeignCoinCapReached)
		assert.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with zrc20 paused", func(t *testing.T) {
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
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrPausedZRC20)

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
		assert.ErrorIs(t, err, fungibletypes.ErrPausedZRC20)
		assert.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with calling a non-contract address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrCallNonContract)

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
		assert.ErrorIs(t, err, fungibletypes.ErrCallNonContract)
		assert.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should fail if can't parse address and data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		senderChain := getValidEthChain(t)

		_, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: sample.EthAddress().String(),
				Amount:   math.NewUint(42),
				CoinType: common.CoinType_Gas,
				Message:  "not_hex",
				Asset:    "",
			},
			senderChain,
		)
		assert.ErrorIs(t, err, types.ErrUnableToParseAddress)
	})

	t.Run("should deposit into address if address is parsed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		assert.NoError(t, err)
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			data,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: sample.EthAddress().String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_ERC20,
				Message:  receiver.Hex()[2:] + "DEADBEEF",
				Asset:    "",
			},
			senderChain,
		)
		assert.NoError(t, err)
		assert.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should deposit into receiver with specified data if no address parsed with data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChain(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		assert.NoError(t, err)
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			data,
			common.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		reverted, err := k.HandleEVMDeposit(
			ctx,
			sample.CrossChainTx(t, "foo"),
			types.MsgVoteOnObservedInboundTx{
				Sender:   sample.EthAddress().String(),
				Receiver: receiver.String(),
				Amount:   math.NewUintFromBigInt(amount),
				CoinType: common.CoinType_ERC20,
				Message:  "DEADBEEF",
				Asset:    "",
			},
			senderChain,
		)
		assert.NoError(t, err)
		assert.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	// TODO: add test cases for testing logs process
	// https://github.com/zeta-chain/node/issues/1207
}
