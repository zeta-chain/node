package staking

import (
	"testing"

	_ "github.com/ethereum/go-ethereum/core/vm"
)

func Test_MoveStake(t *testing.T) {
	//t.Run("should fail in read only method", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//	mockVMContract.CallerAddress = stakerAddr
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//	mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)
	//
	//	// ACT
	//	_, err = contract.Run(mockEVM, mockVMContract, true)
	//
	//	// ASSERT
	//	require.ErrorIs(t, err, precompiletypes.ErrWriteMethod{Method: MoveStakeMethodName})
	//})
	//
	//t.Run("should fail if validator dest doesn't exist", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//	mockVMContract.CallerAddress = stakerAddr
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//	mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)
	//
	//	// ACT
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//})

	//t.Run("should move stake", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	err := sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	require.NoError(t, err)
	//	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorDest := sample.Validator(t, r)
	//	err = sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//	require.NoError(t, err)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err = sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//	mockVMContract.CallerAddress = stakerAddr
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//	require.NoError(t, err)
	//
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//	mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)
	//
	//	// ACT
	//	// move stake to validator dest
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//
	//	// ASSERT
	//	require.NoError(t, err)
	//})

	//t.Run("should fail if staker is invalid arg", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, argsStake)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		42,
	//		validatorSrc.OperatorAddress,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// ACT
	//	_, err = contract.MoveStake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, argsMoveStake)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//})
	//
	//t.Run("should fail if validator src is invalid arg", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, argsStake)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		42,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// ACT
	//	_, err = contract.MoveStake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, argsMoveStake)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//})
	//
	//t.Run("should fail if validator dest is invalid arg", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, argsStake)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		42,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// ACT
	//	_, err = contract.MoveStake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, argsMoveStake)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//})
	//
	//t.Run("should fail if amount is invalid arg", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, argsStake)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).Uint64(),
	//	}
	//
	//	// ACT
	//	_, err = contract.MoveStake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, argsMoveStake)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//})
	//
	//t.Run("should fail if wrong args amount", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, argsStake)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{stakerEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress}
	//
	//	// ACT
	//	_, err = contract.MoveStake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, argsMoveStake)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//})
	//
	//t.Run("should fail if caller is not staker", func(t *testing.T) {
	//	// ARRANGE
	//	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	//	methodID := abi.Methods[MoveStakeMethodName]
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	validatorSrc := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorSrc)
	//	validatorDest := sample.Validator(t, r)
	//	sdkKeepers.StakingKeeper.SetValidator(ctx, validatorDest)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//	mockVMContract.CallerAddress = stakerAddr
	//	argsStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//
	//	// stake to validator src
	//	stakeMethodID := abi.Methods[StakeMethodName]
	//	mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//	require.NoError(t, err)
	//
	//	argsMoveStake := []interface{}{
	//		stakerEthAddr,
	//		validatorSrc.OperatorAddress,
	//		validatorDest.OperatorAddress,
	//		coins.AmountOf(config.BaseDenom).BigInt(),
	//	}
	//	mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)
	//
	//	callerEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
	//	mockVMContract.CallerAddress = callerEthAddr
	//
	//	// ACT
	//	_, err = contract.Run(mockEVM, mockVMContract, false)
	//
	//	// ASSERT
	//	require.ErrorContains(t, err, "caller is not staker")
	//})
}
