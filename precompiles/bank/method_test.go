package bank

import (
	"math/big"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	evmkeeper "github.com/zeta-chain/ethermint/x/evm/keeper"
	"github.com/zeta-chain/ethermint/x/evm/statedb"
	"github.com/zeta-chain/node/pkg/chains"
	erc1967proxy "github.com/zeta-chain/node/pkg/contracts/erc1967proxy"
	ptypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	gatewayzevm "github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
)

func Test_Deposit(t *testing.T) {
	t.Run("should fail when caller has 0 token balance", func(t *testing.T) {
		ts := setupChain(t)

		methodID := ts.abi.Methods[DepositMethodName]

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = fungibletypes.ModuleAddressEVM
		ts.mockEVM.Origin = fungibletypes.ModuleAddressEVM

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(0)}...,
		)

		_, err := ts.contract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.ErrorAs(
			t,
			ptypes.ErrInvalidAmount{
				Got: "0",
			},
			err,
		)
	})

	t.Run("should fail when bank has 0 token allowance", func(t *testing.T) {
		ts := setupChain(t)
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(1000))

		methodID := ts.abi.Methods[DepositMethodName]

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = fungibletypes.ModuleAddressEVM
		ts.mockEVM.Origin = fungibletypes.ModuleAddressEVM

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(0)}...,
		)

		_, err := ts.contract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.ErrorAs(
			t,
			ptypes.ErrInvalidAmount{
				Got: "0",
			},
			err,
		)
	})
}

/*
	Functions to set up the test environment.
*/

type testSuite struct {
	ctx              sdk.Context
	fungibleKeeper   *fungiblekeeper.Keeper
	sdkKeepers       keeper.SDKKeepers
	contract         *Contract
	abi              abi.ABI
	mockEVM          *vm.EVM
	mockVMContract   *vm.Contract
	zrc20Address     common.Address
}

func setupChain(t *testing.T) testSuite {
	// Initialize basic parameters to mock the chain.
	fungibleKeeper, ctx, sdkKeepers, _ := keeper.FungibleKeeper(t)
	chainID := getValidChainID(t)

	// Make sure the account store is initialized.
	// This is completely needed for accounts to be created in the state.
	fungibleKeeper.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

	// Deploy system contracts in order to deploy a ZRC20 token.
	deploySystemContracts(t, ctx, fungibleKeeper, *sdkKeepers.EvmKeeper)
	zrc20Address := setupGasCoin(t, ctx, fungibleKeeper, sdkKeepers.EvmKeeper, chainID, "ZRC20", "ZRC20")

	// Keepers and chain configuration.
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	gasConfig := storetypes.TransientGasConfig()

	// Create the bank contract.
	contract := NewIBankContract(ctx, sdkKeepers.BankKeeper, *fungibleKeeper, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIBankContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	address := contract.Address()
	require.NotNil(t, address, "contract address should not be nil")

	mockEVM := vm.NewEVM(
		vm.BlockContext{},
		vm.TxContext{},
		statedb.New(ctx, sdkKeepers.EvmKeeper, statedb.TxConfig{}),
		&params.ChainConfig{},
		vm.Config{},
	)

	mockVMContract := vm.NewContract(
		contractRef{address: common.Address{}},
		contractRef{address: ContractAddress},
		big.NewInt(0),
		0,
	)

	return  testSuite{
		ctx,
		fungibleKeeper,
		sdkKeepers,
		contract,
		abi,
		mockEVM,
		mockVMContract,
		zrc20Address,
	}
}

// setupGasCoin is a helper function to setup the gas coin for testing
func setupGasCoin(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	chainID int64,
	assetName string,
	symbol string,
) (zrc20 common.Address) {
	addr, err := k.SetupChainGasCoinAndPool(
		ctx,
		chainID,
		assetName,
		symbol,
		8,
		nil,
	)
	require.NoError(t, err)
	assertContractDeployment(t, *evmk, ctx, addr)
	return addr
}

// get a valid chain id independently of the build flag
func getValidChainID(t *testing.T) int64 {
	list := chains.DefaultChainsList()
	require.NotEmpty(t, list)
	require.NotNil(t, list[0])
	return list[0].ChainId
}

// require that a contract has been deployed by checking stored code is non-empty.
func assertContractDeployment(t *testing.T, k evmkeeper.Keeper, ctx sdk.Context, contractAddress common.Address) {
	acc := k.GetAccount(ctx, contractAddress)
	require.NotNil(t, acc)
	code := k.GetCode(ctx, common.BytesToHash(acc.CodeHash))
	require.NotEmpty(t, code)
}

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk evmkeeper.Keeper,
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

	// deploy the gateway contract
	contract := deployGatewayContract(t, ctx, k, &evmk, wzeta, sample.EthAddress())
	require.NotEmpty(t, contract)

	return
}

// deploy upgradable gateway contract and return its address
func deployGatewayContract(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	wzeta, admin common.Address,
) common.Address {
	// Deploy the gateway contract
	implAddr, err := k.DeployContract(ctx, gatewayzevm.GatewayZEVMMetaData)
	require.NoError(t, err)
	require.NotEmpty(t, implAddr)
	assertContractDeployment(t, *evmk, ctx, implAddr)

	// Deploy the proxy contract
	gatewayABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	require.NoError(t, err)

	// Encode the initializer data
	initializerData, err := gatewayABI.Pack("initialize", wzeta, admin)
	require.NoError(t, err)

	gatewayContract, err := k.DeployContract(ctx, erc1967proxy.ERC1967ProxyMetaData, implAddr, initializerData)
	require.NoError(t, err)
	require.NotEmpty(t, gatewayContract)
	assertContractDeployment(t, *evmk, ctx, gatewayContract)

	// store the gateway in the system contract object
	sys, found := k.GetSystemContract(ctx)
	if !found {
		sys = fungibletypes.SystemContract{}
	}
	sys.Gateway = gatewayContract.Hex()
	k.SetSystemContract(ctx, sys)

	return gatewayContract
}

func packInputArgs(t *testing.T, methodID abi.Method, args ...interface{}) []byte {
	input, err := methodID.Inputs.Pack(args...)
	require.NoError(t, err)
	return append(methodID.ID, input...)
}

type contractRef struct {
	address common.Address
}

func (c contractRef) Address() common.Address {
	return c.address
}
