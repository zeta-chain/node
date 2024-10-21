package staking

import (
	"encoding/json"
	"testing"

	"math/big"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/ethermint/x/evm/statedb"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/precompiles/prototype"
	"github.com/zeta-chain/node/testutil/keeper"
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
