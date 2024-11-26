package staking

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetValidators(t *testing.T) {
	t.Run("should return an empty list for a non staker address", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create validator.
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// Create staker.
		stakerEVMAddr := sample.EthAddress()

		/* ACT */
		// Call getValidatorListForDelegator.
		getValidatorsMethod := s.stkContractABI.Methods[GetValidatorsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			getValidatorsMethod,
			[]interface{}{stakerEVMAddr}...,
		)

		bytes, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		res, err := getValidatorsMethod.Outputs.Unpack(bytes)
		require.NoError(t, err)
		require.NotEmpty(t, res)

		list, ok := res[0].([]string)
		require.True(t, ok)
		require.Len(t, list, 0)
	})

	t.Run("should return an error for zero address", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create validator.
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		/* ACT */
		// Call getValidatorListForDelegator.
		getValidatorsMethod := s.stkContractABI.Methods[GetValidatorsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			getValidatorsMethod,
			[]interface{}{common.Address{}}...,
		)

		_, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.Error(t, err)
		require.Contains(
			t,
			err.Error(),
			"invalid address 0x0000000000000000000000000000000000000000, reason: empty address",
		)
	})

	t.Run("should return staker's validator list", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create validator.
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// Create staker.
		stakerEVMAddr := sample.EthAddress()
		stakerCosmosAddr, err := precompiletypes.GetCosmosAddress(s.sdkKeepers.BankKeeper, stakerEVMAddr)
		require.NoError(t, err)

		// Become a staker.
		stakeThroughCosmosAPI(
			t,
			s.ctx,
			s.sdkKeepers.BankKeeper,
			s.sdkKeepers.StakingKeeper,
			validator,
			stakerCosmosAddr,
			math.NewInt(100),
		)

		/* ACT */
		// Call getValidatorListForDelegator.
		getValidatorsMethod := s.stkContractABI.Methods[GetValidatorsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			getValidatorsMethod,
			[]interface{}{stakerEVMAddr}...,
		)

		bytes, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		res, err := getValidatorsMethod.Outputs.Unpack(bytes)
		require.NoError(t, err)
		require.NotEmpty(t, res)

		list, ok := res[0].([]string)
		require.True(t, ok)
		require.Len(t, list, 1)
		require.Equal(t, validator.GetOperator().String(), list[0])
	})

	t.Run(" should return staker's validator list - heavy test with 100 validators", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create staker.
		stakerEVMAddr := sample.EthAddress()
		stakerCosmosAddr, err := precompiletypes.GetCosmosAddress(s.sdkKeepers.BankKeeper, stakerEVMAddr)
		require.NoError(t, err)

		// Create 100 validators, and stake on each of them.
		for n := range 100 {
			validator := sample.Validator(t, rand.New(rand.NewSource(int64(n))))
			s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

			stakeThroughCosmosAPI(
				t,
				s.ctx,
				s.sdkKeepers.BankKeeper,
				s.sdkKeepers.StakingKeeper,
				validator,
				stakerCosmosAddr,
				math.NewInt(100),
			)
		}

		/* ACT */
		// Call getValidatorListForDelegator.
		getValidatorsMethod := s.stkContractABI.Methods[GetValidatorsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			getValidatorsMethod,
			[]interface{}{stakerEVMAddr}...,
		)

		bytes, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		res, err := getValidatorsMethod.Outputs.Unpack(bytes)
		require.NoError(t, err)
		require.NotEmpty(t, res)

		list, ok := res[0].([]string)
		require.True(t, ok)

		// The returned list should contain 100 entries.
		require.Len(t, list, 100)
	})
}
