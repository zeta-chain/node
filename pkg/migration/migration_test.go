package migration_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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
	payee, err := btcutil.NewAddressWitnessPubKeyHash(make([]byte, 20), &chaincfg.RegressionNetParams)
	require.NoError(t, err)

	t.Run("computes output and fee for a sweep", func(t *testing.T) {
		// ARRANGE
		totalInputSats := int64(100_000_000)
		feeRate := int64(10)
		numInputs := 3

		// ACT
		outputSats, feeSats, err := migration.ComputeBTCMigration(totalInputSats, feeRate, numInputs, payee)

		// ASSERT
		require.NoError(t, err)
		require.Positive(t, feeSats)
		require.Equal(t, totalInputSats-feeSats, outputSats)
	})

	t.Run("errors when inputs cannot cover fee", func(t *testing.T) {
		// ARRANGE
		totalInputSats := int64(1)
		feeRate := int64(100)
		numInputs := 20

		// ACT
		_, _, err := migration.ComputeBTCMigration(totalInputSats, feeRate, numInputs, payee)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, migration.ErrInsufficientFunds)
	})
}
