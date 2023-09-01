package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_EnsureGasStabilityPoolAccountCreated(t *testing.T) {
	t.Run("can create the gas stability pool account if doesn't exist", func(t *testing.T) {
		k, ctx, _ := testkeeper.FungibleKeeper(t)

		// account doesn't exist
		acc := k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress())
		require.Nil(t, acc)

		// create the account
		k.EnsureGasStabilityPoolAccountCreated(ctx)
		acc = k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress())
		require.NotNil(t, acc)
		require.Equal(t, types.GasStabilityPoolAddress(), acc.GetAddress())

		// can call the method again without side effects
		k.EnsureGasStabilityPoolAccountCreated(ctx)
		acc2 := k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress())
		require.NotNil(t, acc2)
		require.True(t, acc.GetAddress().Equals(acc2.GetAddress()))
		require.Equal(t, acc.GetAccountNumber(), acc2.GetAccountNumber())
		require.Equal(t, acc.GetSequence(), acc2.GetSequence())
	})
}

func TestKeeper_FundGasStabilityPool(t *testing.T) {
	t.Run("can fund the gas stability pool and withdraw", func(t *testing.T) {
		k, ctx, sdkk := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		// deploy the system contracts and gas coin
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// balance is initially 0
		balance, err := k.GetGasStabilityPoolBalance(ctx, chainID)
		require.NoError(t, err)
		require.Equal(t, int64(0), balance.Int64())

		// fund the gas stability pool
		err = k.FundGasStabilityPool(ctx, chainID, big.NewInt(100))
		require.NoError(t, err)

		// balance is now 100
		balance, err = k.GetGasStabilityPoolBalance(ctx, chainID)
		require.NoError(t, err)
		require.Equal(t, int64(100), balance.Int64())

		// withdraw from the gas stability pool
		err = k.WithdrawFromGasStabilityPool(ctx, chainID, big.NewInt(50))
		require.NoError(t, err)

		// balance is now 50
		balance, err = k.GetGasStabilityPoolBalance(ctx, chainID)
		require.NoError(t, err)
		require.Equal(t, int64(50), balance.Int64())
	})
}
