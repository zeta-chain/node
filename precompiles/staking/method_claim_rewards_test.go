package staking

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_ClaimRewards(t *testing.T) {
	t.Run("should return an error when passing empty delegator", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))

		/* ACT */
		// Call claimRewardsMethod.
		claimRewardsMethod := s.stkContractABI.Methods[ClaimRewardsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			claimRewardsMethod,
			[]interface{}{common.Address{}, validator.OperatorAddress}...,
		)

		_, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.Error(t, err)
		require.Contains(
			t,
			err.Error(),
			"invalid address 0x0000000000000000000000000000000000000000, reason: empty address",
		)
	})

	t.Run("should return an error when passing incorrect validator", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create staker.
		stakerEVMAddr := sample.EthAddress()

		/* ACT */
		// Call claimRewardsMethod.
		claimRewardsMethod := s.stkContractABI.Methods[ClaimRewardsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			claimRewardsMethod,
			[]interface{}{stakerEVMAddr, "cosmosvaloper100000000000000000000000000000000000000"}...,
		)

		_, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "decoding bech32 failed")
	})

	t.Run("should return an error when there's no delegation", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))

		// Create staker.
		stakerEVMAddr := sample.EthAddress()

		/* ACT */
		// Call claimRewardsMethod.
		claimRewardsMethod := s.stkContractABI.Methods[ClaimRewardsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			claimRewardsMethod,
			[]interface{}{stakerEVMAddr, validator.OperatorAddress}...,
		)

		_, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.Error(t, err)
		require.Contains(
			t,
			err.Error(),
			"unexpected error in WithdrawDelegationRewards: no delegation distribution info",
		)
	})
}
