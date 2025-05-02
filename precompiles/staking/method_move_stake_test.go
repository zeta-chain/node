package staking

import (
	"math/rand"
	"testing"
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	_ "github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_MoveStake(t *testing.T) {
	t.Run("should fail in read only method", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		s.mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}
		s.mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, true)

		// ASSERT
		require.ErrorIs(t, err, precompiletypes.ErrWriteMethod{Method: MoveStakeMethodName})
	})

	t.Run("should fail if destination validator doesn't exist", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded

		s.mockVMContract.CallerAddress = stakerEthAddr

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		s.mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}
		s.mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "redelegation destination validator not found")
	})

	t.Run("should move stake", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		s.mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		argsMoveStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}
		s.mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)
	})

	t.Run("should fail if staker address is invalid arg", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		_, err = s.stakingContract.Stake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&stakeMethodID,
			argsStake,
		)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			42,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// ACT
		_, err = s.stakingContract.MoveStake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			argsMoveStake,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 42 (type int)")
	})

	t.Run("should fail if validator src is invalid arg", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		_, err = s.stakingContract.Stake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&stakeMethodID,
			argsStake,
		)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			stakerEthAddr,
			42,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// ACT
		_, err = s.stakingContract.MoveStake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			argsMoveStake,
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if validator dest is invalid arg", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		_, err = s.stakingContract.Stake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&stakeMethodID,
			argsStake,
		)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			42,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// ACT
		_, err = s.stakingContract.MoveStake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			argsMoveStake,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 42 (type int)")
	})

	t.Run("should fail if amount is invalid arg", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		_, err = s.stakingContract.Stake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&stakeMethodID,
			argsStake,
		)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).Uint64(),
		}

		// ACT
		_, err = s.stakingContract.MoveStake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			argsMoveStake,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 1000000000000000000 (type uint64)")
	})

	t.Run("should fail if wrong args number", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr

		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		_, err = s.stakingContract.Stake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&stakeMethodID,
			argsStake,
		)
		require.NoError(t, err)

		argsMoveStake := []interface{}{stakerEthAddr, validatorSrc.OperatorAddress, validatorDest.OperatorAddress}

		// ACT
		_, err = s.stakingContract.MoveStake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&methodID,
			argsMoveStake,
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid number of arguments")
	})

	t.Run("should fail if caller is not staker", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[MoveStakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validatorSrc, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorSrc)
		require.NoError(t, err)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		validatorDest := sample.Validator(t, r)
		validatorDest.Status = stakingtypes.Bonded
		err = s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validatorDest)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr
		argsStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// stake to validator src
		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]
		s.mockVMContract.Input = packInputArgs(t, stakeMethodID, argsStake...)
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		argsMoveStake := []interface{}{
			stakerEthAddr,
			validatorSrc.OperatorAddress,
			validatorDest.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}
		s.mockVMContract.Input = packInputArgs(t, methodID, argsMoveStake...)

		callerEthAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		s.mockVMContract.CallerAddress = callerEthAddr

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "caller is not staker")
	})
}
