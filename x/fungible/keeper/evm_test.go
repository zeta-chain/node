package keeper_test

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/server/config"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_DeploySystemContract(t *testing.T) {
	t.Run("deploy the system contract at the given address", func(t *testing.T) {
		// k, ctx := testkeeper.FungibleNoMocks(t)
	})
}

func TestKeeper_Deploy(t *testing.T) {
	t.Run("deploy the contracts", func(t *testing.T) {
		k, ctx, sdkk := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// can deploy wzeta
		wzeta, err := k.DeployWZETA(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, wzeta)

		// can deploy uniswap v2 factory
		uniswapV2Factory, err := k.DeployUniswapV2Factory(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, uniswapV2Factory)

		// can deploy uniswap v2 router
		uniswapV2Router, err := k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
		require.NoError(t, err)
		require.NotEmpty(t, uniswapV2Router)

		// can deploy connector zevm
		connector, err := k.DeployConnectorZEVM(ctx, wzeta)
		require.NoError(t, err)
		require.NotEmpty(t, connector)

		// can deploy system contract
		systemContract, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, uniswapV2Router)
		require.NoError(t, err)
		require.NotEmpty(t, systemContract)

		// can find system contract address
		found, err := k.GetSystemContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, systemContract, found)

		acc := sdkk.EvmKeeper.GetAccount(ctx, systemContract)
		require.NotNil(t, acc)
		code := sdkk.EvmKeeper.GetCode(ctx, common.BytesToHash(acc.CodeHash))
		_ = code

		// can find factory address
		found, err = k.GetUniswapV2Router02Address(ctx)
		require.NoError(t, err)
		require.Equal(t, uniswapV2Router, found)

		// can find the contract addresses
		found, err = k.GetWZetaContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, wzeta, found)
	})
}

func TestKeeper_CallEVMWithData(t *testing.T) {
	t.Run("apply new message without gas limit estimates gas", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountMock(t, k)
		mockEVMKeeper := testkeeper.GetFungibleEVMMock(t, k)

		// Set up values
		fromAddr := sample.EthAddress()
		contractAddress := sample.EthAddress()
		data := sample.Bytes()
		args, _ := json.Marshal(evmtypes.TransactionArgs{
			From: &fromAddr,
			To:   &contractAddress,
			Data: (*hexutil.Bytes)(&data),
		})
		gasRes := &evmtypes.EstimateGasResponse{Gas: 1000}
		msgRes := &evmtypes.MsgEthereumTxResponse{}

		// Set up mocked methods
		mockAuthKeeper.On(
			"GetSequence",
			mock.Anything,
			sdk.AccAddress(fromAddr.Bytes()),
		).Return(uint64(1), nil)
		mockEVMKeeper.On(
			"EstimateGas",
			mock.Anything,
			&evmtypes.EthCallRequest{Args: args, GasCap: config.DefaultGasCap},
		).Return(gasRes, nil)
		mockEVMKeeper.On(
			"ApplyMessage",
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(msgRes, nil)

		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.On("SetBlockBloomTransient", mock.Anything).Maybe()
		mockEVMKeeper.On("SetLogSizeTransient", mock.Anything).Maybe()
		mockEVMKeeper.On("GetLogSizeTransient", mock.Anything, mock.Anything).Maybe()

		// Call the method
		res, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			&contractAddress,
			data,
			true,
			big.NewInt(100),
			nil,
		)
		require.NoError(t, err)
		require.Equal(t, msgRes, res)

		// Assert that the expected methods were called
		mockAuthKeeper.AssertExpectations(t)
		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("apply new message with gas limit skip gas estimation", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountMock(t, k)
		mockEVMKeeper := testkeeper.GetFungibleEVMMock(t, k)

		// Set up values
		fromAddr := sample.EthAddress()
		msgRes := &evmtypes.MsgEthereumTxResponse{}

		// Set up mocked methods
		mockAuthKeeper.On(
			"GetSequence",
			mock.Anything,
			sdk.AccAddress(fromAddr.Bytes()),
		).Return(uint64(1), nil)
		mockEVMKeeper.On(
			"ApplyMessage",
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(msgRes, nil)

		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.On("SetBlockBloomTransient", mock.Anything).Maybe()
		mockEVMKeeper.On("SetLogSizeTransient", mock.Anything).Maybe()
		mockEVMKeeper.On("GetLogSizeTransient", mock.Anything, mock.Anything).Maybe()

		// Call the method
		contractAddress := sample.EthAddress()
		res, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			&contractAddress,
			sample.Bytes(),
			true,
			big.NewInt(100),
			big.NewInt(1000),
		)
		require.NoError(t, err)
		require.Equal(t, msgRes, res)

		// Assert that the expected methods were called
		mockAuthKeeper.AssertExpectations(t)
		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("GetSequence failure returns error", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountMock(t, k)
		mockAuthKeeper.On("GetSequence", mock.Anything, mock.Anything).Return(uint64(1), sample.ErrSample)

		// Call the method
		contractAddress := sample.EthAddress()
		_, err := k.CallEVMWithData(
			ctx,
			sample.EthAddress(),
			&contractAddress,
			sample.Bytes(),
			true,
			big.NewInt(100),
			nil,
		)
		require.ErrorIs(t, err, sample.ErrSample)
	})

	t.Run("EstimateGas failure returns error", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountMock(t, k)
		mockEVMKeeper := testkeeper.GetFungibleEVMMock(t, k)

		// Set up values
		fromAddr := sample.EthAddress()

		// Set up mocked methods
		mockAuthKeeper.On(
			"GetSequence",
			mock.Anything,
			sdk.AccAddress(fromAddr.Bytes()),
		).Return(uint64(1), nil)
		mockEVMKeeper.On(
			"EstimateGas",
			mock.Anything,
			mock.Anything,
		).Return(nil, sample.ErrSample)

		// Call the method
		contractAddress := sample.EthAddress()
		_, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			&contractAddress,
			sample.Bytes(),
			true,
			big.NewInt(100),
			nil,
		)
		require.ErrorIs(t, err, sample.ErrSample)
	})

	t.Run("ApplyMessage failure returns error", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountMock(t, k)
		mockEVMKeeper := testkeeper.GetFungibleEVMMock(t, k)

		// Set up values
		fromAddr := sample.EthAddress()
		contractAddress := sample.EthAddress()
		data := sample.Bytes()
		args, _ := json.Marshal(evmtypes.TransactionArgs{
			From: &fromAddr,
			To:   &contractAddress,
			Data: (*hexutil.Bytes)(&data),
		})
		gasRes := &evmtypes.EstimateGasResponse{Gas: 1000}
		msgRes := &evmtypes.MsgEthereumTxResponse{}

		// Set up mocked methods
		mockAuthKeeper.On(
			"GetSequence",
			mock.Anything,
			sdk.AccAddress(fromAddr.Bytes()),
		).Return(uint64(1), nil)
		mockEVMKeeper.On(
			"EstimateGas",
			mock.Anything,
			&evmtypes.EthCallRequest{Args: args, GasCap: config.DefaultGasCap},
		).Return(gasRes, nil)
		mockEVMKeeper.On(
			"ApplyMessage",
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(msgRes, sample.ErrSample)

		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		// Call the method
		_, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			&contractAddress,
			data,
			true,
			big.NewInt(100),
			nil,
		)
		require.ErrorIs(t, err, sample.ErrSample)
	})
}
