package staking

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_GetShares(t *testing.T) {
	// Disabled temporarily because the staking functions were disabled.
	// Issue: https://github.com/zeta-chain/node/issues/3009
	// t.Run("should return stakes", func(t *testing.T) {
	// 	// ARRANGE
	// 	s := newTestSuite(t)
	// 	methodID := s.contractABI.Methods[GetSharesMethodName]
	// 	r := rand.New(rand.NewSource(42))
	// 	validator := sample.Validator(t, r)
	// 	s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

	// 	staker := sample.Bech32AccAddress()
	// 	stakerEthAddr := common.BytesToAddress(staker.Bytes())
	// 	coins := sample.Coins()
	// 	err := s.sdkKeepers.BankKeeper.MintCoins(s.ctx, fungibletypes.ModuleName, sample.Coins())
	// 	require.NoError(t, err)
	// 	err = s.sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, fungibletypes.ModuleName, staker, coins)
	// 	require.NoError(t, err)

	// 	stakerAddr := common.BytesToAddress(staker.Bytes())

	// 	stakeArgs := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

	// 	stakeMethodID := s.contractABI.Methods[StakeMethodName]

	// 	// ACT
	// 	_, err = s.contract.Stake(s.ctx, s.mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, stakeArgs)
	// 	require.NoError(t, err)

	// 	// ASSERT
	// 	args := []interface{}{stakerEthAddr, validator.OperatorAddress}
	// 	s.mockVMContract.Input = packInputArgs(t, methodID, args...)
	// 	stakes, err := s.contract.Run(s.mockEVM, s.mockVMContract, false)
	// 	require.NoError(t, err)

	// 	res, err := methodID.Outputs.Unpack(stakes)
	// 	require.NoError(t, err)
	// 	require.Equal(
	// 		t,
	// 		fmt.Sprintf("%d000000000000000000", coins.AmountOf(config.BaseDenom).BigInt().Int64()),
	// 		res[0].(*big.Int).String(),
	// 	)
	// })

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.contractABI.Methods[GetSharesMethodName]
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		args := []interface{}{stakerEthAddr}

		// ACT
		_, err := s.contract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if invalid staker arg", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.contractABI.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		args := []interface{}{42, validator.OperatorAddress}

		// ACT
		_, err := s.contract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if invalid val address", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.contractABI.Methods[GetSharesMethodName]
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		args := []interface{}{stakerEthAddr, staker.String()}

		// ACT
		_, err := s.contract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if invalid val address format", func(t *testing.T) {
		// ARRANGE
		s := newTestSuite(t)
		methodID := s.contractABI.Methods[GetSharesMethodName]
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		args := []interface{}{stakerEthAddr, 42}

		// ACT
		_, err := s.contract.GetShares(s.ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})
}
