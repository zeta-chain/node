package staking

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_Unstake(t *testing.T) {
	t.Run("should fail in read only method", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, true)

		// ASSERT
		require.ErrorIs(t, err, precompiletypes.ErrWriteMethod{Method: UnstakeMethodName})
	})

	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)

		s.mockVMContract.CallerAddress = stakerEthAddr
		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err := s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "validator does not exist")
	})

	t.Run("should unstake", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		// stake first
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		s.mockVMContract.Input = packInputArgs(t, stakeMethodID, args...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		// ACT
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)
	})

	t.Run("should fail if caller is not staker", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
		// stake first
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		s.mockVMContract.Input = packInputArgs(t, stakeMethodID, args...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		callerEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		s.mockVMContract.CallerAddress = callerEthAddr
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.ErrorContains(t, err, "caller is not staker address")
	})

	t.Run("should fail if no previous staking", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, _ := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress}

		// ACT
		_, err = s.stakingContract.Unstake(s.ctx, s.mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid number of arguments")
	})

	t.Run("should fail if staker is not eth addr", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddress, validator, coins := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{
			sample.Bech32AccAddress(),
			validator.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// ACT
		_, err := s.stakingContract.Unstake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddress},
			&methodID,
			args,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("invalid argument: got %v (type types.AccAddress)", args[0]))
	})

	t.Run("should fail if validator is not valid string", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

		// ACT
		_, err = s.stakingContract.Unstake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			args,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 42 (type int)")
	})

	t.Run("should fail if amount is not int64", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[UnstakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		// ACT
		_, err := s.stakingContract.Unstake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			args,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 1000000000000000000 (type uint64)")
	})
}
