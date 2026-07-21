package migration_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/migration"
)

func TestComputeEVMMigration(t *testing.T) {
	t.Run("computes amount, gas price and gas limit", func(t *testing.T) {
		// ARRANGE
		balance := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100_000)

		// ACT
		amount, gasPrice, gasLimit, err := migration.ComputeEVMMigration(balance, medianGasPrice)

		// ASSERT
		require.NoError(t, err)
		require.EqualValues(t, 21_000, gasLimit)
		// medianGasPrice × 2.5
		require.Equal(t, sdkmath.NewUint(250_000).String(), gasPrice.String())
		// fee = 21000 × 250000 + 2_100_000_000 = 7_350_000_000
		expectedFee := sdkmath.NewUint(7_350_000_000)
		require.Equal(t, balance.Sub(expectedFee).String(), amount.String())
	})

	t.Run("errors when fee exceeds balance", func(t *testing.T) {
		// ARRANGE
		balance := sdkmath.NewUint(100_000_000)
		medianGasPrice := sdkmath.NewUint(100_000)

		// ACT
		_, _, _, err := migration.ComputeEVMMigration(balance, medianGasPrice)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, migration.ErrInsufficientFunds)
	})
}

func TestComputeBTCMigration(t *testing.T) {
	t.Run("computes output and fee for a sweep", func(t *testing.T) {
		// ARRANGE
		totalInputSats := int64(100_000_000)
		feeRate := int64(10)

		// ACT
		outputSats, feeSats, err := migration.ComputeBTCMigration(totalInputSats, feeRate)

		// ASSERT
		require.NoError(t, err)
		require.Positive(t, feeSats)
		require.Equal(t, totalInputSats-feeSats, outputSats)
	})

	t.Run("overhead at conservative fee rate matches tss-balances", func(t *testing.T) {
		// ARRANGE: at 50 sat/vB the total overhead must equal the value tss-balances reports:
		// OutboundBytesMax*50 (77_150) + RBF reserve (1_000_000) + nonce-mark buffer (3_000).
		totalInputSats := int64(1_000_000_000)

		// ACT
		outputSats, feeSats, err := migration.ComputeBTCMigration(totalInputSats, migration.BTCConservativeFeeRate)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, int64(1_080_150), feeSats)
		require.Equal(t, totalInputSats-int64(1_080_150), outputSats)
	})

	t.Run("errors when inputs cannot cover fee", func(t *testing.T) {
		// ARRANGE
		totalInputSats := int64(1)
		feeRate := int64(100)

		// ACT
		_, _, err := migration.ComputeBTCMigration(totalInputSats, feeRate)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, migration.ErrInsufficientFunds)
	})
}
