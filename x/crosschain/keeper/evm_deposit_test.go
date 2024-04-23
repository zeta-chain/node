package keeper_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
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
		sender := sample.EthAddress()
		senderChainId := int64(0)

		// expect DepositCoinZeta to be called
		fungibleMock.On("ZETADepositAndCallContract", ctx, ethcommon.HexToAddress(sender.String()), receiver, senderChainId, amount, mock.Anything, mock.Anything).Return(nil, nil)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = senderChainId
		cctx.InboundTxParams.Sender = sender.String()
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
		sender := sample.EthAddress()
		senderChainId := int64(0)
		amount := big.NewInt(42)
		cctx := sample.CrossChainTx(t, "foo")
		// expect DepositCoinZeta to be called
		errDeposit := errors.New("deposit failed")
		fungibleMock.On("ZETADepositAndCallContract", ctx, ethcommon.HexToAddress(sender.String()), receiver, senderChainId, amount, mock.Anything, mock.Anything).Return(nil, errDeposit)
		// call HandleEVMDeposit

		cctx.InboundTxParams.Sender = sender.String()
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = senderChainId
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

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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

	t.Run("should error on processing ERC20 deposit calling fungible method for contract call if process logs fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{
			Logs: []*evmtypes.Log{
				{
					Address:     receiver.Hex(),
					Topics:      []string{},
					Data:        []byte{},
					BlockNumber: uint64(ctx.BlockHeight()),
					TxHash:      sample.Hash().Hex(),
					TxIndex:     1,
					BlockHash:   sample.Hash().Hex(),
					Index:       1,
				},
			},
		}, true, nil)

		fungibleMock.On("GetSystemContract", mock.Anything).Return(fungibletypes.SystemContract{}, false)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.InboundTxParams.TxOrigin = ""
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
		cctx.GetInboundTxParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.Error(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("can process ERC20 deposit calling fungible method for contract call if process logs doesnt fail", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{
			Logs: []*evmtypes.Log{
				{
					Address:     receiver.Hex(),
					Topics:      []string{},
					Data:        []byte{},
					BlockNumber: uint64(ctx.BlockHeight()),
					TxHash:      sample.Hash().Hex(),
					TxIndex:     1,
					BlockHash:   sample.Hash().Hex(),
					Index:       1,
				},
			},
		}, true, nil)

		fungibleMock.On("GetSystemContract", mock.Anything).Return(fungibletypes.SystemContract{
			ConnectorZevm: sample.EthAddress().Hex(),
		}, true)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.InboundTxParams.TxOrigin = ""
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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

	t.Run("should error if invalid sender", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.InboundTxParams.TxOrigin = ""
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundTxParams().Sender = "invalid"
		cctx.GetInboundTxParams().SenderChainId = 987
		cctx.RelayedMessage = ""
		cctx.GetInboundTxParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.Error(t, err)
		require.False(t, reverted)
	})

	t.Run("should return error with non-reverted if deposit ERC20 fails with tx non-failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, errDeposit)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.InboundTxParams.CoinType = coin.CoinType_ERC20
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

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrForeignCoinCapReached)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrPausedZRC20)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrCallNonContract)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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
		senderChain := getValidEthChainID()

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = sample.EthAddress().String()
		cctx.GetInboundTxParams().Amount = math.NewUint(42)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_Gas
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

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		require.NoError(t, err)
		cctx := sample.CrossChainTx(t, "foo")
		b, err := cctx.Marshal()
		require.NoError(t, err)
		ctx = ctx.WithTxBytes(b)
		fungibleMock.On(
			"ZRC20DepositAndCallContract",
			ctx,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			data,
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		cctx.GetCurrentOutTxParam().Receiver = sample.EthAddress().String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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
		require.Equal(t, uint64(ctx.BlockHeight()), cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight)
	})

	t.Run("should deposit into receiver with specified data if no address parsed with data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

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
			coin.CoinType_ERC20,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundTxParams().CoinType = coin.CoinType_ERC20
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
