package keeper_test

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_RefundRemainingGasFees(t *testing.T) {
	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		receiver := ethcommon.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")

		err := k.RefundRemainingGasFees(ctx, chainID, big.NewInt(100), receiver)
		require.Error(t, err)
	})

	t.Run("can refund remaining gas fees to receiver", func(t *testing.T) {
		// Arrange
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		receiver := ethcommon.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		receiverBalanceBefore, err := k.BalanceOfZRC4(ctx, gasZRC20, receiver)
		require.NoError(t, err)
		require.Equal(t, int64(0), receiverBalanceBefore.Int64())

		// Act
		err = k.RefundRemainingGasFees(ctx, chainID, big.NewInt(100), receiver)
		require.NoError(t, err)

		// Assert
		gasZRC20, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		receiverBalanceAfter, err := k.BalanceOfZRC4(ctx, gasZRC20, receiver)
		require.NoError(t, err)
		require.Equal(t, int64(100), receiverBalanceAfter.Int64())
	})
}
