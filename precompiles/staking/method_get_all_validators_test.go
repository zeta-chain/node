package staking

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetAllValidators(t *testing.T) {
	t.Run("should return empty array if validators not set", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)

		// Clean all validators.
		validatorsList := s.sdkKeepers.StakingKeeper.GetAllValidators(s.ctx)
		for _, v := range validatorsList {
			s.sdkKeepers.StakingKeeper.RemoveValidator(s.ctx, v.GetOperator())
		}

		methodID := s.contractABI.Methods[GetAllValidatorsMethodName]
		s.mockVMContract.Input = methodID.ID

		// ACT
		validators, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.Empty(t, res[0])
	})

	t.Run("should return validators if set", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.contractABI.Methods[GetAllValidatorsMethodName]
		s.mockVMContract.Input = methodID.ID
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// ACT
		validators, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.NotEmpty(t, res[0])
	})
}
