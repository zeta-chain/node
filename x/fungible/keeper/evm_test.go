package keeper_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc1967proxy.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/wzeta.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/node/e2e/contracts/dapp"
	"github.com/zeta-chain/node/e2e/contracts/dappreverter"
	"github.com/zeta-chain/node/e2e/contracts/example"
	"github.com/zeta-chain/node/e2e/contracts/reverter"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/ptr"
	"github.com/zeta-chain/node/server/config"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

// get a valid chain id independently of the build flag
func getValidChainID(t *testing.T) int64 {
	list := chains.DefaultChainsList()
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
	mockEVMKeeper *keepertest.FungibleMockEVMKeeper,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	return deploySystemContracts(t, ctx, k, mockEVMKeeper)
}

// deploy upgradable gateway contract and return its address
func deployGatewayContract(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk types.EVMKeeper,
	wzeta, admin common.Address,
) common.Address {
	// Deploy the gateway contract
	implAddr, err := k.DeployContract(ctx, gatewayzevm.GatewayZEVMMetaData)
	require.NoError(t, err)
	require.NotEmpty(t, implAddr)
	assertContractDeployment(t, evmk, ctx, implAddr)

	// Deploy the proxy contract
	gatewayABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	require.NoError(t, err)

	// Encode the initializer data
	initializerData, err := gatewayABI.Pack("initialize", wzeta, admin)
	require.NoError(t, err)

	gatewayContract, err := k.DeployContract(ctx, erc1967proxy.ERC1967ProxyMetaData, implAddr, initializerData)
	require.NoError(t, err)
	require.NotEmpty(t, gatewayContract)
	assertContractDeployment(t, evmk, ctx, gatewayContract)

	// store the gateway in the system contract object
	sys, found := k.GetSystemContract(ctx)
	if !found {
		sys = types.SystemContract{}
	}
	sys.Gateway = gatewayContract.Hex()
	k.SetSystemContract(ctx, sys)

	return gatewayContract
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

	k.SetGatewayGasLimit(ctx, types.DefaultGatewayGasLimit)

	// deploy the gateway contract
	contract := deployGatewayContract(t, ctx, k, evmk, wzeta, sample.EthAddress())
	require.NotEmpty(t, contract)

	return
}

type SystemContractDeployConfig struct {
	DeployWZeta            bool
	DeployUniswapV2Factory bool
	DeployUniswapV2Router  bool
}

// deploySystemContractsConfigurable deploys the system contracts and returns their addresses
// while having a possibility to skip some deployments to test different scenarios
func deploySystemContractsConfigurable(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk types.EVMKeeper,
	config *SystemContractDeployConfig,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	var err error

	if config.DeployWZeta {
		wzeta, err = k.DeployWZETA(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, wzeta)
		assertContractDeployment(t, evmk, ctx, wzeta)
	}

	if config.DeployUniswapV2Factory {
		uniswapV2Factory, err = k.DeployUniswapV2Factory(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, uniswapV2Factory)
		assertContractDeployment(t, evmk, ctx, uniswapV2Factory)
	}

	if config.DeployUniswapV2Router {
		uniswapV2Router, err = k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
		require.NoError(t, err)
		require.NotEmpty(t, uniswapV2Router)
		assertContractDeployment(t, evmk, ctx, uniswapV2Router)
	}

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
	exampleABI, err := example.ExampleMetaData.GetAbi()
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
	t.Run("should error if chain not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			987,
			coin.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			chainID,
			coin.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if deploy contract fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)
		chainID := getValidChainID(t)
		mockFailedContractDeployment(ctx, t, k)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			chainID,
			coin.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("can deploy the zrc20 contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			chainID,
			coin.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
			ptr.Ptr(sdkmath.NewUint(2000)),
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
		require.Equal(t, coin.CoinType_Gas, foreignCoins.CoinType)
		require.Equal(t, uint64(1000), foreignCoins.GasLimit)
		require.Equal(t, uint64(2000), foreignCoins.LiquidityCap.Uint64())

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

	t.Run("can deploy the zrc20 contract without a gateway address", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		systemContract, found := k.GetSystemContract(ctx)
		require.True(t, found)
		systemContract.Gateway = ""
		k.SetSystemContract(ctx, systemContract)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			chainID,
			coin.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
			ptr.Ptr(sdkmath.NewUint(2000)),
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
		require.Equal(t, coin.CoinType_Gas, foreignCoins.CoinType)
		require.Equal(t, uint64(1000), foreignCoins.GasLimit)
		require.Equal(t, uint64(2000), foreignCoins.LiquidityCap.Uint64())

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

	t.Run("can deploy the zrc20 contract with default liquidity cap", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		addr, err := k.DeployZRC20Contract(
			ctx,
			"foo",
			"bar",
			8,
			chainID,
			coin.CoinType_Gas,
			"foobar",
			big.NewInt(1000),
			nil,
		)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, addr)

		foreignCoins, found := k.GetForeignCoins(ctx, addr.Hex())
		require.True(t, found)
		require.Greater(t, foreignCoins.LiquidityCap.Uint64(), uint64(0))
	})
}

