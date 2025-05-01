package staking

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetAllValidators(t *testing.T) {
	t.Run("should return empty array if validators not set", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)

		// Clean all validators.
		validatorsList, err := s.sdkKeepers.StakingKeeper.GetAllValidators(s.ctx)
		require.NoError(t, err)
		for _, v := range validatorsList {
			valAddr, err := sdk.ValAddressFromBech32(v.GetOperator())
			require.NoError(t, err)
			s.sdkKeepers.StakingKeeper.RemoveValidator(s.ctx, valAddr)
		}

		methodID := s.stkContractABI.Methods[GetAllValidatorsMethodName]
		s.mockVMContract.Input = methodID.ID

		// ACT
		validators, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.Empty(t, res[0])
	})

	t.Run("should return validators if set", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.stkContractABI.Methods[GetAllValidatorsMethodName]
		s.mockVMContract.Input = methodID.ID
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		validator := sample.Validator(t, r)
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// ACT
		validators, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.NotEmpty(t, res[0])
	})
}
