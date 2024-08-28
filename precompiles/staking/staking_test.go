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

		require.NotNil(t, abi.Methods[DelegateMethodName], "delegate method should be present in the ABI")
		require.NotNil(t, abi.Methods[UndelegateMethodName], "undelegate method should be present in the ABI")
		require.NotNil(
			t,
			abi.Methods[RedelegateMethodName],
			"redelegate method should be present in the ABI",
		)
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		abi := contract.Abi()
		var method [4]byte

		t.Run("delegate", func(t *testing.T) {
			delegate := contract.RequiredGas(abi.Methods[DelegateMethodName].ID)
			copy(method[:], abi.Methods[DelegateMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				delegate,
				"delegate method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				delegate,
			)
		})

		t.Run("undelegate", func(t *testing.T) {
			undelegate := contract.RequiredGas(abi.Methods[UndelegateMethodName].ID)
			copy(method[:], abi.Methods[UndelegateMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				undelegate,
				"undelegate method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				undelegate,
			)
		})

		t.Run("redelegate", func(t *testing.T) {
			redelegate := contract.RequiredGas(abi.Methods[RedelegateMethodName].ID)
			copy(method[:], abi.Methods[RedelegateMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				redelegate,
				"redelegate method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				redelegate,
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

func Test_Delegate(t *testing.T) {
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

	methodID := abi.Methods[DelegateMethodName]

	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should delegate", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.NoError(t, err)
	})

	t.Run("should fail if origin is not delegator", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

		args := []interface{}{originEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.ErrorContains(t, err, "origin is not delegator address")
	})

	t.Run("should fail if delegation fails", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		// delegator without funds
		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, int64(42)}

		_, err := contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if delegator is not eth addr", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegator, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if validator is not valid string", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, 42, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if amount is not int64", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		_, err = contract.Delegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})
}

func Test_Undelegate(t *testing.T) {
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

	methodID := abi.Methods[UndelegateMethodName]

	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should undelegate", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate first
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, args)
		require.NoError(t, err)

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.NoError(t, err)
	})

	t.Run("should fail if origin is not delegator", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate first
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, args)
		require.NoError(t, err)

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

		_, err = contract.Undelegate(ctx, originEthAddr, &methodID, args)
		require.ErrorContains(t, err, "origin is not delegator address")
	})

	t.Run("should fail if no previous delegation", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress}

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if delegator is not eth addr", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegator, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if validator is not valid string", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, 42, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})

	t.Run("should fail if amount is not int64", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		args := []interface{}{delegatorEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		_, err = contract.Undelegate(ctx, delegatorAddr, &methodID, args)
		require.Error(t, err)
	})
}

func Test_Redelegate(t *testing.T) {
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

	methodID := abi.Methods[RedelegateMethodName]

	t.Run("should fail if validator dest doesn't exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.Error(t, err)
	})

	t.Run("should redelegate", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// redelegate to validator dest
		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.NoError(t, err)
	})

	t.Run("should fail if delegator is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{42, validatorSrc.OperatorAddress, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.Error(t, err)
	})

	t.Run("should fail if validator src is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, 42, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.Error(t, err)
	})

	t.Run("should fail if validator dest is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.NoError(t, err)
	})

	t.Run("should fail if amount is invalid arg", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress}

		_, err = contract.Redelegate(ctx, delegatorAddr, &methodID, argsRedelegate)
		require.Error(t, err)
	})

	t.Run("should fail if origin is not delegator", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		validatorSrc := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
		validatorDest := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)

		delegator := sample.Bech32AccAddress()
		delegatorEthAddr := common.BytesToAddress(delegator.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, delegator, coins)
		require.NoError(t, err)

		delegatorAddr := common.BytesToAddress(delegator.Bytes())

		argsDelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		// delegate to validator src
		delegateMethodID := abi.Methods[DelegateMethodName]
		_, err = contract.Delegate(ctx, delegatorAddr, &delegateMethodID, argsDelegate)
		require.NoError(t, err)

		argsRedelegate := []interface{}{delegatorEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress, coins.AmountOf(config.BaseDenom).Int64()}

		originEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		_, err = contract.Redelegate(ctx, originEthAddr, &methodID, argsRedelegate)
		require.ErrorContains(t, err, "origin is not delegator")
	})
}
