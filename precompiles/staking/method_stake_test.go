package staking

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/sample"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func Test_Stake(t *testing.T) {
	// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
	t.Run("should fail with error disabled", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stkContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := s.sdkKeepers.BankKeeper.MintCoins(s.ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = s.sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())
		s.mockVMContract.CallerAddress = stakerAddr
		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err = s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, precompiletypes.ErrDisabledMethod{
			Method: StakeMethodName,
		})
	})

	// t.Run("should fail in read only mode", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	require.ErrorIs(t, err, ptypes.ErrWriteMethod{Method: StakeMethodName})
	// })

	// t.Run("should fail if validator doesn't exist", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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

	// t.Run("should stake", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	require.NoError(t, err)
	// })

	// t.Run("should fail if no input args", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	mockVMContract.CallerAddress = stakerAddr
	// 	mockVMContract.Input = methodID.ID

	// 	// ACT
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if caller is not staker", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	nonStakerAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
	// 	args := []interface{}{nonStakerAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)

	// 	// ACT
	// 	_, err = contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.ErrorContains(t, err, "caller is not staker address")
	// })

	// t.Run("should fail if staking fails", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

	// 	// staker without funds
	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())
	// 	mockVMContract.CallerAddress = stakerAddr

	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
	// 	mockVMContract.Input = packInputArgs(t, methodID, args...)

	// 	// ACT
	// 	_, err := contract.Run(mockEVM, mockVMContract, false)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if wrong args amount", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if staker is not eth addr", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if validator is not valid string", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })

	// t.Run("should fail if amount is not int64", func(t *testing.T) {
	// 	// ARRANGE
	// 	ctx, contract, abi, sdkKeepers, mockEVM, _ := setup(t)
	// 	methodID := abi.Methods[StakeMethodName]
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
	// 	_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

	// 	// ASSERT
	// 	require.Error(t, err)
	// })
}
