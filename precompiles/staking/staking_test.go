package staking

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	"math/big"
// 	"math/rand"

// 	tmdb "github.com/cometbft/cometbft-db"
// 	"github.com/cosmos/cosmos-sdk/store"

// 	storetypes "github.com/cosmos/cosmos-sdk/store/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/stretchr/testify/require"
// 	ethermint "github.com/zeta-chain/ethermint/types"
// 	"github.com/zeta-chain/node/cmd/zetacored/config"
// 	"github.com/zeta-chain/node/testutil/keeper"
// 	"github.com/zeta-chain/node/testutil/sample"
// 	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
// )

// func setup(t *testing.T) (sdk.Context, *Contract, abi.ABI, keeper.SDKKeepers) {
// 	var encoding ethermint.EncodingConfig
// 	appCodec := encoding.Codec

// 	cdc := keeper.NewCodec()

// 	db := tmdb.NewMemDB()
// 	stateStore := store.NewCommitMultiStore(db)
// 	sdkKeepers := keeper.NewSDKKeepers(cdc, db, stateStore)
// 	gasConfig := storetypes.TransientGasConfig()
// 	ctx := keeper.NewContext(stateStore)
// 	require.NoError(t, stateStore.LoadLatestVersion())

// 	stakingGenesisState := stakingtypes.DefaultGenesisState()
// 	stakingGenesisState.Params.BondDenom = config.BaseDenom
// 	sdkKeepers.StakingKeeper.InitGenesis(ctx, stakingGenesisState)

// 	contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
// 	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

// 	abi := contract.Abi()
// 	require.NotNil(t, abi, "contract ABI should not be nil")

// 	address := contract.Address()
// 	require.NotNil(t, address, "contract address should not be nil")

// 	return ctx, contract, abi, sdkKeepers
// }

// func Test_IStakingContract(t *testing.T) {
// 	_, contract, abi, _ := setup(t)
// 	gasConfig := storetypes.TransientGasConfig()

// 	t.Run("should check methods are present in ABI", func(t *testing.T) {
// 		require.NotNil(t, abi.Methods[StakeMethodName], "stake method should be present in the ABI")
// 		require.NotNil(t, abi.Methods[UnstakeMethodName], "unstake method should be present in the ABI")
// 		require.NotNil(
// 			t,
// 			abi.Methods[MoveStakeMethodName],
// 			"moveStake method should be present in the ABI",
// 		)

// 		require.NotNil(
// 			t,
// 			abi.Methods[GetAllValidatorsMethodName],
// 			"getAllValidators method should be present in the ABI",
// 		)
// 		require.NotNil(t, abi.Methods[GetSharesMethodName], "getShares method should be present in the ABI")
// 	})

// 	t.Run("should check gas requirements for methods", func(t *testing.T) {
// 		var method [4]byte

// 		t.Run("stake", func(t *testing.T) {
// 			// ACT
// 			stake := contract.RequiredGas(abi.Methods[StakeMethodName].ID)
// 			// ASSERT
// 			copy(method[:], abi.Methods[StakeMethodName].ID[:4])
// 			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
// 			require.Equal(
// 				t,
// 				GasRequiredByMethod[method]+baseCost,
// 				stake,
// 				"stake method should require %d gas, got %d",
// 				GasRequiredByMethod[method]+baseCost,
// 				stake,
// 			)
// 		})

// 		t.Run("unstake", func(t *testing.T) {
// 			// ACT
// 			unstake := contract.RequiredGas(abi.Methods[UnstakeMethodName].ID)
// 			// ASSERT
// 			copy(method[:], abi.Methods[UnstakeMethodName].ID[:4])
// 			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
// 			require.Equal(
// 				t,
// 				GasRequiredByMethod[method]+baseCost,
// 				unstake,
// 				"unstake method should require %d gas, got %d",
// 				GasRequiredByMethod[method]+baseCost,
// 				unstake,
// 			)
// 		})

// 		t.Run("moveStake", func(t *testing.T) {
// 			// ACT
// 			moveStake := contract.RequiredGas(abi.Methods[MoveStakeMethodName].ID)
// 			// ASSERT
// 			copy(method[:], abi.Methods[MoveStakeMethodName].ID[:4])
// 			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
// 			require.Equal(
// 				t,
// 				GasRequiredByMethod[method]+baseCost,
// 				moveStake,
// 				"moveStake method should require %d gas, got %d",
// 				GasRequiredByMethod[method]+baseCost,
// 				moveStake,
// 			)
// 		})

