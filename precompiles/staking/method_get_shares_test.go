package staking

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetShares(t *testing.T) {
	t.Run("should return stakes", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, validator, coins := s.setupStakerDefaultAmount(t, r)
		err := s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)
		require.NoError(t, err)

		stakeArgs := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		stakeMethodID := s.stakingContractABI.Methods[StakeMethodName]

		// ACT
		_, err = s.stakingContract.Stake(
			s.ctx,
			s.mockEVM,
			&vm.Contract{CallerAddress: stakerEthAddr},
			&stakeMethodID,
			stakeArgs,
		)
		require.NoError(t, err)

		// ASSERT
		args := []interface{}{stakerEthAddr, validator.OperatorAddress}
		s.mockVMContract.Input = packInputArgs(t, methodID, args...)
		stakes, err := s.stakingContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(stakes)
		require.NoError(t, err)
		require.Equal(
			t,
			fmt.Sprintf("%d000000000000000000", coins.AmountOf(config.BaseDenom).BigInt().Int64()),
			res[0].(*big.Int).String(),
		)
	})

	t.Run("should fail if wrong args number", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, _, _ := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{stakerEthAddr}
		// ACT
		_, err := s.stakingContract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid number of arguments")
	})

	t.Run("should fail if invalid staker arg", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		_, validator, _ := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{42, validator.OperatorAddress}

		// ACT
		_, err := s.stakingContract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 42 (type int)")
	})

	t.Run("should fail if invalid val address", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, _, _ := s.setupStakerDefaultAmount(t, r)

		// Set AccAddress instead of ValAddress
		args := []interface{}{stakerEthAddr, sample.Bech32AccAddress().String()}

		// ACT
		_, err := s.stakingContract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid Bech32 prefix; expected zetavaloper, got zeta")
	})

	t.Run("should fail if invalid val address format", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stakingContractABI.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stakerEthAddr, _, _ := s.setupStakerDefaultAmount(t, r)

		args := []interface{}{stakerEthAddr, 42}

		// ACT
		_, err := s.stakingContract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid argument: got 42 (type int)")
	})
}
