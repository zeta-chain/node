package staking

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_Stake(t *testing.T) {
	t.Run("should fail in read only mode", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
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
		require.ErrorIs(t, err, precompiletypes.ErrWriteMethod{Method: StakeMethodName})
	})

	t.Run("should fail if validator doesn't exist", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
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

	t.Run("should stake", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
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
		require.NoError(t, err)
	})

	t.Run("should fail if no input args", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]

		stakerAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		s.mockVMContract.CallerAddress = stakerAddr
		s.mockVMContract.Input = methodID.ID

		// ACT
		_, err := s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "attempting to unmarshal an empty string while arguments are expected")
	})

	t.Run("should fail if caller is not staker", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddress, validator, coins := s.setupStakerDefaultAmount(t, r)

		s.mockVMContract.CallerAddress = stakerEthAddress
		nonStakerAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())
		args := []interface{}{nonStakerAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err := s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "caller is not staker address")
	})

	t.Run("should fail if staking fails because of trying to stake more than available balance", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, math.ZeroInt()))
		stakerEthAddr, validator := s.setupStaker(t, r, coins)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		s.mockVMContract.CallerAddress = stakerEthAddr

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, math.OneInt().BigInt()}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)

		// ACT
		_, err = s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "insufficient funds")
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, _ := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(sample.Bech32AccAddress().Bytes())

		args := []interface{}{stakerEthAddr, validator.OperatorAddress}

		// ACT
		_, err = s.stakingContract.Stake(s.ctx, s.mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid number of arguments")
	})

	t.Run("should fail if staker is not eth addr", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddress, validator, coins := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{
			sample.Bech32AccAddress(),
			validator.OperatorAddress,
			coins.AmountOf(config.BaseDenom).BigInt(),
		}

		// ACT
		_, err := s.stakingContract.Stake(
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
		methodID := s.stakingContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		args := []interface{}{stakerEthAddr, 42, coins.AmountOf(config.BaseDenom).BigInt()}

		// ACT
		_, err = s.stakingContract.Stake(s.ctx, s.mockEVM, &vm.Contract{CallerAddress: stakerEthAddr}, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 42 (type int)")

	})

	t.Run("should fail if amount is invalid", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[StakeMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).Uint64()}

		// ACT
		_, err := s.stakingContract.Stake(s.ctx, s.mockEVM, &vm.Contract{CallerAddress: stakerEthAddr}, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 1000000000000000000 (type uint64)")
	})
}