// 		t.Run("getAllValidators", func(t *testing.T) {
// 			// ACT
// 			getAllValidators := contract.RequiredGas(abi.Methods[GetAllValidatorsMethodName].ID)
// 			// ASSERT
// 			copy(method[:], abi.Methods[GetAllValidatorsMethodName].ID[:4])
// 			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
// 			require.Equal(
// 				t,
// 				GasRequiredByMethod[method]+baseCost,
// 				getAllValidators,
// 				"getAllValidators method should require %d gas, got %d",
// 				GasRequiredByMethod[method]+baseCost,
// 				getAllValidators,
// 			)
// 		})

// 		t.Run("getShares", func(t *testing.T) {
// 			// ACT
// 			getShares := contract.RequiredGas(abi.Methods[GetSharesMethodName].ID)
// 			// ASSERT
// 			copy(method[:], abi.Methods[GetSharesMethodName].ID[:4])
// 			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
// 			require.Equal(
// 				t,
// 				GasRequiredByMethod[method]+baseCost,
// 				getShares,
// 				"getShares method should require %d gas, got %d",
// 				GasRequiredByMethod[method]+baseCost,
// 				getShares,
// 			)
// 		})

// 		t.Run("invalid method", func(t *testing.T) {
// 			// ARRANGE
// 			invalidMethodBytes := []byte("invalidMethod")
// 			// ACT
// 			gasInvalidMethod := contract.RequiredGas(invalidMethodBytes)
// 			// ASSERT
// 			require.Equal(
// 				t,
// 				uint64(0),
// 				gasInvalidMethod,
// 				"invalid method should require %d gas, got %d",
// 				uint64(0),
// 				gasInvalidMethod,
// 			)
// 		})
// 	})
// }

// func Test_InvalidMethod(t *testing.T) {
// 	_, _, abi, _ := setup(t)

// 	_, doNotExist := abi.Methods["invalidMethod"]
// 	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
// }

// func Test_InvalidABI(t *testing.T) {
// 	IStakingMetaData.ABI = "invalid json"
// 	defer func() {
// 		if r := recover(); r != nil {
// 			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
// 		}
// 	}()

// 	initABI()
// }

// func Test_Stake(t *testing.T) {
// 	ctx, contract, abi, sdkKeepers := setup(t)
// 	methodID := abi.Methods[StakeMethodName]

// 	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should stake", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.NoError(t, err)
// 	})

// 	t.Run("should fail if origin is not staker", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

// 		args := []interface{}{originEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.ErrorContains(t, err, "origin is not staker address")
// 	})

// 	t.Run("should fail if staking fails", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		// staker without funds
// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, int64(42)}

// 		// ACT
// 		_, err := contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if wrong args amount", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if staker is not eth addr", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if validator is not valid string", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if amount is not int64", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})
// }

// func Test_Unstake(t *testing.T) {
// 	ctx, contract, abi, sdkKeepers := setup(t)
// 	methodID := abi.Methods[UnstakeMethodName]

// 	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should unstake", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// stake first
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, args)
// 		require.NoError(t, err)

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.NoError(t, err)
// 	})

// 	t.Run("should fail if origin is not staker", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// stake first
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, args)
// 		require.NoError(t, err)

// 		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

// 		// ACT
// 		_, err = contract.Unstake(ctx, originEthAddr, &methodID, args)

// 		// ASSERT
// 		require.ErrorContains(t, err, "origin is not staker address")
// 	})

// 	t.Run("should fail if no previous staking", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if wrong args amount", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress}

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if staker is not eth addr", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if validator is not valid string", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if amount is not int64", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

// 		// ACT
// 		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})
// }

// func Test_MoveStake(t *testing.T) {
// 	ctx, contract, abi, sdkKeepers := setup(t)
// 	methodID := abi.Methods[MoveStakeMethodName]

// 	t.Run("should fail if validator dest doesn't exist", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			validatorDest.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// ACT
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should move stake", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			validatorDest.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// ACT
// 		// move stake to validator dest
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.NoError(t, err)
// 	})

