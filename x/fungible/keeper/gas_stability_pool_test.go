package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_EnsureGasStabilityPoolAccountCreated(t *testing.T) {
	t.Run("can create the gas stability pool account if doesn't exist", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.FungibleKeeper(t)

		// account doesn't exist
		acc := k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress())
		assert.Nil(t, acc)

		// create the account
		k.EnsureGasStabilityPoolAccountCreated(ctx)
		acc = k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress())
		assert.NotNil(t, acc)
		assert.Equal(t, types.GasStabilityPoolAddress(), acc.GetAddress())

		// can call the method again without side effects
		k.EnsureGasStabilityPoolAccountCreated(ctx)
		acc2 := k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress())
		assert.NotNil(t, acc2)
		assert.True(t, acc.GetAddress().Equals(acc2.GetAddress()))
		assert.Equal(t, acc.GetAccountNumber(), acc2.GetAccountNumber())
		assert.Equal(t, acc.GetSequence(), acc2.GetSequence())
	})
}

func TestKeeper_FundGasStabilityPool(t *testing.T) {
	t.Run("can fund the gas stability pool and withdraw", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		// deploy the system contracts and gas coin
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// balance is initially 0
		balance, err := k.GetGasStabilityPoolBalance(ctx, chainID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), balance.Int64())

		// fund the gas stability pool
		err = k.FundGasStabilityPool(ctx, chainID, big.NewInt(100))
		assert.NoError(t, err)

		// balance is now 100
		balance, err = k.GetGasStabilityPoolBalance(ctx, chainID)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), balance.Int64())

		// withdraw from the gas stability pool
		err = k.WithdrawFromGasStabilityPool(ctx, chainID, big.NewInt(50))
		assert.NoError(t, err)

		// balance is now 50
		balance, err = k.GetGasStabilityPoolBalance(ctx, chainID)
		assert.NoError(t, err)
		assert.Equal(t, int64(50), balance.Int64())
	})
}
