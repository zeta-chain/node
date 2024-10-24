package staking

import (
	"encoding/json"
	"testing"

	"math/big"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	evmkeeper "github.com/zeta-chain/ethermint/x/evm/keeper"
	"github.com/zeta-chain/ethermint/x/evm/statedb"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/contracts/erc1967proxy"
	"github.com/zeta-chain/node/precompiles/prototype"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"
)

func Test_IStakingContract(t *testing.T) {
	s := newTestSuite(t)
	gasConfig := storetypes.TransientGasConfig()

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		require.NotNil(t, s.contractABI.Methods[StakeMethodName], "stake method should be present in the ABI")
		require.NotNil(t, s.contractABI.Methods[UnstakeMethodName], "unstake method should be present in the ABI")
		require.NotNil(
			t,
			s.contractABI.Methods[MoveStakeMethodName],
			"moveStake method should be present in the ABI",
		)

		require.NotNil(
			t,
			s.contractABI.Methods[GetAllValidatorsMethodName],
			"getAllValidators method should be present in the ABI",
		)
		require.NotNil(t, s.contractABI.Methods[GetSharesMethodName], "getShares method should be present in the ABI")
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		var method [4]byte

		t.Run("stake", func(t *testing.T) {
			// ACT
			stake := s.contract.RequiredGas(s.contractABI.Methods[StakeMethodName].ID)
			// ASSERT
			copy(method[:], s.contractABI.Methods[StakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				stake,
				"stake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				stake,
			)
		})

		t.Run("unstake", func(t *testing.T) {
			// ACT
			unstake := s.contract.RequiredGas(s.contractABI.Methods[UnstakeMethodName].ID)
			// ASSERT
			copy(method[:], s.contractABI.Methods[UnstakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				unstake,
				"unstake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				unstake,
			)
		})

		t.Run("moveStake", func(t *testing.T) {
			// ACT
			moveStake := s.contract.RequiredGas(s.contractABI.Methods[MoveStakeMethodName].ID)
			// ASSERT
			copy(method[:], s.contractABI.Methods[MoveStakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				moveStake,
				"moveStake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				moveStake,
			)
		})

		t.Run("getAllValidators", func(t *testing.T) {
			// ACT
			getAllValidators := s.contract.RequiredGas(s.contractABI.Methods[GetAllValidatorsMethodName].ID)
			// ASSERT
			copy(method[:], s.contractABI.Methods[GetAllValidatorsMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				getAllValidators,
				"getAllValidators method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				getAllValidators,
			)
		})

		t.Run("getShares", func(t *testing.T) {
			// ACT
			getShares := s.contract.RequiredGas(s.contractABI.Methods[GetSharesMethodName].ID)
			// ASSERT
			copy(method[:], s.contractABI.Methods[GetSharesMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				getShares,
				"getShares method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				getShares,
			)
		})

		t.Run("invalid method", func(t *testing.T) {
			// ARRANGE
			invalidMethodBytes := []byte("invalidMethod")
			// ACT
			gasInvalidMethod := s.contract.RequiredGas(invalidMethodBytes)
			// ASSERT
			require.Equal(
				t,
				uint64(0),
				gasInvalidMethod,
				"invalid method should require %d gas, got %d",
				uint64(0),
				gasInvalidMethod,
			)
		})
	})
}

func Test_InvalidMethod(t *testing.T) {
	s := newTestSuite(t)

	_, doNotExist := s.contractABI.Methods["invalidMethod"]
	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
}

func Test_InvalidABI(t *testing.T) {
	IStakingMetaData.ABI = "invalid json"
	defer func() {
		if r := recover(); r != nil {
			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
		}
	}()

	initABI()
}

func Test_RunInvalidMethod(t *testing.T) {
	// ARRANGE
	s := newTestSuite(t)

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	gasConfig := storetypes.TransientGasConfig()

	prototype := prototype.NewIPrototypeContract(s.fungibleKeeper, appCodec, gasConfig)

	prototypeAbi := prototype.Abi()
	methodID := prototypeAbi.Methods["bech32ToHexAddr"]
	args := []interface{}{"123"}
	s.mockVMContract.Input = packInputArgs(t, methodID, args...)

	// ACT
	_, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

	// ASSERT
	require.Error(t, err)
}

func setup(t *testing.T) (sdk.Context, *Contract, abi.ABI, keeper.SDKKeepers, *vm.EVM, *vm.Contract) {
	// Initialize state.
	// Get sdk keepers initialized with this state and the context.
	cdc := keeper.NewCodec()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	keys, memKeys, tkeys, allKeys := keeper.StoreKeys()

	sdkKeepers := keeper.NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys, allKeys)

	for _, key := range keys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}
	for _, key := range tkeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeTransient, nil)
	}
	for _, key := range memKeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeMemory, nil)
	}

	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := keeper.NewContext(stateStore)

	// Intiliaze codecs and gas config.
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	gasConfig := storetypes.TransientGasConfig()

	stakingGenesisState := stakingtypes.DefaultGenesisState()
	stakingGenesisState.Params.BondDenom = config.BaseDenom
	sdkKeepers.StakingKeeper.InitGenesis(ctx, stakingGenesisState)

	// Get the fungible keeper.
	fungibleKeeper, _, _, _ := keeper.FungibleKeeper(t)

	accAddress := sdk.AccAddress(ContractAddress.Bytes())
	fungibleKeeper.GetAuthKeeper().SetAccount(ctx, authtypes.NewBaseAccount(accAddress, nil, 0, 0))

	// Initialize staking contract.
	contract := NewIStakingContract(
		ctx,
		&sdkKeepers.StakingKeeper,
		*fungibleKeeper,
		sdkKeepers.BankKeeper,
		appCodec,
		gasConfig,
	)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

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

	return ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract
}

/*
	Complete Test Suite
	TODO: Migrate all staking tests to this suite.
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

func packInputArgs(t *testing.T, methodID abi.Method, args ...interface{}) []byte {
	input, err := methodID.Inputs.Pack(args...)
	require.NoError(t, err)
	return append(methodID.ID, input...)
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

type contractRef struct {
	address common.Address
}

func (c contractRef) Address() common.Address {
	return c.address
}
