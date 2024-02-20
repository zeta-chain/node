package keeper_test

import (
	"encoding/json"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	zetacommon "github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/server/config"
	"github.com/zeta-chain/zetacore/testutil/contracts"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// get a valid chain id independently of the build flag
func getValidChainID(t *testing.T) int64 {
	list := zetacommon.DefaultChainsList()
	require.NotEmpty(t, list)
	require.NotNil(t, list[0])

	return list[0].ChainId
}

// require that a contract has been deployed by checking stored code is non-empty.
func assertContractDeployment(t *testing.T, k types.EVMKeeper, ctx sdk.Context, contractAddress common.Address) {
	acc := k.GetAccount(ctx, contractAddress)
	require.NotNil(t, acc)

	code := k.GetCode(ctx, common.BytesToHash(acc.CodeHash))
	require.NotEmpty(t, code)
}

func deploySystemContractsWithMockEvmKeeper(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	mockEVMKeeper *testkeeper.FungibleMockEVMKeeper,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	return deploySystemContracts(t, ctx, k, mockEVMKeeper)
}

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk types.EVMKeeper,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	var err error

	wzeta, err = k.DeployWZETA(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wzeta)
	assertContractDeployment(t, evmk, ctx, wzeta)

	uniswapV2Factory, err = k.DeployUniswapV2Factory(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Factory)
	assertContractDeployment(t, evmk, ctx, uniswapV2Factory)

	uniswapV2Router, err = k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Router)
	assertContractDeployment(t, evmk, ctx, uniswapV2Router)

	connector, err = k.DeployConnectorZEVM(ctx, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, connector)
	assertContractDeployment(t, evmk, ctx, connector)

	systemContract, err = k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, uniswapV2Router)
	require.NoError(t, err)
	require.NotEmpty(t, systemContract)
	assertContractDeployment(t, evmk, ctx, systemContract)

	return
}

// assertExampleBarValue asserts value Bar of the contract Example, used to test onCrossChainCall
func assertExampleBarValue(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	address common.Address,
	expected int64,
) {
	exampleABI, err := contracts.ExampleMetaData.GetAbi()
	require.NoError(t, err)
	res, err := k.CallEVM(
		ctx,
		*exampleABI,
		types.ModuleAddressEVM,
		address,
		big.NewInt(0),
		nil,
		false,
		false,
		"bar",
	)
	require.NoError(t, err)
	unpacked, err := exampleABI.Unpack("bar", res.Ret)
	require.NoError(t, err)
	require.NotZero(t, len(unpacked))
	bar, ok := unpacked[0].(*big.Int)
	require.True(t, ok)
	require.Equal(t, big.NewInt(expected), bar)
}

func TestKeeper_DeployZRC20Contract(t *testing.T) {
	t.Run("can deploy the zrc20 contract", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			chainID,
			zetacommon.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
		)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, addr)

		// check foreign coin
		foreignCoins, found := k.GetForeignCoins(ctx, addr.Hex())
		require.True(t, found)
		require.Equal(t, "foobar", foreignCoins.Asset)
		require.Equal(t, chainID, foreignCoins.ForeignChainId)
		require.Equal(t, uint32(8), foreignCoins.Decimals)
		require.Equal(t, "foo", foreignCoins.Name)
		require.Equal(t, "bar", foreignCoins.Symbol)
		require.Equal(t, zetacommon.CoinType_Gas, foreignCoins.CoinType)
		require.Equal(t, uint64(1000), foreignCoins.GasLimit)

		// can get the zrc20 data
		zrc20Data, err := k.QueryZRC20Data(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, "foo", zrc20Data.Name)
		require.Equal(t, "bar", zrc20Data.Symbol)
		require.Equal(t, uint8(8), zrc20Data.Decimals)

		// can deposit tokens
		to := sample.EthAddress()
		oldBalance, err := k.BalanceOfZRC4(ctx, addr, to)
		require.NoError(t, err)
		require.NotNil(t, oldBalance)
		require.Equal(t, int64(0), oldBalance.Int64())

		amount := big.NewInt(100)
		_, err = k.DepositZRC20(ctx, addr, to, amount)
		require.NoError(t, err)

		newBalance, err := k.BalanceOfZRC4(ctx, addr, to)
		require.NoError(t, err)
		require.NotNil(t, newBalance)
		require.Equal(t, amount.Int64(), newBalance.Int64())
	})
}

func TestKeeper_DeploySystemContract(t *testing.T) {
	t.Run("can deploy the system contracts", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy the system contracts
		wzeta, uniswapV2Factory, uniswapV2Router, _, systemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

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

	t.Run("can deposit into wzeta", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		balance, err := k.BalanceOfZRC4(ctx, wzeta, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotNil(t, balance)
		require.Equal(t, int64(0), balance.Int64())

		amount := big.NewInt(100)
		err = sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("azeta", sdk.NewIntFromBigInt(amount))))
		require.NoError(t, err)

		err = k.CallWZetaDeposit(ctx, types.ModuleAddressEVM, amount)
		require.NoError(t, err)

		balance, err = k.BalanceOfZRC4(ctx, wzeta, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotNil(t, balance)
		require.Equal(t, amount.Int64(), balance.Int64())
	})
}

