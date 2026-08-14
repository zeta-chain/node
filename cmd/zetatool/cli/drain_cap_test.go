package cli

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/drain"
)

func TestParseEVMMaxAmount(t *testing.T) {
	t.Run("empty means no cap", func(t *testing.T) {
		v, err := parseEVMMaxAmount("")
		require.NoError(t, err)
		require.True(t, v.IsZero())
	})

	t.Run("whitespace is trimmed", func(t *testing.T) {
		v, err := parseEVMMaxAmount("  1000  ")
		require.NoError(t, err)
		require.Equal(t, sdkmath.NewUint(1000), v)
	})

	t.Run("parses a wei amount beyond int64", func(t *testing.T) {
		v, err := parseEVMMaxAmount("100000000000000000000")
		require.NoError(t, err)
		require.Equal(t, "100000000000000000000", v.String())
	})

	t.Run("an explicit zero is rejected", func(t *testing.T) {
		// a zero cap would pin a transfer that pays gas to move nothing; the operator almost
		// certainly meant to omit the flag
		_, err := parseEVMMaxAmount("0")
		require.ErrorContains(t, err, "must be positive")
	})

	t.Run("a value above 256 bits is rejected, not panicked on", func(t *testing.T) {
		// sdkmath.Uint's constructor panics past 256 bits, so an extra-zeros typo would take
		// down zetatool mid-drain instead of reporting a bad flag
		maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

		v, err := parseEVMMaxAmount(maxUint256.String())
		require.NoError(t, err)
		require.Equal(t, maxUint256.String(), v.String())

		require.NotPanics(t, func() {
			_, err = parseEVMMaxAmount(new(big.Int).Add(maxUint256, big.NewInt(1)).String())
		})
		require.ErrorContains(t, err, "exceeds the maximum 256-bit")
	})

	t.Run("rejects non-decimal and negative values", func(t *testing.T) {
		for _, s := range []string{"0x64", "1e18", "abc", "-1", "1.5"} {
			_, err := parseEVMMaxAmount(s)
			require.ErrorContains(t, err, "invalid --evm-max-amount", "input %q", s)
		}
	})
}

func TestWarnIfCapped(t *testing.T) {
	t.Run("silent when nothing is capped", func(t *testing.T) {
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.ZeroUint(), 0)
		require.Empty(t, buf.String())
	})

	t.Run("an unset cap reads as uncapped, not as zero", func(t *testing.T) {
		// "btc 0 sats total" would read as "sends no BTC", which is the opposite of the truth:
		// BTC is swept in full when only the EVM cap is set
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.NewUint(1_000), 0)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL PAYLOAD")
		require.Contains(t, out, "evm: 1000 wei per chain")
		require.Contains(t, out, "btc: uncapped")
	})

	t.Run("btc-only cap", func(t *testing.T) {
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.ZeroUint(), 100_000)

		out := buf.String()
		require.Contains(t, out, "evm: uncapped")
		require.Contains(t, out, "btc: 100000 sats total")
	})

	t.Run("does not claim nothing is drained", func(t *testing.T) {
		// a cap is an upper bound: a chain holding less than it is swept in full
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.NewUint(1_000), 100_000)
		require.Contains(t, buf.String(), "except on chains already below the cap")
	})
}

func TestReportUnviableBTCCap(t *testing.T) {
	payee, err := btcutil.NewAddressWitnessPubKeyHash(make([]byte, 20), &chaincfg.RegressionNetParams)
	require.NoError(t, err)

	const feeRate = int64(10)
	minViable, err := drain.MinViableSweepSats(feeRate, payee)
	require.NoError(t, err)

	t.Run("points at the smallest cap that would work", func(t *testing.T) {
		// mainnet-shaped: large UTXOs far above the cap, the rest dust far below the threshold
		utxos := []drain.UTXO{
			{TxID: "large", Vout: 0, AmountSats: 8_828_175},
			{TxID: "dust", Vout: 1, AmountSats: 931},
		}

		var buf bytes.Buffer
		reportUnviableBTCCap(&buf, utxos, 100_000, feeRate, payee, 8332)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL SWEPT NO BTC")
		require.Contains(t, out, fmt.Sprintf("at least %d sats", minViable))
		// the actionable number: the smallest UTXO that can be swept on its own
		require.Contains(t, out, "8828175")
	})

	t.Run("says so when no cap can rehearse BTC", func(t *testing.T) {
		// every UTXO is below the viability threshold, so raising the cap cannot help
		utxos := []drain.UTXO{
			{TxID: "dust-0", Vout: 0, AmountSats: 931},
			{TxID: "dust-1", Vout: 1, AmountSats: 500},
		}

		var buf bytes.Buffer
		reportUnviableBTCCap(&buf, utxos, 100_000, feeRate, payee, 8332)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL SWEPT NO BTC")
		require.Contains(t, out, "no cap can rehearse BTC")
		require.Contains(t, out, "can only be drained uncapped")
	})
}
