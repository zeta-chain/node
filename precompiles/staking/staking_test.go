package staking

import (
	"encoding/json"
	"fmt"
	"testing"

	"math/big"
	"math/rand"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func setup(t *testing.T) (sdk.Context, *Contract, abi.ABI, keeper.SDKKeepers) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec

	cdc := keeper.NewCodec()

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	sdkKeepers := keeper.NewSDKKeepers(cdc, db, stateStore)
	gasConfig := storetypes.TransientGasConfig()
	ctx := keeper.NewContext(stateStore)
	require.NoError(t, stateStore.LoadLatestVersion())

	stakingGenesisState := stakingtypes.DefaultGenesisState()
	stakingGenesisState.Params.BondDenom = config.BaseDenom
	sdkKeepers.StakingKeeper.InitGenesis(ctx, stakingGenesisState)

	contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	return ctx, contract, abi, sdkKeepers
}

func Test_IStakingContract(t *testing.T) {
	_, contract, abi, _ := setup(t)
	gasConfig := storetypes.TransientGasConfig()

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		require.NotNil(t, abi.Methods[StakeMethodName], "stake method should be present in the ABI")
		require.NotNil(t, abi.Methods[UnstakeMethodName], "unstake method should be present in the ABI")
		require.NotNil(
			t,
			abi.Methods[TransferStakeMethodName],
			"transferStake method should be present in the ABI",
		)

		require.NotNil(
			t,
			abi.Methods[GetAllValidatorsMethodName],
			"getAllValidators method should be present in the ABI",
		)
		require.NotNil(t, abi.Methods[GetStakesMethodName], "getStakes method should be present in the ABI")
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		var method [4]byte

		t.Run("stake", func(t *testing.T) {
			stake := contract.RequiredGas(abi.Methods[StakeMethodName].ID)
			copy(method[:], abi.Methods[StakeMethodName].ID[:4])
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
			unstake := contract.RequiredGas(abi.Methods[UnstakeMethodName].ID)
			copy(method[:], abi.Methods[UnstakeMethodName].ID[:4])
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

		t.Run("transferStake", func(t *testing.T) {
			transferStake := contract.RequiredGas(abi.Methods[TransferStakeMethodName].ID)
			copy(method[:], abi.Methods[TransferStakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				transferStake,
				"transferStake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				transferStake,
			)
		})

		t.Run("getAllValidators", func(t *testing.T) {
			getAllValidators := contract.RequiredGas(abi.Methods[GetAllValidatorsMethodName].ID)
			copy(method[:], abi.Methods[GetAllValidatorsMethodName].ID[:4])
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

		t.Run("getStakes", func(t *testing.T) {
			getStakes := contract.RequiredGas(abi.Methods[GetStakesMethodName].ID)
			copy(method[:], abi.Methods[GetStakesMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				getStakes,
				"getStakes method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				getStakes,
			)
		})

		t.Run("invalid method", func(t *testing.T) {
			invalidMethodBytes := []byte("invalidMethod")
			gasInvalidMethod := contract.RequiredGas(invalidMethodBytes)
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
	_, _, abi, _ := setup(t)

	_, doNotExist := abi.Methods["invalidMethod"]
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

func Test_Stake(t *testing.T) {
	ctx, contract, abi, sdkKeepers := setup(t)
	methodID := abi.Methods[StakeMethodName]

	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should stake", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.NoError(t, err)
	})

	t.Run("should fail if origin is not staker", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

		args := []interface{}{originEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.ErrorContains(t, err, "origin is not staker address")
	})

	t.Run("should fail if staking fails", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		// staker without funds
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, int64(42)}

		_, err := contract.Stake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if staker is not eth addr", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if validator is not valid string", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if amount is not int64", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		_, err = contract.Stake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})
}

func Test_Unstake(t *testing.T) {
	ctx, contract, abi, sdkKeepers := setup(t)
	methodID := abi.Methods[UnstakeMethodName]

	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should unstake", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		// stake first
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, args)
		require.NoError(t, err)

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.NoError(t, err)
	})

	t.Run("should fail if origin is not staker", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		// stake first
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, args)
		require.NoError(t, err)

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

		_, err = contract.Unstake(ctx, originEthAddr, &methodID, args)
		require.ErrorContains(t, err, "origin is not staker address")
	})

	t.Run("should fail if no previous staking", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress}

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if staker is not eth addr", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if validator is not valid string", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if amount is not int64", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		_, err = contract.Unstake(ctx, stakerAddr, &methodID, args)
		require.Error(t, err)
	})
}

func Test_TransferStake(t *testing.T) {
	ctx, contract, abi, sdkKeepers := setup(t)
	methodID := abi.Methods[TransferStakeMethodName]

	t.Run("should fail if validator dest doesn't exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.Error(t, err)
	})

	t.Run("should transfer stake", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// transfer stake to validator dest
		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.NoError(t, err)
	})

	t.Run("should fail if staker is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			42,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.Error(t, err)
	})

	t.Run("should fail if validator src is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			42,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.Error(t, err)
	})

	t.Run("should fail if validator dest is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.NoError(t, err)
	})

	t.Run("should fail if amount is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Uint64(),
		}

		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{stakerEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress}

		_, err = contract.TransferStake(ctx, stakerAddr, &methodID, argsTransferStake)
		require.Error(t, err)
	})

	t.Run("should fail if origin is not staker", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		_, err = contract.TransferStake(ctx, originEthAddr, &methodID, argsTransferStake)
		require.ErrorContains(t, err, "origin is not staker")
	})
}

func Test_GetAllValidators(t *testing.T) {
	ctx, contract, abi, sdkKeepers := setup(t)
	methodID := abi.Methods[GetAllValidatorsMethodName]

	t.Run("should return empty array if validators not set", func(t *testing.T) {
		validators, err := contract.GetAllValidators(ctx, &methodID)
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.Empty(t, res[0])
	})

	t.Run("should return validators if set", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		validators, err := contract.GetAllValidators(ctx, &methodID)
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.NotEmpty(t, res[0])
	})
}

func Test_GetStakes(t *testing.T) {
	ctx, contract, abi, sdkKeepers := setup(t)
	methodID := abi.Methods[GetStakesMethodName]

	t.Run("should return stakes", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		stakeArgs := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, stakeArgs)
		require.NoError(t, err)

		args := []interface{}{stakerEthAddr, validator.OperatorAddress}
		stakes, err := contract.GetStake(ctx, &methodID, args)
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(stakes)
		require.NoError(t, err)
		require.Equal(
			t,
			fmt.Sprintf("%d000000000000000000", coins.AmountOf(config.BaseDenom).BigInt().Int64()),
			res[0].(*big.Int).String(),
		)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr}
		_, err := contract.GetStake(ctx, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if invalid staker arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)

		args := []interface{}{42, validator.OperatorAddress}
		_, err := contract.GetStake(ctx, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if invalid val address", func(t *testing.T) {
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())

		args := []interface{}{stakerEthAddr, staker}
		_, err := contract.GetStake(ctx, &methodID, args)
		require.Error(t, err)
	})
}