// 	t.Run("should fail if staker is invalid arg", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			42,
// 			validatorSrc.OperatorAddress,
// 			validatorDest.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// ACT
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if validator src is invalid arg", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			stakerEthAddr,
// 			42,
// 			validatorDest.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// ACT
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if validator dest is invalid arg", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			42,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// ACT
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if amount is invalid arg", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			validatorDest.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).Uint64(),
// 		}

// 		// ACT
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if wrong args amount", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{stakerEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress}

// 		// ACT
// 		_, err = contract.MoveStake(ctx, stakerAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if origin is not staker", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validatorSrc := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
// 		validatorDest := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		argsStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		// stake to validator src
// 		stakeMethodID := abi.Methods[StakeMethodName]
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
// 		require.NoError(t, err)

// 		argsMoveStake := []interface{}{
// 			stakerEthAddr,
// 			validatorSrc.OperatorAddress,
// 			validatorDest.OperatorAddress,
// 			coins.AmountOf(config.BaseDenom).BigInt(),
// 		}

// 		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

// 		// ACT
// 		_, err = contract.MoveStake(ctx, originEthAddr, &methodID, argsMoveStake)

// 		// ASSERT
// 		require.ErrorContains(t, err, "origin is not staker")
// 	})
// }

// func Test_GetAllValidators(t *testing.T) {
// 	ctx, contract, abi, sdkKeepers := setup(t)
// 	methodID := abi.Methods[GetAllValidatorsMethodName]

// 	t.Run("should return empty array if validators not set", func(t *testing.T) {
// 		// ACT
// 		validators, err := contract.GetAllValidators(ctx, &methodID)

// 		// ASSERT
// 		require.NoError(t, err)

// 		res, err := methodID.Outputs.Unpack(validators)
// 		require.NoError(t, err)

// 		require.Empty(t, res[0])
// 	})

// 	t.Run("should return validators if set", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		// ACT
// 		validators, err := contract.GetAllValidators(ctx, &methodID)

// 		// ASSERT
// 		require.NoError(t, err)

// 		res, err := methodID.Outputs.Unpack(validators)
// 		require.NoError(t, err)

// 		require.NotEmpty(t, res[0])
// 	})
// }

// func Test_GetShares(t *testing.T) {
// 	ctx, contract, abi, sdkKeepers := setup(t)
// 	methodID := abi.Methods[GetSharesMethodName]

// 	t.Run("should return stakes", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		coins := sample.Coins()
// 		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
// 		require.NoError(t, err)
// 		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
// 		require.NoError(t, err)

// 		stakerAddr := common.BytesToAddress(staker.Bytes())

// 		stakeArgs := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

// 		stakeMethodID := abi.Methods[StakeMethodName]

// 		// ACT
// 		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, stakeArgs)
// 		require.NoError(t, err)

// 		// ASSERT
// 		args := []interface{}{stakerEthAddr, validator.OperatorAddress}
// 		stakes, err := contract.GetShares(ctx, &methodID, args)
// 		require.NoError(t, err)

// 		res, err := methodID.Outputs.Unpack(stakes)
// 		require.NoError(t, err)
// 		require.Equal(
// 			t,
// 			fmt.Sprintf("%d000000000000000000", coins.AmountOf(config.BaseDenom).BigInt().Int64()),
// 			res[0].(*big.Int).String(),
// 		)
// 	})

// 	t.Run("should fail if wrong args amount", func(t *testing.T) {
// 		// ARRANGE
// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		args := []interface{}{stakerEthAddr}

// 		// ACT
// 		_, err := contract.GetShares(ctx, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if invalid staker arg", func(t *testing.T) {
// 		// ARRANGE
// 		r := rand.New(rand.NewSource(42))
// 		validator := sample.Validator(t, r)
// 		args := []interface{}{42, validator.OperatorAddress}

// 		// ACT
// 		_, err := contract.GetShares(ctx, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})

// 	t.Run("should fail if invalid val address", func(t *testing.T) {
// 		// ARRANGE
// 		staker := sample.Bech32AccAddress()
// 		stakerEthAddr := common.BytesToAddress(staker.Bytes())
// 		args := []interface{}{stakerEthAddr, staker.String()}

// 		// ACT
// 		_, err := contract.GetShares(ctx, &methodID, args)

// 		// ASSERT
// 		require.Error(t, err)
// 	})
// }