func TestKeeper_DeploySystemContracts(t *testing.T) {
	t.Run("system contract deployment should error if deploy contract fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		wzeta, uniswapV2Factory, uniswapV2Router, _, _ := deploySystemContractsWithMockEvmKeeper(
			t,
			ctx,
			k,
			mockEVMKeeper,
		)
		mockFailedContractDeployment(ctx, t, k)

		res, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, uniswapV2Router)
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("router deployment should error if deploy contract fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockFailedContractDeployment(ctx, t, k)

		res, err := k.DeployUniswapV2Router02(ctx, sample.EthAddress(), sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("wzeta deployment should error if deploy contract fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockFailedContractDeployment(ctx, t, k)

		res, err := k.DeployWZETA(ctx)
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("connector deployment should error if deploy contract fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockFailedContractDeployment(ctx, t, k)

		res, err := k.DeployConnectorZEVM(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("can deploy the system contracts", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
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
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		balance, err := k.BalanceOfZRC4(ctx, wzeta, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotNil(t, balance)
		require.Equal(t, int64(0), balance.Int64())

		amount := big.NewInt(100)
		err = sdkk.BankKeeper.MintCoins(
			ctx,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin("azeta", sdkmath.NewIntFromBigInt(amount))),
		)
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
	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		example, err := k.DeployContract(ctx, example.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, example)

		res, err := k.CallDepositAndCall(
			ctx,
			systemcontract.ZContext{
				Origin:  sample.EthAddress().Bytes(),
				Sender:  sample.EthAddress(),
				ChainID: big.NewInt(chainID),
			},
			sample.EthAddress(),
			example,
			big.NewInt(42),
			[]byte(""),
		)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should deposit and call the contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")

		exampleContract, err := k.DeployContract(ctx, example.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, exampleContract)

		res, err := k.CallDepositAndCall(
			ctx,
			systemcontract.ZContext{
				Origin:  sample.EthAddress().Bytes(),
				Sender:  sample.EthAddress(),
				ChainID: big.NewInt(chainID),
			},
			zrc20,
			exampleContract,
			big.NewInt(42),
			[]byte(""),
		)
		require.NoError(t, err)
		require.False(t, types.IsContractReverted(res, err))
		balance, err := k.BalanceOfZRC4(ctx, zrc20, exampleContract)
		require.NoError(t, err)
		require.Equal(t, int64(42), balance.Int64())

		// check onCrossChainCall has been called
		exampleABI, err := example.ExampleMetaData.GetAbi()
		require.NoError(t, err)
		res, err = k.CallEVM(
			ctx,
			*exampleABI,
			types.ModuleAddressEVM,
			exampleContract,
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
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")

		// Deploy reverter
		reverter, err := k.DeployContract(ctx, reverter.ReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, reverter)

		res, err := k.CallDepositAndCall(
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
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")

		res, err := k.CallDepositAndCall(
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
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// Deploy example
		contract, err := k.DeployContract(ctx, example.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, contract)
		abi, err := example.ExampleMetaData.GetAbi()
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
		require.Contains(t, err.Error(), fmt.Sprintf("\"revert_reason\":\"%s\"", utils.ErrHashRevertFoo))

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

		// No revert with successful call
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
		k, ctx := keepertest.FungibleKeeperAllMocks(t)

		mockAuthKeeper := keepertest.GetFungibleAccountMock(t, k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)

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
		k, ctx := keepertest.FungibleKeeperAllMocks(t)

		mockAuthKeeper := keepertest.GetFungibleAccountMock(t, k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)

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
		k, ctx := keepertest.FungibleKeeperAllMocks(t)

		mockAuthKeeper := keepertest.GetFungibleAccountMock(t, k)
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
		k, ctx := keepertest.FungibleKeeperAllMocks(t)

		mockAuthKeeper := keepertest.GetFungibleAccountMock(t, k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)

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
		k, ctx := keepertest.FungibleKeeperAllMocks(t)

		mockAuthKeeper := keepertest.GetFungibleAccountMock(t, k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)

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

func TestKeeper_DeployContract(t *testing.T) {
	t.Run("should error if pack ctor args fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		addr, err := k.DeployContract(ctx, zrc20.ZRC20MetaData, "")
		require.ErrorIs(t, err, types.ErrABIGet)
		require.Empty(t, addr)
	})

	t.Run("should error if metadata bin empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		metadata := &bind.MetaData{
			ABI: wzeta.WETH9MetaData.ABI,
			Bin: "",
		}
		addr, err := k.DeployContract(ctx, metadata)
		require.ErrorIs(t, err, types.ErrABIGet)
		require.Empty(t, addr)
	})

	t.Run("should error if metadata cant be decoded", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		metadata := &bind.MetaData{
			ABI: wzeta.WETH9MetaData.ABI,
			Bin: "0x1",
		}
		addr, err := k.DeployContract(ctx, metadata)
		require.ErrorIs(t, err, types.ErrABIPack)
		require.Empty(t, addr)
	})

	t.Run("should error if module acc not set up", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		addr, err := k.DeployContract(ctx, wzeta.WETH9MetaData)
		require.Error(t, err)
		require.Empty(t, addr)
	})
}

func TestKeeper_QueryProtocolFlatFee(t *testing.T) {
	t.Run("should error if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.QueryProtocolFlatFee(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.QueryProtocolFlatFee(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return fee", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		fee := big.NewInt(42)
		protocolFlatFee, err := zrc20ABI.Methods["PROTOCOL_FLAT_FEE"].Outputs.Pack(fee)
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: protocolFlatFee})

		res, err := k.QueryProtocolFlatFee(ctx, sample.EthAddress())
		require.NoError(t, err)
		require.Equal(t, fee, res)
	})
}

func TestKeeper_QueryGasLimit(t *testing.T) {
	t.Run("should error if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.QueryGasLimit(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.QueryGasLimit(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return gas limit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		limit := big.NewInt(42)
		gasLimit, err := zrc20ABI.Methods["GAS_LIMIT"].Outputs.Pack(limit)
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: gasLimit})

		res, err := k.QueryGasLimit(ctx, sample.EthAddress())
		require.NoError(t, err)
		require.Equal(t, limit, res)
	})
}

