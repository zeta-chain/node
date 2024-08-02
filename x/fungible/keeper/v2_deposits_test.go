package keeper_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
	"testing"
)

func TestKeeper_ProcessV2Deposit(t *testing.T) {
	t.Run("should process no-call deposit", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		receiver := sample.EthAddress()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, contractCall, err := k.ProcessV2Deposit(ctx, zrc20, receiver, big.NewInt(42), []byte{})

		// ASSERT
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, receiver)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should fail if not recognized action", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		receiver := sample.EthAddress()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, _, err := k.ProcessV2Deposit(ctx, zrc20, receiver, big.NewInt(42), sample.Bytes())

		// ASSERT
		require.ErrorContains(t, err, "not implemented")
	})
}
