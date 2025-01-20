package bank

import (
	"math/big"
	"testing"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	evmkeeper "github.com/zeta-chain/ethermint/x/evm/keeper"
	"github.com/zeta-chain/ethermint/x/evm/statedb"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/contracts/erc1967proxy"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"
)

func Test_Methods(t *testing.T) {
	t.Run("should fail when trying to run deposit as read only method", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		methodID := ts.bankABI.Methods[DepositMethodName]

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(0)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, true)
		require.ErrorIs(
			t,
			precompiletypes.ErrWriteMethod{
				Method: "deposit",
			},
			err)
		require.Empty(t, success)
	})

	t.Run("should fail when trying to run withdraw as read only method", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		methodID := ts.bankABI.Methods[WithdrawMethodName]

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(0)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, true)
		require.ErrorIs(
			t,
			precompiletypes.ErrWriteMethod{
				Method: "withdraw",
			},
			err)
		require.Empty(t, success)
	})

	t.Run("should fail when caller has 0 token balance", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		methodID := ts.bankABI.Methods[DepositMethodName]

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(0)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.Error(t, err)
		require.ErrorAs(
			t,
			precompiletypes.ErrInvalidAmount{
				Got: "0",
			},
			err,
		)

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail when bank has 0 token allowance", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))

		methodID := ts.bankABI.Methods[DepositMethodName]

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(1000)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 0")

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail when trying to deposit 0", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))

		methodID := ts.bankABI.Methods[DepositMethodName]

		// Allow bank to spend 500 ZRC20 tokens.
		allowBank(t, ts, big.NewInt(500))

		// Set CallerAddress and evm.Origin to the caller address.
		// Caller does not have any balance, and bank does not have any allowance.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Set the input arguments for the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(0)}...,
		)

		_, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid token amount: 0")
	})

	t.Run("should fail when trying to deposit more than allowed to bank", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))

		// Allow bank to spend 500 ZRC20 tokens.
		allowBank(t, ts, big.NewInt(500))

		// Set CallerAddress and evm.Origin to the caller address.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Prepare and call the deposit method.
		methodID := ts.bankABI.Methods[DepositMethodName]

		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(501)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.Error(t, err)
		require.Contains(
			t,
			err.Error(),
			"unexpected error in LockZRC20InBank: failed allowance check: invalid allowance, got 500, wanted 501",
		)

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err := ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok := resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.EqualValues(t, big.NewInt(0).Uint64(), balance.Uint64())
	})

	t.Run("should fail when trying to deposit more than user balance", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))

		// Allow bank to spend 500 ZRC20 tokens.
		allowBank(t, ts, big.NewInt(1000))

		// Set CallerAddress and evm.Origin to the caller address.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Prepare and call the deposit method.
		methodID := ts.bankABI.Methods[DepositMethodName]

		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(1001)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.Error(t, err)
		require.ErrorAs(
			t,
			precompiletypes.ErrInvalidAmount{
				Got: "1000",
			},
			err,
		)

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err := ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok := resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.EqualValues(t, big.NewInt(0).Uint64(), balance.Uint64())
	})

	t.Run("should deposit tokens and retrieve balance of cosmos coin", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))
		methodID := ts.bankABI.Methods[DepositMethodName]

		// Allow bank to spend 500 ZRC20 tokens.
		allowBank(t, ts, big.NewInt(500))

		// Set CallerAddress and evm.Origin to the caller address.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Prepare and call the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(500)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.NoError(t, err)

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.True(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err := ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok := resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.Equal(t, big.NewInt(500), balance)
	})

	t.Run("should deposit tokens, withdraw and check with balanceOf", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))
		methodID := ts.bankABI.Methods[DepositMethodName]

		// Allow bank to spend 500 ZRC20 tokens.
		allowBank(t, ts, big.NewInt(500))

		// Set CallerAddress and evm.Origin to the caller address.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Prepare and call the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(500)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.NoError(t, err)

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.True(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err := ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok := resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.Equal(t, big.NewInt(500).Uint64(), balance.Uint64())

		// Prepare and call the withdraw method.
		methodID = ts.bankABI.Methods[WithdrawMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(500)}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.NoError(t, err)

		res, err = ts.bankABI.Methods[WithdrawMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok = res[0].(bool)
		require.True(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err = ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok = resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.Equal(t, big.NewInt(0).Uint64(), balance.Uint64())
	})

	t.Run("should deposit tokens and fail when withdrawing more than depositted", func(t *testing.T) {
		ts := setupChain(t)
		caller := fungibletypes.ModuleAddressEVM
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, caller, big.NewInt(1000))
		methodID := ts.bankABI.Methods[DepositMethodName]

		// Allow bank to spend 500 ZRC20 tokens.
		allowBank(t, ts, big.NewInt(500))

		// Set CallerAddress and evm.Origin to the caller address.
		ts.mockVMContract.CallerAddress = caller
		ts.mockEVM.Origin = caller

		// Prepare and call the deposit method.
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(500)}...,
		)

		success, err := ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.NoError(t, err)

		res, err := ts.bankABI.Methods[DepositMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.True(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err := ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok := resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.Equal(t, big.NewInt(500).Uint64(), balance.Uint64())

		// Prepare and call the withdraw method.
		methodID = ts.bankABI.Methods[WithdrawMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, big.NewInt(501)}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		require.Error(t, err)
		require.ErrorAs(
			t,
			precompiletypes.ErrInsufficientBalance{
				Requested: "501",
				Got:       "500",
			},
			err,
		)

		res, err = ts.bankABI.Methods[WithdrawMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		ok = res[0].(bool)
		require.False(t, ok)

		// Prepare and call the balanceOf method.
		methodID = ts.bankABI.Methods[BalanceOfMethodName]
		ts.mockVMContract.Input = packInputArgs(
			t,
			methodID,
			[]interface{}{ts.zrc20Address, caller}...,
		)

		success, err = ts.bankContract.Run(ts.mockEVM, ts.mockVMContract, false)
		resultBalanceOf, err = ts.bankABI.Methods[BalanceOfMethodName].Outputs.Unpack(success)
		require.NoError(t, err)

		balance, ok = resultBalanceOf[0].(*big.Int)
		require.True(t, ok)
		require.Equal(t, big.NewInt(500).Uint64(), balance.Uint64())
	})
}