func TestKeeper_QueryChainIDFromContract(t *testing.T) {
	t.Run("should error if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.QueryChainIDFromContract(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.QueryChainIDFromContract(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		chainId := big.NewInt(42)
		chainIdFromContract, err := zrc20ABI.Methods["GAS_LIMIT"].Outputs.Pack(chainId)
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: chainIdFromContract})

		res, err := k.QueryChainIDFromContract(ctx, sample.EthAddress())
		require.NoError(t, err)
		require.Equal(t, chainId, res)
	})
}

func TestKeeper_TotalSupplyZRC4(t *testing.T) {
	t.Run("should error if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.TotalSupplyZRC4(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.TotalSupplyZRC4(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return total supply", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		supply := big.NewInt(42)
		supplyFromContract, err := zrc20ABI.Methods["GAS_LIMIT"].Outputs.Pack(supply)
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: supplyFromContract})

		res, err := k.TotalSupplyZRC4(ctx, sample.EthAddress())
		require.NoError(t, err)
		require.Equal(t, supply, res)
	})
}

func TestKeeper_BalanceOfZRC4(t *testing.T) {
	t.Run("should error if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.BalanceOfZRC4(ctx, sample.EthAddress(), sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.BalanceOfZRC4(ctx, sample.EthAddress(), sample.EthAddress())
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return balance", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		balance := big.NewInt(42)
		balanceFromContract, err := zrc20ABI.Methods["GAS_LIMIT"].Outputs.Pack(balance)
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: balanceFromContract})

		res, err := k.BalanceOfZRC4(ctx, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.Equal(t, balance, res)
	})
}

func TestKeeper_QueryZRC20Data(t *testing.T) {
	t.Run("should error if evm call fails for name", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("should error if unpack fails for name", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("should error if evm call fails for symbol", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc4ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		name, err := zrc4ABI.Methods["name"].Outputs.Pack("name")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: name})

		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("should error if unpack for symbol", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc4ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		name, err := zrc4ABI.Methods["name"].Outputs.Pack("name")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: name})

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("should error if evm call fails for decimals", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc4ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		name, err := zrc4ABI.Methods["name"].Outputs.Pack("name")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: name})

		symbol, err := zrc4ABI.Methods["symbol"].Outputs.Pack("symbol")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: symbol})

		mockEVMKeeper.MockEVMFailCallOnce()

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("should error if unpack fails for decimals", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc4ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		name, err := zrc4ABI.Methods["name"].Outputs.Pack("name")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: name})

		symbol, err := zrc4ABI.Methods["symbol"].Outputs.Pack("symbol")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: symbol})

		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: []byte{}})

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("should return zrc20 data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		zrc4ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		name, err := zrc4ABI.Methods["name"].Outputs.Pack("name")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: name})

		symbol, err := zrc4ABI.Methods["symbol"].Outputs.Pack("symbol")
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: symbol})

		decimals, err := zrc4ABI.Methods["decimals"].Outputs.Pack(uint8(8))
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: decimals})

		res, err := k.QueryZRC20Data(ctx, sample.EthAddress())
		require.NoError(t, err)
		require.Equal(t, uint8(8), res.Decimals)
		require.Equal(t, "name", res.Name)
		require.Equal(t, "symbol", res.Symbol)
	})
}

