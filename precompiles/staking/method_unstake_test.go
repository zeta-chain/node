package staking

import (
	"testing"
)

func Test_Unstake(t *testing.T) {
	// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
	//t.Run("should fail with error disabled", func(t *testing.T) {
	//	// ARRANGE
	//	s := newTestSuite(t)
	//	methodID := s.stkContractABI.Methods[UnstakeMethodName]
	//	r := rand.New(rand.NewSource(42))
	//	validator := sample.Validator(t, r)
	//
	//	staker := sample.Bech32AccAddress()
	//	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	//	coins := sample.Coins()
	//	err := s.sdkKeepers.BankKeeper.MintCoins(s.ctx, fungibletypes.ModuleName, sample.Coins())
	//	require.NoError(t, err)
	//	err = s.sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, fungibletypes.ModuleName, staker, coins)
	//	require.NoError(t, err)
	//
	//	stakerAddr := common.BytesToAddress(staker.Bytes())
	//	s.mockVMContract.CallerAddress = stakerAddr
	//
	//	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	//	s.mockVMContract.Input = packInputArgs(t, methodID, args...)
	//
	//	// ACT
	//	_, err = s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
	//
	//	// ASSERT
	//	require.Error(t, err)
	//	require.ErrorIs(t, err, precompiletypes.ErrDisabledMethod{
	//		Method: UnstakeMethodName,
	//	})
	//})

	// Tests have been disabled since
	// https://github.com/zeta-chain/node/pull/3019

	// t.Run("should fail in read only method", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)

	// 	// ACT
	// 	_, err = contract.Run(mockEVM, mockVMContract, true)

	// 	// ASSERT
	// 	require.ErrorIs(t, err, ptypes.ErrWriteMethod{Method: UnstakeMethodName})
	// })

	// t.Run("should fail if validator doesn't exist", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)

	// 	// ACT
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should unstake", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

	// 	// stake first
	// 	stakeMethodID := abi.Methods[StakeMethodName]
	// 	mockVMContract.Input = packInputArgs(t, stakeMethodID, args...)
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)
	// 	require.NoError(t, err)

	// 	// ACT
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.NoError(t, err)
	// })

	// t.Run("should fail if caller is not staker", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	// 	// stake first
	// 	stakeMethodID := abi.Methods[StakeMethodName]
	// 	mockVMContract.Input = packInputArgs(t, stakeMethodID, args...)
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)
	// 	require.NoError(t, err)

	// 	callerEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
	// 	mockVMContract.CallerAddress = callerEthAddr
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)

	// 	// ACT
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.ErrorContains(t, err, "caller is not staker address")
	// })

	// t.Run("should fail if no previous staking", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)

	// 	// ACT
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if wrong args amount", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress}

	// 	// ACT
	// 	_, err = contract.Unstake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if staker is not eth addr", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())

	// 	args := []interface{}{staker, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

	// 	// ACT
	// 	_, err = contract.Unstake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if validator is not valid string", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())

	// 	args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

	// 	// ACT
	// 	_, err = contract.Unstake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if amount is not int64", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[UnstakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

	// 	// ACT
	// 	_, err = contract.Unstake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })
}
