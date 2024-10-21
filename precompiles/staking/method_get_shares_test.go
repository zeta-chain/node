package staking

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/testutil/sample"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func Test_GetShares(t *testing.T) {
	t.Run("should return stakes", func(t *testing.T) {
		// ARRANGE
		ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract := setup(t)
		methodID := abi.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		sdkKeepers.StakingKeeper.SetValidator(ctx, validator)

		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		coins := sample.Coins()
		err := sdkKeepers.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
		require.NoError(t, err)

		stakerAddr := common.BytesToAddress(staker.Bytes())

		stakeArgs := []interface{}{stakerEthAddr, validator.OperatorAddress, coins.AmountOf(config.BaseDenom).BigInt()}

		stakeMethodID := abi.Methods[StakeMethodName]

		// ACT
		_, err = contract.Stake(ctx, mockEVM, &vm.Contract{CallerAddress: stakerAddr}, &stakeMethodID, stakeArgs)
		require.NoError(t, err)

		// ASSERT
		args := []interface{}{stakerEthAddr, validator.OperatorAddress}
		mockVMContract.Input = packInputArgs(t, methodID, args...)
		stakes, err := contract.Run(mockEVM, mockVMContract, false)
		require.NoError(t, err)

		res, err := methodID.Outputs.Unpack(stakes)
		require.NoError(t, err)
		require.Equal(
			t,
			fmt.Sprintf("%d000000000000000000", coins.AmountOf(config.BaseDenom).BigInt().Int64()),
			res[0].(*big.Int).String(),
		)
	})

	t.Run("should fail if wrong args amount", func(t *testing.T) {
		// ARRANGE
		ctx, contract, abi, _, _, _ := setup(t)
		methodID := abi.Methods[GetSharesMethodName]
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		args := []interface{}{stakerEthAddr}

		// ACT
		_, err := contract.GetShares(ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if invalid staker arg", func(t *testing.T) {
		// ARRANGE
		ctx, contract, abi, _, _, _ := setup(t)
		methodID := abi.Methods[GetSharesMethodName]
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		args := []interface{}{42, validator.OperatorAddress}

		// ACT
		_, err := contract.GetShares(ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if invalid val address", func(t *testing.T) {
		// ARRANGE
		ctx, contract, abi, _, _, _ := setup(t)
		methodID := abi.Methods[GetSharesMethodName]
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		args := []interface{}{stakerEthAddr, staker.String()}

		// ACT
		_, err := contract.GetShares(ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should fail if invalid val address format", func(t *testing.T) {
		// ARRANGE
		ctx, contract, abi, _, _, _ := setup(t)
		methodID := abi.Methods[GetSharesMethodName]
		staker := sample.Bech32AccAddress()
		stakerEthAddr := common.BytesToAddress(staker.Bytes())
		args := []interface{}{stakerEthAddr, 42}

		// ACT
		_, err := contract.GetShares(ctx, &methodID, args)

		// ASSERT
		require.Error(t, err)
	})
}