func TestKeeper_CallOnReceiveZevmConnector(t *testing.T) {
	t.Run("should call on receive on connector which calls onZetaMessage on sample DAPP", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		senderAddress := sample.EthAddress().Bytes()
		sourceChainID := big.NewInt(1)
		destinationAddress := dAppContract
		zetaValue := big.NewInt(45)
		data := []byte("message")
		internalSendHash := [32]byte{}

		_, err = k.CallOnReceiveZevmConnector(
			ctx,
			senderAddress,
			sourceChainID,
			destinationAddress,
			zetaValue,
			data,
			internalSendHash,
		)
		require.NoError(t, err)

		require.NoError(t, err)

		dappAbi, err := dapp.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*dappAbi,
			types.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"zetaTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("zetaTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, senderAddress, valSenderAddress)
	})

	t.Run("should error if system contract not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		dAppContract, err := k.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		_, err = k.CallOnReceiveZevmConnector(ctx,
			sample.EthAddress().Bytes(),
			big.NewInt(1),
			dAppContract,
			big.NewInt(45), []byte("message"), [32]byte{})
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should error in contract call reverts", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, dappreverter.DappReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		_, err = k.CallOnReceiveZevmConnector(ctx,
			sample.EthAddress().Bytes(),
			big.NewInt(1),
			dAppContract,
			big.NewInt(45), []byte("message"), [32]byte{})
		require.ErrorContains(t, err, "execution reverted")
	})
}

func TestKeeper_CallOnRevertZevmConnector(t *testing.T) {
	t.Run("should call on revert on connector which calls onZetaRevert on sample DAPP", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)
		senderAddress := dAppContract
		sourceChainID := big.NewInt(1)
		destinationAddress := sample.EthAddress().Bytes()
		destinationChainID := big.NewInt(1)
		zetaValue := big.NewInt(45)
		data := []byte("message")
		internalSendHash := [32]byte{}
		_, err = k.CallOnRevertZevmConnector(
			ctx,
			senderAddress,
			sourceChainID,
			destinationAddress,
			destinationChainID,
			zetaValue,
			data,
			internalSendHash,
		)
		require.NoError(t, err)

		dappAbi, err := dapp.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*dappAbi,
			types.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"zetaTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("zetaTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, senderAddress.Bytes(), valSenderAddress)
	})

	t.Run("should error if system contract not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		dAppContract, err := k.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		_, err = k.CallOnRevertZevmConnector(ctx,
			dAppContract,
			big.NewInt(1),
			sample.EthAddress().Bytes(),
			big.NewInt(1),
			big.NewInt(45), []byte("message"), [32]byte{})
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should error in contract call reverts", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, dappreverter.DappReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		_, err = k.CallOnRevertZevmConnector(ctx,
			dAppContract,
			big.NewInt(1),
			sample.EthAddress().Bytes(),
			big.NewInt(1),
			big.NewInt(45), []byte("message"), [32]byte{})
		require.ErrorContains(t, err, "execution reverted")
	})
}

func TestKeeper_RefundRemainingGasFees(t *testing.T) {
	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		receiver := ethcommon.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")

		err := k.DepositChainGasToken(ctx, chainID, big.NewInt(100), receiver)
		require.Error(t, err)
	})

	t.Run("can refund remaining gas fees to receiver", func(t *testing.T) {
		// Arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		receiver := ethcommon.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		receiverBalanceBefore, err := k.BalanceOfZRC4(ctx, gasZRC20, receiver)
		require.NoError(t, err)
		require.Equal(t, int64(0), receiverBalanceBefore.Int64())

		// Act
		err = k.DepositChainGasToken(ctx, chainID, big.NewInt(100), receiver)
		require.NoError(t, err)

		// Assert
		gasZRC20, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		receiverBalanceAfter, err := k.BalanceOfZRC4(ctx, gasZRC20, receiver)
		require.NoError(t, err)
		require.Equal(t, int64(100), receiverBalanceAfter.Int64())
	})
}
