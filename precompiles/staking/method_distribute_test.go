package staking

import (
	"math/big"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	evmkeeper "github.com/zeta-chain/ethermint/x/evm/keeper"
	"github.com/zeta-chain/ethermint/x/evm/statedb"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/contracts/erc1967proxy"
	ptypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"
)

func Test_Distribute(t *testing.T) {
	t.Run("should fail to run distribute as read only method", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(0)}...,
		)

		// Call method as read only.
		result, err := s.contract.Run(s.mockEVM, s.mockVMContract, true)

		// Check error and result.
		require.ErrorIs(t, err, ptypes.ErrWriteMethod{
			Method: DistributeMethodName,
		})

		// Result is empty as the write check is done before executing distribute() function.
		// On-chain this would look like reverting, so staticcall is properly reverted.
		require.Empty(t, result)
	})

	t.Run("should fail to distribute with 0 token balance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(0)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.ErrorAs(
			t,
			ptypes.ErrInvalidAmount{
				Got: "0",
			},
			err,
		)

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute with 0 allowance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1000)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 0")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute 0 token", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(0)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid token amount: 0")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute more than allowed to staking", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(999))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1000)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 999, wanted 1000")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute more than user balance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(100000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1001)}...,
		)

		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should distribute and lock ZRC20 under the bank account", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1000)}...,
		)

		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.NoError(t, err)

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.True(t, ok)

		// balance, err := s.fungibleKeeper.ZRC20BalanceOf(s.ctx, s.zrc20ABI, s.zrc20Address, s.defaultCaller)
		// require.NoError(t, err)

		// check it was really distributed
	})
}

/*
	Helpers
*/

type testSuite struct {
	ctx            sdk.Context
	contract       *Contract
	contractABI    *abi.ABI
	fungibleKeeper *fungiblekeeper.Keeper
	sdkKeepers     keeper.SDKKeepers
	mockEVM        *vm.EVM
	mockVMContract *vm.Contract
	methodID       abi.Method
	defaultCaller  common.Address
	defaultLocker  common.Address
	zrc20Address   common.Address
	zrc20ABI       *abi.ABI
}

func newTestSuite(t *testing.T) testSuite {
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

	// Create the staking contract.
	contract := NewIStakingContract(
		ctx,
		&sdkKeepers.StakingKeeper,
		*fungibleKeeper,
		sdkKeepers.BankKeeper,
		appCodec,
		gasConfig,
	)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

	accAddress := sdk.AccAddress(ContractAddress.Bytes())
	fungibleKeeper.GetAuthKeeper().SetAccount(ctx, authtypes.NewBaseAccount(accAddress, nil, 0, 0))

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

	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	// Default locker is the bank address.
	locker := common.HexToAddress("0x0000000000000000000000000000000000000067")

	// Set default caller.
	caller := fungibletypes.ModuleAddressEVM
	mockVMContract.CallerAddress = caller
	mockEVM.Origin = caller

	return testSuite{
		ctx,
		contract,
		&abi,
		fungibleKeeper,
		sdkKeepers,
		mockEVM,
		mockVMContract,
		abi.Methods[DistributeMethodName],
		caller,
		locker,
		zrc20Address,
		zrc20ABI,
	}
}

func allowStaking(t *testing.T, ts testSuite, amount *big.Int) {
	resAllowance, err := callEVM(
		t,
		ts.ctx,
		ts.fungibleKeeper,
		ts.zrc20ABI,
		fungibletypes.ModuleAddressEVM,
		ts.zrc20Address,
		"approve",
		[]interface{}{ts.contract.Address(), amount},
	)
	require.NoError(t, err, "error allowing staking to spend ZRC20 tokens")

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
