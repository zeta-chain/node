package staking

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	ptypes "github.com/zeta-chain/node/precompiles/types"
)

func Test_Distribute(t *testing.T) {
	t.Run("should fail to run distribute as read only method", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(0)}...,
		)

		// Call method as read only.
		result, err := s.contract.Run(s.mockEVM, s.mockVMContract, true)

		// Check error and result.
		require.ErrorIs(t, err, ptypes.ErrWriteMethod{
			Method: DistributeMethodName,
		})

		// Result is empty as the write check is done before executing distribute() function.
		// On-chain this would look like reverting, so staticcall is properly reverted.
		require.Empty(t, result)
	})

	t.Run("should fail to distribute with 0 token balance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(0)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.ErrorAs(
			t,
			ptypes.ErrInvalidAmount{
				Got: "0",
			},
			err,
		)

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute with 0 allowance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1000)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 0")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute 0 token", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(0)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid token amount: 0")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute more than allowed to staking", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(999))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1000)}...,
		)

		// Call method.
		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 999, wanted 1000")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should fail to distribute more than user balance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(100000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1001)}...,
		)

		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)
	})

	t.Run("should distribute and lock ZRC20", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)

		// Set caller balance.
		s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, big.NewInt(1000))

		// Allow staking to spend ZRC20 tokens.
		allowStaking(t, s, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			s.methodID,
			[]interface{}{s.zrc20Address, big.NewInt(1000)}...,
		)

		success, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.NoError(t, err)

		// Unpack and check result boolean.
		res, err := s.methodID.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.True(t, ok)

		// feeCollectorAddress := s.sdkKeepers.AuthKeeper.GetModuleAccount(s.ctx, types.FeeCollectorName).GetAddress()
		// coin := s.sdkKeepers.BankKeeper.GetBalance(s.ctx, feeCollectorAddress, ptypes.ZRC20ToCosmosDenom(s.zrc20Address))
		// fmt.Println(coin)
	})
}
