package keeper_test

import (
	"encoding/hex"
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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
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

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, errDeposit)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("can process ERC20 deposit calling fungible method", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with non-reverted if deposit ERC20 fails with tx non-failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, errDeposit)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with tx failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, errDeposit)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with liquidity cap reached", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, fungibletypes.ErrForeignCoinCapReached)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with zrc20 paused", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, fungibletypes.ErrPausedZRC20)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with calling a non-contract address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

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
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, fungibletypes.ErrCallNonContract)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should fail if can't parse address and data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		senderChain := getValidEthChainID(t)

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = sample.EthAddress().String()
		cctx.GetInboundTxParams().Amount = math.NewUint(42)
		cctx.GetInboundTxParams().CoinType = common.CoinType_Gas
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = "not_hex"
		cctx.GetInboundTxParams().Asset = ""
		_, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, types.ErrUnableToParseAddress)
	})

	t.Run("should deposit into address if address is parsed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		require.NoError(t, err)
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

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = sample.EthAddress().String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = receiver.Hex()[2:] + "DEADBEEF"
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should deposit into receiver with specified data if no address parsed with data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID(t)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		require.NoError(t, err)
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

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = common.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = "DEADBEEF"
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	// TODO: add test cases for testing logs process
	// https://github.com/zeta-chain/node/issues/1207
}