func TestKeeper_DepositZRC20AndCallContract(t *testing.T) {
	t.Run("should deposit and call the contract", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")

		example, err := k.DeployContract(ctx, contracts.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, example)

		res, err := k.DepositZRC20AndCallContract(
			ctx,
			systemcontract.ZContext{
				Origin:  sample.EthAddress().Bytes(),
				Sender:  sample.EthAddress(),
				ChainID: big.NewInt(chainID),
			},
			zrc20,
			example,
			big.NewInt(42),
			[]byte(""),
		)
		require.NoError(t, err)
		require.False(t, types.IsContractReverted(res, err))
		balance, err := k.BalanceOfZRC4(ctx, zrc20, example)
		require.NoError(t, err)
		require.Equal(t, int64(42), balance.Int64())

		// check onCrossChainCall has been called
		exampleABI, err := contracts.ExampleMetaData.GetAbi()
		require.NoError(t, err)
		res, err = k.CallEVM(
			ctx,
			*exampleABI,
			types.ModuleAddressEVM,
			example,
			big.NewInt(0),
			nil,
			false,
			false,
			"bar",
		)
		require.NoError(t, err)
		unpacked, err := exampleABI.Unpack("bar", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		bar, ok := unpacked[0].(*big.Int)
		require.True(t, ok)
		require.Equal(t, big.NewInt(42), bar)
	})

	t.Run("should return a revert error when the underlying contract call revert", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")

		// Deploy reverter
		reverter, err := k.DeployContract(ctx, contracts.ReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, reverter)

		res, err := k.DepositZRC20AndCallContract(
			ctx,
			systemcontract.ZContext{
				Origin:  sample.EthAddress().Bytes(),
				Sender:  sample.EthAddress(),
				ChainID: big.NewInt(chainID),
			},
			zrc20,
			reverter,
			big.NewInt(42),
			[]byte(""),
		)
		require.True(t, types.IsContractReverted(res, err))
		balance, err := k.BalanceOfZRC4(ctx, zrc20, reverter)
		require.NoError(t, err)
		require.Zero(t, balance.Int64())
	})

	t.Run("should revert if the underlying contract doesn't exist", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")

		res, err := k.DepositZRC20AndCallContract(
			ctx,
			systemcontract.ZContext{
				Origin:  sample.EthAddress().Bytes(),
				Sender:  sample.EthAddress(),
				ChainID: big.NewInt(chainID),
			},
			zrc20,
			sample.EthAddress(),
			big.NewInt(42),
			[]byte(""),
		)
		require.True(t, types.IsContractReverted(res, err))
	})
}

func TestKeeper_CallEVMWithData(t *testing.T) {
	t.Run("should return a revert error when the contract call revert", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// Deploy example
		contract, err := k.DeployContract(ctx, contracts.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, contract)
		abi, err := contracts.ExampleMetaData.GetAbi()
		require.NoError(t, err)

		// doRevert make contract reverted
		res, err := k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			contract,
			big.NewInt(0),
			nil,
			true,
			false,
			"doRevert",
		)
		require.Nil(t, res)
		require.True(t, types.IsContractReverted(res, err))

		// check reason is included for revert error
		// 0xbfb4ebcf is the hash of "Foo()"
		require.Contains(t, err.Error(), "reason: 0xbfb4ebcf")

		res, err = k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			contract,
			big.NewInt(0),
			nil,
			true,
			false,
			"doRevertWithMessage",
		)
		require.Nil(t, res)
		require.True(t, types.IsContractReverted(res, err))

		res, err = k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			contract,
			big.NewInt(0),
			nil,
			true,
			false,
			"doRevertWithRequire",
		)
		require.Nil(t, res)
		require.True(t, types.IsContractReverted(res, err))

		// Not a revert error if another type of error
		res, err = k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			contract,
			big.NewInt(0),
			nil,
			true,
			false,
			"doNotExist",
		)
		require.Nil(t, res)
		require.Error(t, err)
		require.False(t, types.IsContractReverted(res, err))
		require.NotContains(t, err.Error(), "reason:")

		// No revert with successfull call
		res, err = k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			contract,
			big.NewInt(0),
			nil,
			true,
			false,
			"doSucceed",
		)
		require.NotNil(t, res)
		require.NoError(t, err)
		require.False(t, types.IsContractReverted(res, err))
	})

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
		mockEVMKeeper.MockEVMSuccessCallOnce()

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
			false,
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
		mockEVMKeeper.MockEVMSuccessCallOnce()

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
			false,
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
			false,
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
			false,
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
		mockEVMKeeper.MockEVMFailCallOnce()

		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		// Call the method
		_, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			&contractAddress,
			data,
			true,
			false,
			big.NewInt(100),
			nil,
		)
		require.ErrorIs(t, err, sample.ErrSample)
	})
}
