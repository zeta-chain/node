package keeper_test

import (
	"encoding/json"
	zetacommon "github.com/zeta-chain/zetacore/common"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/server/config"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	var err error

	wzeta, err = k.DeployWZETA(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wzeta)

	uniswapV2Factory, err = k.DeployUniswapV2Factory(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Factory)

	uniswapV2Router, err = k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Router)

	connector, err = k.DeployConnectorZEVM(ctx, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, connector)

	systemContract, err = k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, uniswapV2Router)
	require.NoError(t, err)
	require.NotEmpty(t, systemContract)

	return
}

func TestKeeper_DeployZRC20Contract(t *testing.T) {
	t.Run("can deploy the zrc20 contract", func(t *testing.T) {
		k, ctx, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			1,
			zetacommon.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
		)
		require.NoError(t, err)

		found, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(1))
		require.NoError(t, err)
		require.Equal(t, addr, found)
	})
}

func TestKeeper_DeploySystemContract(t *testing.T) {
	t.Run("can deploy the system contracts", func(t *testing.T) {
		k, ctx, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy the system contracts
		wzeta, uniswapV2Factory, uniswapV2Router, _, systemContract := deploySystemContracts(t, ctx, k)

		// can find system contract address
		found, err := k.GetSystemContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, systemContract, found)

		// can find factory address
		found, err = k.GetUniswapV2FactoryAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, uniswapV2Factory, found)

		// can find router address
		found, err = k.GetUniswapV2Router02Address(ctx)
		require.NoError(t, err)
		require.Equal(t, uniswapV2Router, found)

		// can find the wzeta contract address
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
