package staking

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetAllValidators(t *testing.T) {
	t.Run("should return empty array if validators not set", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)
		validatorsList := ts.sdkKeepers.StakingKeeper.GetAllValidators(ts.ctx)
		for _, v := range validatorsList {
			fmt.Println(v.OperatorAddress)
			ts.sdkKeepers.StakingKeeper.RemoveValidator(ts.ctx, types.ValAddress(v.OperatorAddress))
		}

		methodID := ts.contractABI.Methods[GetAllValidatorsMethodName]
		ts.mockVMContract.Input = methodID.ID

		// ACT
		validators, err := ts.contract.Run(ts.mockEVM, ts.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.Empty(t, res[0])
	})

	t.Run("should return validators if set", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)
		methodID := ts.contractABI.Methods[GetAllValidatorsMethodName]
		ts.mockVMContract.Input = methodID.ID
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		ts.sdkKeepers.StakingKeeper.SetValidator(ts.ctx, validator)

		// ACT
		validators, err := ts.contract.Run(ts.mockEVM, ts.mockVMContract, false)

		// ASSERT
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(validators)
		require.NoError(t, err)

		require.NotEmpty(t, res[0])
	})
}