/*
	Test utils.
*/

type testSuite struct {
	ctx            sdk.Context
	fungibleKeeper *fungiblekeeper.Keeper
	sdkKeepers     keeper.SDKKeepers
	bankContract   *Contract
	bankABI        abi.ABI
	mockEVM        *vm.EVM
	mockVMContract *vm.Contract
	zrc20Address   common.Address
	zrc20ABI       abi.ABI
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

	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return testSuite{}
	}

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
		uint256.NewInt(0),
		0,
	)

	return testSuite{
		ctx,
		fungibleKeeper,
		sdkKeepers,
		contract,
		abi,
		mockEVM,
		mockVMContract,
		zrc20Address,
		*zrc20ABI,
	}
}

func allowBank(t *testing.T, ts testSuite, amount *big.Int) {
	resAllowance, err := callEVM(
		t,
		ts.ctx,
		ts.fungibleKeeper,
		&ts.zrc20ABI,
		fungibletypes.ModuleAddressEVM,
		ts.zrc20Address,
		"approve",
		[]interface{}{ts.bankContract.Address(), amount},
	)
	require.NoError(t, err, "error allowing bank to spend ZRC20 tokens")

	allowed, ok := resAllowance[0].(bool)
	require.True(t, ok)
	require.True(t, allowed)
}

func callEVM(
	t *testing.T,
	ctx sdk.Context,
	fungibleKeeper *fungiblekeeper.Keeper,
	abi *abi.ABI,
	from common.Address,
	dst common.Address,
	method string,
	args []interface{},
) ([]interface{}, error) {
	res, err := fungibleKeeper.CallEVM(
		ctx,           // ctx
		*abi,          // abi
		from,          // from
		dst,           // to
		big.NewInt(0), // value
		nil,           // gasLimit
		true,          // commit
		true,          // noEthereumTxEvent
		method,        // method
		args...,       // args
	)
	require.NoError(t, err, "CallEVM error")
	require.Equal(t, "", res.VmError, "res.VmError should be empty")

	ret, err := abi.Methods[method].Outputs.Unpack(res.Ret)
	require.NoError(t, err, "Unpack error")

	return ret, nil
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
