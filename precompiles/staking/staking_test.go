package staking

import (
	"encoding/json"
	"testing"

	"math/rand"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func Test_IStakingContract(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	keys, memKeys, tkeys, allKeys := keeper.StoreKeys()
	cdc := keeper.NewCodec()
	sdkKeepers := keeper.NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys, allKeys)

	gasConfig := storetypes.TransientGasConfig()

	t.Run("should create contract and check address and ABI", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

		address := contract.Address()
		require.Equal(t, ContractAddress, address, "contract address should match the precompiled address")

		abi := contract.Abi()
		require.NotNil(t, abi, "contract ABI should not be nil")
	})

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		abi := contract.Abi()

		require.NotNil(t, abi.Methods[StakeMethodName], "stake method should be present in the ABI")
		require.NotNil(t, abi.Methods[UnstakeMethodName], "unstake method should be present in the ABI")
		require.NotNil(
			t,
			abi.Methods[TransferStakeMethodName],
			"transferStake method should be present in the ABI",
		)
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		abi := contract.Abi()
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
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	keys, memKeys, tkeys, allKeys := keeper.StoreKeys()
	cdc := keeper.NewCodec()
	sdkKeepers := keeper.NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys, allKeys)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

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

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{originEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

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

		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).Int64()}

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
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			42,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			42,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
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
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		// stake to validator src
		stakeMethodID := abi.Methods[StakeMethodName]
		_, err = contract.Stake(ctx, stakerAddr, &stakeMethodID, argsStake)
		require.NoError(t, err)

		argsTransferStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Int64(),
		}

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		_, err = contract.TransferStake(ctx, originEthAddr, &methodID, argsTransferStake)
		require.ErrorContains(t, err, "origin is not staker")
	})
}